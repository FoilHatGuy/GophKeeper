package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/google/uuid"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	__ "gophKeeper/src/pb"
	"gophKeeper/src/server/cfg"
	"gophKeeper/src/server/database"
	"gophKeeper/src/server/passwords"
)

type key string

const (
	uidMetaKey key    = "uid"
	sidMetaKey string = "sid"
)

func RunGRPCServer(config *cfg.ConfigT, logger *log.Logger) {
	ctx := context.Background()
	db := database.New(ctx, config)
	auth := AuthGRPC{db: db}
	backend := ServerGRPC{db: db}

	logrusEntry := log.NewEntry(logger)
	opts := []grpc_logrus.Option{
		grpc_logrus.WithLevels(func(code codes.Code) log.Level {
			if code == codes.OK {
				return log.InfoLevel
			}
			return log.ErrorLevel
		}),
	}
	grpc_logrus.ReplaceGrpcLogger(logrusEntry)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(logrusEntry, opts...),
			backend.Authenticate,
		),
		//nolint:godox
		grpc.ChainStreamInterceptor( // todo: add stream auth
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.StreamServerInterceptor(logrusEntry, opts...),
		),
	)
	reflection.Register(server)
	__.RegisterAuthServer(server, &auth)
	__.RegisterGophKeeperServer(server, &backend)

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
	__.UnimplementedAuthServer
	db database.StorageController
}

// Login
// Ping server+database activity
func (s *AuthGRPC) Login(ctx context.Context, in *__.Credentials) (out *__.SessionID_DTO, errRPC error) {
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
	out = &__.SessionID_DTO{SID: sid}
	return
}

// KickOtherSession
// Ping server+database activity
func (s *AuthGRPC) KickOtherSession(ctx context.Context, in *__.Credentials) (out *__.SessionID_DTO, errRPC error) {
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
	out = &__.SessionID_DTO{SID: sid}
	return
}

// Register
// Ping server+database activity
func (s *AuthGRPC) Register(ctx context.Context, in *__.Credentials) (out *__.Empty, errRPC error) {
	login := in.GetLogin()
	password := in.GetPassword()
	hashed, err := passwords.HashPassword(password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "password is incorrect")
	}

	uid := uuid.NewString()
	err = s.db.AddUser(ctx, uid, login, hashed)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error")
	}
	out = &__.Empty{}
	return
}

// ServerGRPC is a structure containing all required services, as well as embedded server
type ServerGRPC struct {
	__.UnimplementedGophKeeperServer
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
		return nil, status.Errorf(codes.Unauthenticated, "missing %q field in metadata", sidMetaKey)
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
func (s *ServerGRPC) Ping(_ context.Context, _ *__.Empty) (out *__.Empty, errRPC error) {
	out = &__.Empty{}
	return
}

// GetCategoryHead returns data about a category
func (s *ServerGRPC) GetCategoryHead(
	ctx context.Context,
	in *__.CategoryType_DTO,
) (out *__.CategoryHead_DTO, errRPC error) {
	var (
		head database.CategoryHead
		err  error
	)
	category := in.GetCategory()

	uidTL := ctx.Value(uidMetaKey)
	uid := uidTL.(string)

	switch category {
	case __.Category_CATEGORY_CRED:
		head, err = s.db.GetCredentialsHead(ctx, uid)

	case __.Category_CATEGORY_TEXT:
		//nolint:godox
		head, err = s.db.GetTextHead(ctx, uid)
	case __.Category_CATEGORY_CARD:
		head, err = s.db.GetCardHead(ctx, uid)
		// case pb.Category_CATEGORY_FILE: // todo: change to correct category
		//	head, err = s.db.GetFileHead(ctx, uid)
	}

	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return nil, status.Errorf(codes.Internal, "login or password incorrect")
	}
	if errors.Is(err, database.ErrNotFound) {
		return &__.CategoryHead_DTO{}, nil
	}

	newInfo := make([]*__.DataInfo, 0, len(head))

	for _, el := range head {
		newInfo = append(newInfo, &__.DataInfo{
			DataID:   el.DataID,
			Metadata: el.Metadata,
		})
	}
	out = &__.CategoryHead_DTO{Info: newInfo}
	return out, nil
}

// Credentials section

// StoreCredentials
// Ping server+database activity
func (s *ServerGRPC) StoreCredentials(ctx context.Context, in *__.SecureData_DTO) (out *__.DataID_DTO, errRPC error) {
	data := in.GetData()
	meta := in.GetMetadata()
	dataID := uuid.NewString()

	uidTL := ctx.Value(uidMetaKey)
	uid := uidTL.(string)

	err := s.db.AddCredentials(ctx, uid, dataID, meta, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "login or password incorrect")
	}

	out = &__.DataID_DTO{
		ID: dataID,
	}
	return
}

// LoadCredentials
// Ping server+database activity
func (s *ServerGRPC) LoadCredentials(ctx context.Context, in *__.DataID_DTO) (out *__.SecureData_DTO, errRPC error) {
	sid := in.GetID()

	uidTL := ctx.Value(uidMetaKey)
	uid := uidTL.(string)

	meta, data, err := s.db.GetCredentials(ctx, uid, sid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error")
	}
	out = &__.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	}
	return
}

// Text section

// StoreText
// Ping server+database activity
func (s *ServerGRPC) StoreText(ctx context.Context, in *__.SecureData_DTO) (out *__.DataID_DTO, errRPC error) {
	data := in.GetData()
	meta := in.GetMetadata()
	dataID := uuid.NewString()

	uidTL := ctx.Value(uidMetaKey)
	uid := uidTL.(string)

	err := s.db.AddText(ctx, uid, dataID, meta, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "login or password incorrect")
	}

	out = &__.DataID_DTO{
		ID: dataID,
	}
	return
}

// LoadText
// Ping server+database activity
func (s *ServerGRPC) LoadText(ctx context.Context, in *__.DataID_DTO) (out *__.SecureData_DTO, errRPC error) {
	sid := in.GetID()

	uidTL := ctx.Value(uidMetaKey)
	uid := uidTL.(string)

	meta, data, err := s.db.GetText(ctx, uid, sid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error")
	}
	out = &__.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	}
	return
}

// Card section

// StoreCard
// Ping server+database activity
func (s *ServerGRPC) StoreCard(ctx context.Context, in *__.SecureData_DTO) (out *__.DataID_DTO, errRPC error) {
	data := in.GetData()
	meta := in.GetMetadata()
	dataID := uuid.NewString()

	uidTL := ctx.Value(uidMetaKey)
	uid := uidTL.(string)

	err := s.db.AddCard(ctx, uid, dataID, meta, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "login or password incorrect")
	}

	out = &__.DataID_DTO{
		ID: dataID,
	}
	return
}

// LoadCard
// Ping server+database activity
func (s *ServerGRPC) LoadCard(ctx context.Context, in *__.DataID_DTO) (out *__.SecureData_DTO, errRPC error) {
	sid := in.GetID()

	uidTL := ctx.Value(uidMetaKey)
	uid := uidTL.(string)

	meta, data, err := s.db.GetCard(ctx, uid, sid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error")
	}
	out = &__.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	}
	return
}
