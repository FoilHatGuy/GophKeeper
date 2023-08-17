//go:build unit

package grpcclient

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gophKeeper/src/client/cfg"
	pb "gophKeeper/src/pb"
	"testing"

	"github.com/stretchr/testify/suite"
)

type GRPCClientUnit struct {
	suite.Suite
	rpc  *GRPCClient
	auth *MockAuthClient
	keep *MockGophKeeperClient
}

var (
	errCall          = errors.New("rpc err")
	errAlreadyExists = status.Errorf(
		codes.AlreadyExists, "already exists: %s", errCall)
)

func (s *GRPCClientUnit) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.auth = NewMockAuthClient(ctrl)
	s.keep = NewMockGophKeeperClient(ctrl)
	s.rpc = &GRPCClient{
		config: cfg.New(cfg.FromDefaults()),
		auth:   s.auth,
		keep:   s.keep,
	}
}

func (s *GRPCClientUnit) TestNew() {
	s.Assert().Panics(func() {
		New(cfg.New())
	})
}

func (s *GRPCClientUnit) TestAuthenticate() {
	ctx := context.Background()
	method := "base.auth/Ping"
	flag := false
	err := s.rpc.Authenticate(ctx, method, "", "", &grpc.ClientConn{},
		func(_ context.Context, _ string,
			_, _ interface{}, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
			flag = true
			return nil
		})
	s.Assert().True(flag)
	s.Assert().NoError(err)

	flag = false
	err = s.rpc.Authenticate(ctx, method, "", "", &grpc.ClientConn{},
		func(_ context.Context, _ string,
			_, _ interface{}, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
			flag = true
			return errCall
		})
	s.Assert().True(flag)
	s.Assert().Error(err)

	method = "base.GophKeeper/Ping"
	flag = false
	err = s.rpc.Authenticate(ctx, method, "", "", &grpc.ClientConn{},
		func(_ context.Context, _ string,
			_, _ interface{}, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
			flag = true
			return nil
		})
	s.Assert().True(flag)
	s.Assert().NoError(err)

	flag = false
	err = s.rpc.Authenticate(ctx, method, "", "", &grpc.ClientConn{},
		func(_ context.Context, _ string,
			_, _ interface{}, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
			flag = true
			return errCall
		})
	s.Assert().True(flag)
	s.Assert().Error(err)
}

func (s *GRPCClientUnit) TestLogin() {
	ctx := context.Background()
	const lLogin = "login"
	const lPass = "pass"

	sid0 := uuid.NewString()
	s.auth.EXPECT().Login(ctx, &pb.Credentials{
		Login:    lLogin,
		Password: lPass,
	}).Return(&pb.SessionID_DTO{SID: sid0}, nil)
	err := s.rpc.Login(ctx, lLogin, lPass)
	s.Assert().NoError(err)
	s.Assert().Equal(sid0, s.rpc.sessionID)

	sid := uuid.NewString()
	s.auth.EXPECT().Login(ctx, &pb.Credentials{
		Login:    lLogin,
		Password: lPass,
	}).Return(&pb.SessionID_DTO{SID: sid}, errAlreadyExists)
	err = s.rpc.Login(ctx, lLogin, lPass)
	s.Assert().ErrorIs(err, ErrAlreadyLoggedIn)
	s.Assert().Equal(sid0, s.rpc.sessionID)

	sid = uuid.NewString()
	s.auth.EXPECT().Login(ctx, &pb.Credentials{
		Login:    lLogin,
		Password: lPass,
	}).Return(&pb.SessionID_DTO{SID: sid}, errCall)
	err = s.rpc.Login(ctx, lLogin, lPass)
	s.Assert().ErrorIs(err, errCall)
	s.Assert().Equal(sid0, s.rpc.sessionID)
}

