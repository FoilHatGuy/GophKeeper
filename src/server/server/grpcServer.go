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

	pb "gophKeeper/src/pb"
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
			// Async,
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
	pb.RegisterAuthServer(server, &auth)
	pb.RegisterGophKeeperServer(server, &backend)

	lis, err := net.Listen("tcp", config.Server.Address)
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

// Ping checks server activity
func (s *AuthGRPC) Ping(_ context.Context, _ *pb.Empty) (out *pb.Empty, errRPC error) {
	out = &pb.Empty{}
	return
}

// Login
// Ping server+database activity
func (s *AuthGRPC) Login(ctx context.Context, in *pb.Credentials) (out *pb.SessionID_DTO, errRPC error) {
	// get all data
	login := in.GetLogin()
	pass := in.GetPassword()
	uid, passStored, err := s.db.GetUserData(ctx, login)
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
	err = s.db.AddSession(ctx, uid, sid)
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
func (s *AuthGRPC) KickOtherSession(ctx context.Context, in *pb.Credentials) (out *pb.SessionID_DTO, errRPC error) {
	// get all data
	login := in.GetLogin()
	pass := in.GetPassword()
	uid, passStored, err := s.db.GetUserData(ctx, login)
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
	err = s.db.UpdateSession(ctx, uid, sid)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "login or password incorrect")
	}
	out = &pb.SessionID_DTO{SID: sid}
	return
}

// Register
// Ping server+database activity
func (s *AuthGRPC) Register(ctx context.Context, in *pb.Credentials) (out *pb.Empty, errRPC error) {
	login := in.GetLogin()
	password := in.GetPassword()
	hashed, err := passwords.HashPassword(password)
	// untestable part, rand.Reader needs to return error
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "password is incorrect")
	}

	uid := uuid.NewString()
	err = s.db.AddUser(ctx, uid, login, hashed)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "database error")
	}
	out = &pb.Empty{}
	return
}

// ServerGRPC is a structure containing all required services, as well as embedded server
type ServerGRPC struct {
	pb.UnimplementedGophKeeperServer
	db database.StorageController
}

// Async launches each handle as a goroutine to prevent panic from stopping the entire server
// func Async(
// 	ctx context.Context,
// 	req interface{},
// 	_ *grpc.UnaryServerInfo,
// 	handler grpc.UnaryHandler,
// ) (response interface{}, errRPC error) {
// 	c := make(chan bool)
// 	go func() {
// 		response, errRPC = handler(ctx, req)
// 		c <- true
// 	}()
// 	<-c
// 	return
// }

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
	ctx = context.WithValue(ctx, uidMetaKey, uid)

	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "session refresh failed")
	}
	if !ok {
		return nil, status.Errorf(codes.ResourceExhausted, "session stale")
	}
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(sidMetaKey, sid))

	return handler(ctx, req)
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

	uidTL := ctx.Value(uidMetaKey)
	uid := uidTL.(string)

	switch category {
	case pb.Category_CATEGORY_CRED:
		head, err = s.db.GetCredentialsHead(ctx, uid)

	case pb.Category_CATEGORY_TEXT:
		head, err = s.db.GetTextHead(ctx, uid)

	case pb.Category_CATEGORY_CARD:
		head, err = s.db.GetCardHead(ctx, uid)

		// case pb.Category_CATEGORY_FILE: // todo: change to correct category
		//	head, err = s.db.GetFileHead(ctx, uid)
	default:
		return nil, status.Errorf(codes.Internal, "database error")
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error")
	}

	newInfo := make([]*pb.DataInfo, 0, len(head))

	for _, el := range head {
		newInfo = append(newInfo, &pb.DataInfo{
			DataID:   el.ID,
			Metadata: el.Metadata,
		})
	}
	out = &pb.CategoryHead_DTO{Info: newInfo}
	return out, nil
}

// Credentials section

