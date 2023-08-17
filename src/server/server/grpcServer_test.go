//go:build unit

package app

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "gophKeeper/src/pb"
	"gophKeeper/src/server/database"
	"gophKeeper/src/server/passwords"
)

var dbErr = errors.New("db err")

type gRPCServerTestSuite struct {
	suite.Suite
	ctx        context.Context
	db         *MockStorageController
	serverKeep pb.GophKeeperServer
	serverAuth pb.AuthServer
}

func (s *gRPCServerTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.db = NewMockStorageController(gomock.NewController(s.T()))
	s.serverKeep = &ServerGRPC{
		db: s.db,
	}
	s.serverAuth = &AuthGRPC{
		db: s.db,
	}
}

func (s *gRPCServerTestSuite) TestPing() {
	_, err := s.serverAuth.Ping(s.ctx, &pb.Empty{})
	s.Assert().NoError(err)
}

func (s *gRPCServerTestSuite) TestAuthRegister() {
	const (
		login = "login"
		pass  = "password"
	)
	s.db.EXPECT().AddUser(s.ctx, gomock.Any(), login, gomock.Any()).Return(nil)

	registerRes, err := s.serverAuth.Register(s.ctx, &pb.Credentials{
		Login:    login,
		Password: pass,
	})
	s.Assert().NoError(err)
	s.Assert().Equal(&pb.Empty{}, registerRes)

	// if db returned err
	s.db.EXPECT().AddUser(s.ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(dbErr)
	_, err = s.serverAuth.Register(s.ctx, &pb.Credentials{
		Login:    login,
		Password: pass,
	})
	s.Assert().Error(err)
	stat, _ := status.FromError(err)
	s.Assert().Equal(codes.Internal, stat.Code())
}

func (s *gRPCServerTestSuite) TestAuthLogin() {
	const (
		login = "login"
		pass  = "password"
	)
	hashed, _ := passwords.HashPassword(pass)

	// login
	uid := uuid.NewString()
	var sid string
	s.db.EXPECT().GetUserData(s.ctx, login).Return(uid, hashed, nil)
	s.db.EXPECT().AddSession(s.ctx, uid, gomock.Any()).DoAndReturn(
		func(_ context.Context, _, in string) error {
			sid = in
			return nil
		},
	)
	loginRes, err := s.serverAuth.Login(s.ctx, &pb.Credentials{
		Login:    login,
		Password: pass,
	})
	s.Assert().NoError(err)
	s.Assert().Equal(sid, loginRes.SID)

	// login not found
	s.db.EXPECT().GetUserData(s.ctx, login).Return(uid, hashed, database.ErrNotFound)
	loginRes, err = s.serverAuth.Login(s.ctx, &pb.Credentials{
		Login:    login,
		Password: pass,
	})
	s.Assert().Error(err)
	stat, _ := status.FromError(err)
	s.Assert().Equal(codes.Unauthenticated, stat.Code())

	// login with wrong pass
	s.db.EXPECT().GetUserData(s.ctx, login).Return(uid, hashed, nil)
	_, err = s.serverAuth.Login(s.ctx, &pb.Credentials{
		Login:    login,
		Password: "wrong " + pass,
	})
	s.Assert().Error(err)
	stat, _ = status.FromError(err)
	s.Assert().Equal(codes.Unauthenticated, stat.Code())

	// login but session for user already present
	s.db.EXPECT().GetUserData(s.ctx, login).Return(uid, hashed, nil)
	s.db.EXPECT().AddSession(s.ctx, uid, gomock.Any()).Return(database.ErrConflict)
	_, err = s.serverAuth.Login(s.ctx, &pb.Credentials{
		Login:    login,
		Password: pass,
	})
	s.Assert().Error(err)
	stat, _ = status.FromError(err)
	s.Assert().Equal(codes.AlreadyExists, stat.Code())

	// login with AddSession returning any error
	s.db.EXPECT().GetUserData(s.ctx, login).Return(uid, hashed, nil)
	s.db.EXPECT().AddSession(s.ctx, uid, gomock.Any()).Return(dbErr)
	_, err = s.serverAuth.Login(s.ctx, &pb.Credentials{
		Login:    login,
		Password: pass,
	})
	s.Assert().Error(err)
	stat, _ = status.FromError(err)
	s.Assert().Equal(codes.Internal, stat.Code())
}

func (s *gRPCServerTestSuite) TestAuthKickOtherSession() {
	const (
		login = "login"
		pass  = "password"
	)
	hashed, _ := passwords.HashPassword(pass)

	// login with kick
	uid := uuid.NewString()
	var sid string
	s.db.EXPECT().GetUserData(s.ctx, login).Return(uid, hashed, nil)
	s.db.EXPECT().UpdateSession(s.ctx, uid, gomock.Any()).DoAndReturn(
		func(_ context.Context, _, in string) error {
			sid = in
			return nil
		},
	)
	loginRes, err := s.serverAuth.KickOtherSession(s.ctx, &pb.Credentials{
		Login:    login,
		Password: pass,
	})
	s.Assert().NoError(err)
	s.Assert().Equal(sid, loginRes.SID)

	// login not found
	s.db.EXPECT().GetUserData(s.ctx, login).Return(uid, hashed, database.ErrNotFound)
	loginRes, err = s.serverAuth.KickOtherSession(s.ctx, &pb.Credentials{
		Login:    login,
		Password: pass,
	})
	s.Assert().Error(err)
	stat, _ := status.FromError(err)
	s.Assert().Equal(codes.Unauthenticated, stat.Code())

	// login with wrong pass
	s.db.EXPECT().GetUserData(s.ctx, login).Return(uid, hashed, nil)
	_, err = s.serverAuth.KickOtherSession(s.ctx, &pb.Credentials{
		Login:    login,
		Password: "wrong " + pass,
	})
	s.Assert().Error(err)
	stat, _ = status.FromError(err)
	s.Assert().Equal(codes.Unauthenticated, stat.Code())

	// login with UpdateSession returning any error
	s.db.EXPECT().GetUserData(s.ctx, login).Return(uid, hashed, nil)
	s.db.EXPECT().UpdateSession(s.ctx, uid, gomock.Any()).Return(dbErr)
	_, err = s.serverAuth.KickOtherSession(s.ctx, &pb.Credentials{
		Login:    login,
		Password: pass,
	})
	s.Assert().Error(err)
	stat, _ = status.FromError(err)
	s.Assert().Equal(codes.Internal, stat.Code())
}

func (s *gRPCServerTestSuite) TestAuthenticator() {
	// setup
	grpcServer := s.serverKeep.(*ServerGRPC)
	info := &grpc.UnaryServerInfo{
		Server:     "Auth",
		FullMethod: "base.Auth/Login",
	}
	handler := grpc.UnaryHandler(func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	})
	sid := uuid.NewString()
	ctx := metadata.NewIncomingContext(s.ctx, metadata.Pairs(sidMetaKey, sid))

	// any call to Auth
	_, err := grpcServer.Authenticate(ctx, nil, info, handler)
	s.Assert().NoError(err)

	info = &grpc.UnaryServerInfo{
		Server:     "GophKeeper",
		FullMethod: "base.GophKeeper/Login",
	}

	// normal call to Keeper
	uid := uuid.NewString()
	s.db.EXPECT().RefreshSession(ctx, sid).Return(uid, true, nil)
	_, err = grpcServer.Authenticate(ctx, nil, info, handler)
	s.Assert().NoError(err)

	// db returned error
	s.db.EXPECT().RefreshSession(ctx, sid).Return(uid, true, dbErr)
	_, err = grpcServer.Authenticate(ctx, nil, info, handler)
	s.Assert().Error(err)
	stat, _ := status.FromError(err)
	s.Assert().Equal(codes.Internal, stat.Code())

	// db returned not ok
	s.db.EXPECT().RefreshSession(ctx, sid).Return(uid, false, nil)
	_, err = grpcServer.Authenticate(ctx, nil, info, handler)
	s.Assert().Error(err)
	stat, _ = status.FromError(err)
	s.Assert().Equal(codes.ResourceExhausted, stat.Code())

	// ctx metadata has no key
	ctx = metadata.NewIncomingContext(s.ctx, metadata.Pairs("k", "v"))
	_, err = grpcServer.Authenticate(ctx, nil, info, handler)
	s.Assert().Error(err)
	stat, _ = status.FromError(err)
	s.Assert().Equal(codes.Unauthenticated, stat.Code())

	// ctx has no metadata
	ctx = s.ctx
	_, err = grpcServer.Authenticate(ctx, nil, info, handler)
	s.Assert().Error(err)
	stat, _ = status.FromError(err)
	s.Assert().Equal(codes.Unauthenticated, stat.Code())
}