func (s *GRPCClientUnit) TestKickOtherSession() {
	ctx := context.Background()
	const kLogin = "login"
	const kPass = "pass"

	sid0 := uuid.NewString()
	s.auth.EXPECT().KickOtherSession(ctx, &pb.Credentials{
		Login:    kLogin,
		Password: kPass,
	}).Return(&pb.SessionID_DTO{SID: sid0}, nil)
	err := s.rpc.KickOtherSession(ctx, kLogin, kPass)
	s.Assert().NoError(err)
	s.Assert().Equal(sid0, s.rpc.sessionID)

	sid1 := uuid.NewString()
	s.auth.EXPECT().KickOtherSession(ctx, &pb.Credentials{
		Login:    kLogin,
		Password: kPass,
	}).Return(&pb.SessionID_DTO{SID: sid1}, errCall)
	err = s.rpc.KickOtherSession(ctx, kLogin, kPass)
	s.Assert().ErrorIs(err, errCall)
	s.Assert().Equal(sid0, s.rpc.sessionID)
}

func (s *GRPCClientUnit) TestRegister() {
	ctx := context.Background()
	const rLogin = "login"
	const rPass = "pass"

	s.auth.EXPECT().Register(ctx, &pb.Credentials{
		Login:    rLogin,
		Password: rPass,
	}).Return(&pb.Empty{}, nil)
	err := s.rpc.Register(ctx, rLogin, rPass)
	s.Assert().NoError(err)

	s.auth.EXPECT().Register(ctx, &pb.Credentials{
		Login:    rLogin,
		Password: rPass,
	}).Return(&pb.Empty{}, errCall)
	err = s.rpc.Register(ctx, rLogin, rPass)
	s.Assert().ErrorIs(err, errCall)
}

func (s *GRPCClientUnit) TestPing() {
	ctx := context.Background()

	s.auth.EXPECT().Ping(ctx, gomock.Any()).Return(&pb.Empty{}, nil)
	err := s.rpc.Ping(ctx)
	s.Assert().NoError(err)

	s.auth.EXPECT().Ping(ctx, gomock.Any()).Return(&pb.Empty{}, errCall)
	err = s.rpc.Ping(ctx)
	s.Assert().ErrorIs(err, errCall)
}

func (s *GRPCClientUnit) TestGetCategoryHead() {
	ctx := context.Background()
	const dataID = "DataID"
	const metadata = "Metadata"

	s.keep.EXPECT().GetCategoryHead(ctx, gomock.Any()).Return(
		&pb.CategoryHead_DTO{
			Info: []*pb.DataInfo{
				{
					DataID:   dataID,
					Metadata: metadata,
				},
			},
		}, nil)
	head, err := s.rpc.GetCategoryHead(ctx, CategoryCred)
	s.Assert().NoError(err)
	s.Assert().Equal(dataID, head[0].DataID)
	s.Assert().Equal(metadata, head[0].Metadata)

	s.keep.EXPECT().GetCategoryHead(ctx, gomock.Any()).Return(
		&pb.CategoryHead_DTO{}, errCall)
	_, err = s.rpc.GetCategoryHead(ctx, CategoryCred)
	s.Assert().ErrorIs(err, errCall)
}

func (s *GRPCClientUnit) TestCredData() {
	ctx := context.Background()
	dataID := uuid.NewString()
	data := []byte(uuid.NewString())
	meta := uuid.NewString()

	// store
	s.keep.EXPECT().StoreCredData(ctx, &pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	}).Return(&pb.DataID_DTO{ID: dataID}, nil)
	newID, newMeta, err := s.rpc.StoreCredData(ctx, data, meta)
	s.Assert().NoError(err)
	s.Assert().Equal(dataID, newID)
	s.Assert().Equal(meta, newMeta)

	s.keep.EXPECT().StoreCredData(ctx, &pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	}).Return(nil, errCall)
	newID, newMeta, err = s.rpc.StoreCredData(ctx, data, meta)
	s.Assert().ErrorIs(err, errCall)

	// load
	s.keep.EXPECT().LoadCredData(ctx, &pb.DataID_DTO{
		ID: dataID,
	}).Return(&pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	}, nil)
	newData, err := s.rpc.LoadCredData(ctx, dataID)
	s.Assert().NoError(err)
	s.Assert().Equal(data, newData)

	s.keep.EXPECT().LoadCredData(ctx, &pb.DataID_DTO{
		ID: dataID,
	}).Return(nil, errCall)
	_, err = s.rpc.LoadCredData(ctx, dataID)
	s.Assert().ErrorIs(err, errCall)
}

