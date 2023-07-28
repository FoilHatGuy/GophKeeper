package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"gophKeeper/src/server/cfg"
	"gophKeeper/src/server/database"
	"gophKeeper/src/server/passwords"
	pb "gophKeeper/src/server/pb"
)

type key string

const (
	uidMetaKey key    = "uid"
	sidMetaKey string = "sid"
)

func RunGRPCServer(config *cfg.ConfigT) {
	db := database.New()
	auth := AuthGRPC{db: db}
	backend := ServerGRPC{db: db}

	//nolint:godox
	server := grpc.NewServer(grpc.UnaryInterceptor(backend.Authenticate)) // TODO add interceptor for session check
	reflection.Register(server)
	pb.RegisterAuthServer(server, &auth)
	pb.RegisterGophKeeperServer(server, &backend)

	lis, err := net.Listen("tcp", config.Server.AddressGRPC)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("GRPC server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// AuthGRPC is a structure containing all required services, as well as embedded server
type AuthGRPC struct {
	pb.UnimplementedAuthServer
	db database.StorageController
}

// Login
// Ping server+database activity
func (s *AuthGRPC) Login(ctx context.Context, in *pb.LoginPassPair) (out *pb.SessionID_DTO, errRPC error) {
	// get all data
	login := in.GetLogin()
	pass := in.GetPassword()
	passStored, err := s.db.GetPassword(ctx, login)
	if errors.Is(err, database.ErrNotFound) {
		return nil, status.Errorf(codes.Unauthenticated, "login or password incorrect")
	}

	// verify password
	ok, err := passwords.ComparePasswordHash(passStored, pass)
	if err != nil || !ok {
		return nil, status.Errorf(codes.Unauthenticated, "login or password incorrect")
	}

	// if ok, add session
	sid := uuid.NewString()
	err = s.db.AddSession(ctx, login, sid)
	switch {
	case errors.Is(err, database.ErrConflict):
		return nil, status.Errorf(codes.AlreadyExists, "already logged in")
	case err != nil:
		return nil, status.Errorf(codes.Internal, "login or password incorrect")
	}
	out = &pb.SessionID_DTO{SID: sid}
	return
}

// KickOtherSession
// Ping server+database activity
func (s *AuthGRPC) KickOtherSession(ctx context.Context, in *pb.LoginPassPair) (out *pb.SessionID_DTO, errRPC error) {
	// get all data
	login := in.GetLogin()
	pass := in.GetPassword()
	passStored, err := s.db.GetPassword(ctx, login)
	if errors.Is(err, database.ErrNotFound) {
		return nil, status.Errorf(codes.Unauthenticated, "login or password incorrect")
	}

	// verify password
	ok, err := passwords.ComparePasswordHash(passStored, pass)
	if err != nil || !ok {
		return nil, status.Errorf(codes.Unauthenticated, "login or password incorrect")
	}

	// if ok, add session
	sid := uuid.NewString()
	err = s.db.UpdateSession(ctx, login, sid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "login or password incorrect")
	}
	out = &pb.SessionID_DTO{SID: sid}
	return
}

// Register
// Ping server+database activity
func (s *AuthGRPC) Register(ctx context.Context, in *pb.LoginPassPair) (out *pb.SessionID_DTO, errRPC error) {
	login := in.GetLogin()
	password := in.GetPassword()
	hashed, err := passwords.HashPassword(password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "password is incorrect")
	}

	err = s.db.AddUser(ctx, login, hashed)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error")
	}
	return
}

// ServerGRPC is a structure containing all required services, as well as embedded server
type ServerGRPC struct {
	pb.UnimplementedGophKeeperServer
	db database.StorageController
}

// Authenticate manages sid cookies.
func (s *ServerGRPC) Authenticate(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (response interface{}, errRPC error) {
	fmt.Printf("%+v\n", info.FullMethod)
	if strings.Contains(strings.ToLower(info.FullMethod), "base.auth") {
		return handler(ctx, req)
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}
	metaValue := md.Get(sidMetaKey)
	if len(metaValue) < 1 {
		return nil, status.Errorf(codes.Unauthenticated, "missing 'user' field in metadata")
	}
	sid := metaValue[0]

	uid, ok, err := s.db.RefreshSession(ctx, sid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "while accessing database")
	}
	if !ok {
		return nil, status.Errorf(codes.ResourceExhausted, "session stale")
	}
	ctx = context.WithValue(ctx, uidMetaKey, uid)
	metadata.AppendToOutgoingContext(ctx, sidMetaKey, sid)

	return handler(ctx, req)
}

// Ping checks server+database activity
func (s *ServerGRPC) Ping(_ context.Context, _ *pb.Empty) (out *pb.Empty, errRPC error) {
	out = &pb.Empty{}
	return
}

// GetCategoryHead returns data about a category
func (s *ServerGRPC) GetCategoryHead(
	ctx context.Context,
	in *pb.CategoryType_DTO,
) (out *pb.CategoryHead_DTO, errRPC error) {
	var (
		head database.CategoryHead
		err  error
	)
	category := in.GetCategory()

	switch category {
	case pb.Category_CATEGORY_PASS:
		head, err = s.db.GetLPPHead(ctx)

	case pb.Category_CATEGORY_TEXT:
		//nolint:godox
		head, err = s.db.GetLPPHead(ctx) // todo: change to correct category
	case pb.Category_CATEGORY_CARD:
		head, err = s.db.GetLPPHead(ctx)
	case pb.Category_CATEGORY_FILE:
		head, err = s.db.GetLPPHead(ctx)
	}

	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return nil, status.Errorf(codes.Internal, "login or password incorrect")
	}
	if errors.Is(err, database.ErrNotFound) {
		return &pb.CategoryHead_DTO{}, nil
	}

	newInfo := make([]*pb.DataInfo, 0, len(head))

	for _, el := range head {
		newInfo = append(newInfo, &pb.DataInfo{
			DataID:   el.DataID,
			Metadata: el.Metadata,
		})
	}
	out = &pb.CategoryHead_DTO{Info: newInfo}
	return out, nil
}

// StoreLoginPassword
// Ping server+database activity
func (s *ServerGRPC) StoreLoginPassword(ctx context.Context, in *pb.SecureData_DTO) (out *pb.DataID_DTO, errRPC error) {
	data := in.GetData()
	meta := in.GetMetadata()
	dataID := uuid.NewString()

	err := s.db.AddLPP(ctx, dataID, meta, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "login or password incorrect")
	}

	out = &pb.DataID_DTO{
		ID: dataID,
	}
	return
}

// LoadLoginPassword
// Ping server+database activity
func (s *ServerGRPC) LoadLoginPassword(ctx context.Context, in *pb.DataID_DTO) (out *pb.SecureData_DTO, errRPC error) {
	id := in.GetID()

	meta, data, err := s.db.GetLPP(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error")
	}
	out = &pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	}
	return
}