func (s *gRPCServerTestSuite) TestGetCategoryHead() {
	uid := uuid.NewString()
	ctx := context.WithValue(s.ctx, uidMetaKey, uid)
	s.db.EXPECT().GetCredentialsHead(ctx, uid).Return(database.CategoryHead{
		{ID: "1", Metadata: "meta"},
	}, nil)
	head, err := s.serverKeep.GetCategoryHead(ctx, &pb.CategoryType_DTO{
		Category: pb.Category_CATEGORY_CRED,
	})
	s.Assert().NoError(err)
	s.Assert().NotNil(head)

	s.db.EXPECT().GetTextHead(ctx, uid).Return(database.CategoryHead{
		{ID: "1", Metadata: "meta"},
	}, nil)
	head, err = s.serverKeep.GetCategoryHead(ctx, &pb.CategoryType_DTO{
		Category: pb.Category_CATEGORY_TEXT,
	})
	s.Assert().NoError(err)
	s.Assert().NotNil(head)

	s.db.EXPECT().GetCardHead(ctx, uid).Return(database.CategoryHead{
		{ID: "1", Metadata: "meta"},
	}, nil)
	head, err = s.serverKeep.GetCategoryHead(ctx, &pb.CategoryType_DTO{
		Category: pb.Category_CATEGORY_CARD,
	})
	s.Assert().NoError(err)
	s.Assert().NotNil(head)

	// nonexistent category
	_, err = s.serverKeep.GetCategoryHead(ctx, &pb.CategoryType_DTO{
		Category: 100,
	})
	s.Assert().Error(err)
	stat, _ := status.FromError(err)
	s.Assert().Equal(codes.Internal, stat.Code())

	// nonexistent category
	s.db.EXPECT().GetCardHead(ctx, uid).Return(nil, dbErr)
	_, err = s.serverKeep.GetCategoryHead(ctx, &pb.CategoryType_DTO{
		Category: pb.Category_CATEGORY_CARD,
	})
	s.Assert().Error(err)
	stat, _ = status.FromError(err)
	s.Assert().Equal(codes.Internal, stat.Code())
}