// StoreCredData
// Ping server+database activity
func (s *ServerGRPC) StoreCredData(ctx context.Context, in *pb.SecureData_DTO) (out *pb.DataID_DTO, errRPC error) {
	data := in.GetData()
	meta := in.GetMetadata()
	dataID := uuid.NewString()

	uidTypeLess := ctx.Value(uidMetaKey)
	uid := uidTypeLess.(string)

	err := s.db.AddCredentials(ctx, uid, dataID, meta, data)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "database error")
	}

	out = &pb.DataID_DTO{
		ID: dataID,
	}
	return
}

// LoadCredData
// Ping server+database activity
func (s *ServerGRPC) LoadCredData(ctx context.Context, in *pb.DataID_DTO) (out *pb.SecureData_DTO, errRPC error) {
	dataID := in.GetID()

	uidTL := ctx.Value(uidMetaKey)
	uid := uidTL.(string)

	meta, data, err := s.db.GetCredentials(ctx, uid, dataID)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "database error")
	}
	out = &pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	}
	return
}

// Text section

// StoreTextData
// Ping server+database activity
func (s *ServerGRPC) StoreTextData(ctx context.Context, in *pb.SecureData_DTO) (out *pb.DataID_DTO, errRPC error) {
	data := in.GetData()
	meta := in.GetMetadata()
	dataID := uuid.NewString()

	uidTL := ctx.Value(uidMetaKey)
	uid := uidTL.(string)

	err := s.db.AddText(ctx, uid, dataID, meta, data)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "login or password incorrect")
	}

	out = &pb.DataID_DTO{
		ID: dataID,
	}
	return
}

// LoadTextData
// Ping server+database activity
func (s *ServerGRPC) LoadTextData(ctx context.Context, in *pb.DataID_DTO) (out *pb.SecureData_DTO, errRPC error) {
	dataID := in.GetID()

	uidTL := ctx.Value(uidMetaKey)
	uid := uidTL.(string)

	meta, data, err := s.db.GetText(ctx, uid, dataID)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "database error")
	}
	out = &pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	}
	return
}

// Card section

// StoreCardData
// Ping server+database activity
func (s *ServerGRPC) StoreCardData(ctx context.Context, in *pb.SecureData_DTO) (out *pb.DataID_DTO, errRPC error) {
	data := in.GetData()
	meta := in.GetMetadata()
	dataID := uuid.NewString()

	uidTL := ctx.Value(uidMetaKey)
	uid := uidTL.(string)

	err := s.db.AddCard(ctx, uid, dataID, meta, data)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "login or password incorrect")
	}

	out = &pb.DataID_DTO{
		ID: dataID,
	}
	return
}

// LoadCardData
// Ping server+database activity
func (s *ServerGRPC) LoadCardData(ctx context.Context, in *pb.DataID_DTO) (out *pb.SecureData_DTO, errRPC error) {
	dataID := in.GetID()

	uidTL := ctx.Value(uidMetaKey)
	uid := uidTL.(string)

	meta, data, err := s.db.GetCard(ctx, uid, dataID)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "database error")
	}
	out = &pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	}
	return
}

// File section

// StoreFileData
// Ping server+database activity
func (s *ServerGRPC) StoreFileData(ctx context.Context, in *pb.SecureData_DTO) (out *pb.DataID_DTO, errRPC error) {
	data := in.GetData()
	meta := in.GetMetadata()
	dataID := uuid.NewString()

	uidTL := ctx.Value(uidMetaKey)
	uid := uidTL.(string)

	err := s.db.AddFile(ctx, uid, dataID, meta, data)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "login or password incorrect")
	}

	out = &pb.DataID_DTO{
		ID: dataID,
	}
	return
}

// LoadFileData
// Ping server+database activity
func (s *ServerGRPC) LoadFileData(ctx context.Context, in *pb.DataID_DTO) (out *pb.SecureData_DTO, errRPC error) {
	dataID := in.GetID()

	uidTL := ctx.Value(uidMetaKey)
	uid := uidTL.(string)

	meta, data, err := s.db.GetFile(ctx, uid, dataID)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "database error")
	}
	out = &pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	}
	return
}