func (s *GRPCClientUnit) TestTextData() {
	ctx := context.Background()
	dataID := uuid.NewString()
	data := []byte(uuid.NewString())
	meta := uuid.NewString()

	// store
	s.keep.EXPECT().StoreTextData(ctx, &pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	}).Return(&pb.DataID_DTO{ID: dataID}, nil)
	newID, newMeta, err := s.rpc.StoreTextData(ctx, data, meta)
	s.Assert().NoError(err)
	s.Assert().Equal(dataID, newID)
	s.Assert().Equal(meta, newMeta)

	s.keep.EXPECT().StoreTextData(ctx, &pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	}).Return(nil, errCall)
	newID, newMeta, err = s.rpc.StoreTextData(ctx, data, meta)
	s.Assert().ErrorIs(err, errCall)

	// load
	s.keep.EXPECT().LoadTextData(ctx, &pb.DataID_DTO{
		ID: dataID,
	}).Return(&pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	}, nil)
	newData, err := s.rpc.LoadTextData(ctx, dataID)
	s.Assert().NoError(err)
	s.Assert().Equal(data, newData)

	s.keep.EXPECT().LoadTextData(ctx, &pb.DataID_DTO{
		ID: dataID,
	}).Return(nil, errCall)
	_, err = s.rpc.LoadTextData(ctx, dataID)
	s.Assert().ErrorIs(err, errCall)
}

func (s *GRPCClientUnit) TestCardData() {
	ctx := context.Background()
	dataID := uuid.NewString()
	data := []byte(uuid.NewString())
	meta := uuid.NewString()

	// store
	s.keep.EXPECT().StoreCardData(ctx, &pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	}).Return(&pb.DataID_DTO{ID: dataID}, nil)
	newID, newMeta, err := s.rpc.StoreCardData(ctx, data, meta)
	s.Assert().NoError(err)
	s.Assert().Equal(dataID, newID)
	s.Assert().Equal(meta, newMeta)

	s.keep.EXPECT().StoreCardData(ctx, &pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	}).Return(nil, errCall)
	newID, newMeta, err = s.rpc.StoreCardData(ctx, data, meta)
	s.Assert().ErrorIs(err, errCall)

	// load
	s.keep.EXPECT().LoadCardData(ctx, &pb.DataID_DTO{
		ID: dataID,
	}).Return(&pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	}, nil)
	newData, err := s.rpc.LoadCardData(ctx, dataID)
	s.Assert().NoError(err)
	s.Assert().Equal(data, newData)

	s.keep.EXPECT().LoadCardData(ctx, &pb.DataID_DTO{
		ID: dataID,
	}).Return(nil, errCall)
	_, err = s.rpc.LoadCardData(ctx, dataID)
	s.Assert().ErrorIs(err, errCall)
}

func (s *GRPCClientUnit) TestFileData() {
	ctx := context.Background()
	dataID := uuid.NewString()
	data := []byte(uuid.NewString())
	meta := uuid.NewString()

	// store
	s.keep.EXPECT().StoreFileData(ctx, &pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	}).Return(&pb.DataID_DTO{ID: dataID}, nil)
	newID, newMeta, err := s.rpc.StoreFileData(ctx, data, meta)
	s.Assert().NoError(err)
	s.Assert().Equal(dataID, newID)
	s.Assert().Equal(meta, newMeta)

	s.keep.EXPECT().StoreFileData(ctx, &pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	}).Return(nil, errCall)
	newID, newMeta, err = s.rpc.StoreFileData(ctx, data, meta)
	s.Assert().ErrorIs(err, errCall)

	// load
	s.keep.EXPECT().LoadFileData(ctx, &pb.DataID_DTO{
		ID: dataID,
	}).Return(&pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	}, nil)
	newData, err := s.rpc.LoadFileData(ctx, dataID)
	s.Assert().NoError(err)
	s.Assert().Equal(data, newData)

	s.keep.EXPECT().LoadFileData(ctx, &pb.DataID_DTO{
		ID: dataID,
	}).Return(nil, errCall)
	_, err = s.rpc.LoadFileData(ctx, dataID)
	s.Assert().ErrorIs(err, errCall)
}

func TestGrPCClientUnit(t *testing.T) {
	suite.Run(t, new(GRPCClientUnit))
}