func (s *gRPCServerTestSuite) TestStoreCredentials() {
	uid := uuid.NewString()
	ctx := context.WithValue(s.ctx, uidMetaKey, uid)
	const meta = "metadata"
	data := []byte("some data")
	var dataID string
	s.db.EXPECT().AddCredentials(ctx, uid, gomock.Any(), meta, data).DoAndReturn(
		func(_ context.Context, _, in, _ string, _ []byte) error {
			dataID = in
			return nil
		},
	)
	resp, err := s.serverKeep.StoreCredData(ctx, &pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	})
	s.Assert().NoError(err)
	s.Assert().Equal(dataID, resp.ID)

	s.db.EXPECT().AddCredentials(ctx, uid, gomock.Any(), meta, data).Return(dbErr)
	resp, err = s.serverKeep.StoreCredData(ctx, &pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	})
	s.Assert().Error(err)
	stat, _ := status.FromError(err)
	s.Assert().Equal(codes.Internal, stat.Code())
	s.Assert().Nil(resp)
}

func (s *gRPCServerTestSuite) TestLoadCredentials() {
	uid := uuid.NewString()
	ctx := context.WithValue(s.ctx, uidMetaKey, uid)
	const meta = "metadata"
	data := []byte("some data")
	dataID := uuid.NewString()
	s.db.EXPECT().GetCredentials(ctx, uid, dataID).Return(meta, data, nil)
	resp, err := s.serverKeep.LoadCredData(ctx, &pb.DataID_DTO{ID: dataID})
	s.Assert().NoError(err)
	s.Assert().Equal(data, resp.Data)
	s.Assert().Equal(meta, resp.Metadata)

	s.db.EXPECT().GetCredentials(ctx, uid, dataID).Return("", nil, dbErr)
	resp, err = s.serverKeep.LoadCredData(ctx, &pb.DataID_DTO{ID: dataID})
	s.Assert().Error(err)
	stat, _ := status.FromError(err)
	s.Assert().Equal(codes.Internal, stat.Code())
	s.Assert().Nil(resp)
}

func (s *gRPCServerTestSuite) TestStoreTextData() {
	uid := uuid.NewString()
	ctx := context.WithValue(s.ctx, uidMetaKey, uid)
	const meta = "metadata"
	data := []byte("some data")
	var dataID string
	s.db.EXPECT().AddText(ctx, uid, gomock.Any(), meta, data).DoAndReturn(
		func(_ context.Context, _, in, _ string, _ []byte) error {
			dataID = in
			return nil
		},
	)
	resp, err := s.serverKeep.StoreTextData(ctx, &pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	})
	s.Assert().NoError(err)
	s.Assert().Equal(dataID, resp.ID)

	s.db.EXPECT().AddText(ctx, uid, gomock.Any(), meta, data).Return(dbErr)
	resp, err = s.serverKeep.StoreTextData(ctx, &pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	})
	s.Assert().Error(err)
	stat, _ := status.FromError(err)
	s.Assert().Equal(codes.Internal, stat.Code())
	s.Assert().Nil(resp)
}

func (s *gRPCServerTestSuite) TestLoadTextData() {
	uid := uuid.NewString()
	ctx := context.WithValue(s.ctx, uidMetaKey, uid)
	const meta = "metadata"
	data := []byte("some data")
	dataID := uuid.NewString()
	s.db.EXPECT().GetText(ctx, uid, dataID).Return(meta, data, nil)
	resp, err := s.serverKeep.LoadTextData(ctx, &pb.DataID_DTO{ID: dataID})
	s.Assert().NoError(err)
	s.Assert().Equal(data, resp.Data)
	s.Assert().Equal(meta, resp.Metadata)

	s.db.EXPECT().GetText(ctx, uid, dataID).Return("", nil, dbErr)
	resp, err = s.serverKeep.LoadTextData(ctx, &pb.DataID_DTO{ID: dataID})
	s.Assert().Error(err)
	stat, _ := status.FromError(err)
	s.Assert().Equal(codes.Internal, stat.Code())
	s.Assert().Nil(resp)
}

func (s *gRPCServerTestSuite) TestStoreCreditCard() {
	uid := uuid.NewString()
	ctx := context.WithValue(s.ctx, uidMetaKey, uid)
	const meta = "metadata"
	data := []byte("some data")
	var dataID string
	s.db.EXPECT().AddCard(ctx, uid, gomock.Any(), meta, data).DoAndReturn(
		func(_ context.Context, _, in, _ string, _ []byte) error {
			dataID = in
			return nil
		},
	)
	resp, err := s.serverKeep.StoreCardData(ctx, &pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	})
	s.Assert().NoError(err)
	s.Assert().Equal(dataID, resp.ID)

	s.db.EXPECT().AddCard(ctx, uid, gomock.Any(), meta, data).Return(dbErr)
	resp, err = s.serverKeep.StoreCardData(ctx, &pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	})
	s.Assert().Error(err)
	stat, _ := status.FromError(err)
	s.Assert().Equal(codes.Internal, stat.Code())
	s.Assert().Nil(resp)
}

func (s *gRPCServerTestSuite) TestLoadCreditCard() {
	uid := uuid.NewString()
	ctx := context.WithValue(s.ctx, uidMetaKey, uid)
	const meta = "metadata"
	data := []byte("some data")
	dataID := uuid.NewString()
	s.db.EXPECT().GetCard(ctx, uid, dataID).Return(meta, data, nil)
	resp, err := s.serverKeep.LoadCardData(ctx, &pb.DataID_DTO{ID: dataID})
	s.Assert().NoError(err)
	s.Assert().Equal(data, resp.Data)
	s.Assert().Equal(meta, resp.Metadata)

	s.db.EXPECT().GetCard(ctx, uid, dataID).Return("", nil, dbErr)
	resp, err = s.serverKeep.LoadCardData(ctx, &pb.DataID_DTO{ID: dataID})
	s.Assert().Error(err)
	stat, _ := status.FromError(err)
	s.Assert().Equal(codes.Internal, stat.Code())
	s.Assert().Nil(resp)
}

func (s *gRPCServerTestSuite) TestStoreFileData() {
	uid := uuid.NewString()
	ctx := context.WithValue(s.ctx, uidMetaKey, uid)
	const meta = "metadata"
	data := []byte("some data")
	var dataID string
	s.db.EXPECT().AddFile(ctx, uid, gomock.Any(), meta, data).DoAndReturn(
		func(_ context.Context, _, in, _ string, _ []byte) error {
			dataID = in
			return nil
		},
	)
	resp, err := s.serverKeep.StoreFileData(ctx, &pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	})
	s.Assert().NoError(err)
	s.Assert().Equal(dataID, resp.ID)

	s.db.EXPECT().AddFile(ctx, uid, gomock.Any(), meta, data).Return(dbErr)
	resp, err = s.serverKeep.StoreFileData(ctx, &pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	})
	s.Assert().Error(err)
	stat, _ := status.FromError(err)
	s.Assert().Equal(codes.Internal, stat.Code())
	s.Assert().Nil(resp)
}

func (s *gRPCServerTestSuite) TestLoadFileData() {
	uid := uuid.NewString()
	ctx := context.WithValue(s.ctx, uidMetaKey, uid)
	const meta = "metadata"
	data := []byte("some data")
	dataID := uuid.NewString()
	s.db.EXPECT().GetFile(ctx, uid, dataID).Return(meta, data, nil)
	resp, err := s.serverKeep.LoadFileData(ctx, &pb.DataID_DTO{ID: dataID})
	s.Assert().NoError(err)
	s.Assert().Equal(data, resp.Data)
	s.Assert().Equal(meta, resp.Metadata)

	s.db.EXPECT().GetFile(ctx, uid, dataID).Return("", nil, dbErr)
	resp, err = s.serverKeep.LoadFileData(ctx, &pb.DataID_DTO{ID: dataID})
	s.Assert().Error(err)
	stat, _ := status.FromError(err)
	s.Assert().Equal(codes.Internal, stat.Code())
	s.Assert().Nil(resp)
}

func TestGRPCServerUnit(t *testing.T) {
	suite.Run(t, new(gRPCServerTestSuite))
}
