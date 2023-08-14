//go:build unit

package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"gophKeeper/src/client/cfg"
	"gophKeeper/src/client/encoding"
	grpcclient "gophKeeper/src/client/grpcClient"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

var errCall = errors.New("error in grpc")

type StateMachineTestSuite struct {
	suite.Suite
	config *cfg.ConfigT
	ctx    context.Context
	grpc   *MockGRPCWrapper
}

func (s *StateMachineTestSuite) SetupSuite() {
	s.config = cfg.New(cfg.FromDefaults())
	s.config.SecretPath = "./TEST_FILE_DELETE_LATER"
	s.ctx = context.Background()
	s.grpc = NewMockGRPCWrapper(gomock.NewController(s.T()))
	os.Remove(s.config.SecretPath)
}

func (s *StateMachineTestSuite) TearDownSuite() {
	os.Remove(s.config.SecretPath)
}

func (s *StateMachineTestSuite) TestLoginStateLogin() {
	const (
		login = "login"
		pass  = "pass"
	)

	app := newApplication(s.config, s.grpc, func() error { return nil })
	testState := newLoginState(app, s.config)

	// wrong login command
	resState, err := testState.execute(s.ctx, "login")
	s.Assert().ErrorIs(err, ErrUnrecognizedCommand)
	s.Assert().IsType(&stateLoginType{}, resState)

	// normal login
	s.grpc.EXPECT().Login(s.ctx, login, pass).Return(nil)
	resState, err = testState.execute(s.ctx,
		fmt.Sprintf("l %s %s", login, pass),
	)
	s.Assert().NoError(err)
	s.Assert().IsType(&stateModal{}, resState)
	// test returned modal
	modal := resState.(*stateModal)
	s.Assert().Contains(modal.stateName, "secret")

	// modal secret too short
	resState, err = modal.execute(s.ctx, "shrt")
	s.Assert().Error(err)
	s.Assert().IsType(&stateModal{}, resState)

	// normal modal
	resState, err = modal.execute(s.ctx, "new_secret")
	s.Assert().NoError(err)
	s.Assert().IsType(&stateMenuType{}, resState)

	// login grpc error
	s.grpc.EXPECT().Login(s.ctx, login, pass).Return(errCall)
	resState, err = testState.execute(s.ctx,
		fmt.Sprintf("l %s %s", login, pass),
	)
	s.Assert().Error(err)
	s.Assert().IsType(&stateLoginType{}, resState)

	// login grpc already logged in
	s.grpc.EXPECT().Login(s.ctx, login, pass).Return(grpcclient.ErrAlreadyLoggedIn)
	resState, err = testState.execute(s.ctx,
		fmt.Sprintf("l %s %s", login, pass),
	)
	s.Assert().NoError(err)
	s.Assert().IsType(&stateModal{}, resState)
	// test returned modal
	modal = resState.(*stateModal)
	s.Assert().Contains(modal.stateName, "Kick")

	// test kick modal
	// return to login
	resState, err = modal.execute(s.ctx, "n")
	s.Assert().NoError(err)
	s.Assert().IsType(&stateLoginType{}, resState)

	// successful kick
	s.grpc.EXPECT().KickOtherSession(s.ctx, login, pass).Return(nil)
	resState, err = modal.execute(s.ctx, "y")
	s.Assert().NoError(err)
	s.Assert().IsType(&stateMenuType{}, resState)

	// error in grpc
	s.grpc.EXPECT().KickOtherSession(s.ctx, login, pass).Return(errCall)
	resState, err = modal.execute(s.ctx, "y")
	s.Assert().Error(err)
	s.Assert().IsType(&stateModal{}, resState)

	//unexpected command
	resState, err = modal.execute(s.ctx, "bruh")
	s.Assert().NoError(err)
	s.Assert().IsType(&stateModal{}, resState)
}

func (s *StateMachineTestSuite) TestLoginStateOther() {
	const (
		login = "login"
		pass  = "pass"
	)

	app := newApplication(s.config, s.grpc, func() error { return nil })
	testState := newLoginState(app, s.config)
	app.encoder = encoding.New("secret")

	// normal register
	s.grpc.EXPECT().Register(s.ctx, login, pass).Return(nil)
	s.grpc.EXPECT().Login(s.ctx, login, pass).Return(nil)
	resState, err := testState.execute(s.ctx,
		fmt.Sprintf("r %s %s", login, pass),
	)
	s.Assert().NoError(err)
	s.Assert().IsType(&stateMenuType{}, resState)

	// register argument count
	resState, err = testState.execute(s.ctx, "reg")
	s.Assert().Error(err)
	s.Assert().IsType(&stateLoginType{}, resState)

	// normal register
	s.grpc.EXPECT().Register(s.ctx, login, pass).Return(errCall)
	resState, err = testState.execute(s.ctx,
		fmt.Sprintf("r %s %s", login, pass),
	)
	s.Assert().Error(err)
	s.Assert().IsType(&stateLoginType{}, resState)

	resState, err = testState.execute(s.ctx, "?")
	s.Assert().ErrorIs(err, ErrUnrecognizedCommand)
	s.Assert().IsType(&stateLoginType{}, resState)

	resState, err = testState.execute(s.ctx, "wrong")
	s.Assert().ErrorIs(err, ErrUnrecognizedCommand)
	s.Assert().IsType(&stateLoginType{}, resState)
}

func (s *StateMachineTestSuite) TestMenuState() {
	app := newApplication(s.config, s.grpc, func() error { return nil })
	testState := newMenuState(app, s.config)
	app.encoder = encoding.New("secret")

	// enter config
	resState, err := testState.execute(s.ctx, "cfg")
	s.Assert().NoError(err)
	s.Assert().IsType(&stateConfigType{}, resState)

	// grpc err
	s.grpc.EXPECT().GetCategoryHead(s.ctx, grpcclient.CategoryCred).Return(nil, errCall)
	resState, err = testState.execute(s.ctx, "open cr")
	s.Assert().Error(err)
	s.Assert().IsType(&stateDataType{}, resState)

	// enter Cred
	s.grpc.EXPECT().GetCategoryHead(s.ctx, grpcclient.CategoryCred).Return([]*grpcclient.CategoryEntry{
		{
			DataID:   "1",
			Metadata: "2",
		},
	}, nil)
	resState, err = testState.execute(s.ctx, "open cr")
	s.Assert().NoError(err)
	s.Assert().IsType(&stateDataType{}, resState)
	newState := resState.(*stateDataType)
	s.Assert().Equal(grpcclient.CategoryCred, newState.category)

	// enter Text
	s.grpc.EXPECT().GetCategoryHead(s.ctx, grpcclient.CategoryText).Return([]*grpcclient.CategoryEntry{
		{
			DataID:   "1",
			Metadata: "2",
		},
	}, nil)
	resState, err = testState.execute(s.ctx, "open t")
	s.Assert().NoError(err)
	s.Assert().IsType(&stateDataType{}, resState)
	newState = resState.(*stateDataType)
	s.Assert().Equal(grpcclient.CategoryText, newState.category)

	// enter Card
	s.grpc.EXPECT().GetCategoryHead(s.ctx, grpcclient.CategoryCard).Return([]*grpcclient.CategoryEntry{
		{
			DataID:   "1",
			Metadata: "2",
		},
	}, nil)
	resState, err = testState.execute(s.ctx, "open ca")
	s.Assert().NoError(err)
	s.Assert().IsType(&stateDataType{}, resState)
	newState = resState.(*stateDataType)
	s.Assert().Equal(grpcclient.CategoryCard, newState.category)

	// enter File
	//	s.grpc.EXPECT().GetCategoryHead(s.ctx, grpcclient.CategoryFile).Return([]*grpcclient.CategoryEntry{
	//		{
	//			DataID:   "1",
	//			Metadata: "2",
	//		},
	//	}, nil)
	//	resState, err = testState.execute(s.ctx, "open f")
	//	s.Assert().NoError(err)
	//	s.Assert().IsType(&stateDataType{}, resState)
	//	newState = resState.(*stateDataType)
	//	s.Assert().Equal(grpcclient.CategoryFile, newState.category)

	// help
	resState, err = testState.execute(s.ctx, "?")
	s.Assert().ErrorIs(err, ErrUnrecognizedCommand)
	s.Assert().IsType(&stateMenuType{}, resState)

	// unrecognized command
	resState, err = testState.execute(s.ctx, "wrong")
	s.Assert().ErrorIs(err, ErrUnrecognizedCommand)
	s.Assert().IsType(&stateMenuType{}, resState)
}

func (s *StateMachineTestSuite) TestDataState() {
	app := newApplication(s.config, s.grpc, func() error { return nil })
	testState := newCredsState(app, s.config)
	app.encoder = encoding.New("secret")
	data := "login\x00pass"
	dataSlice := strings.Split(data, "\x00")
	meta := "metadata"
	dataID := uuid.NewString()

	// head
	s.grpc.EXPECT().GetCategoryHead(s.ctx, grpcclient.CategoryCred).Return([]*grpcclient.CategoryEntry{}, nil)
	_, err := testState.execute(s.ctx, "head")
	s.Assert().Equal("empty list", err.Error())

	s.grpc.EXPECT().GetCategoryHead(s.ctx, grpcclient.CategoryCred).Return([]*grpcclient.CategoryEntry{}, errCall)
	_, err = testState.execute(s.ctx, "head")
	s.Assert().Error(err)

	s.grpc.EXPECT().GetCategoryHead(s.ctx, grpcclient.CategoryCred).Return([]*grpcclient.CategoryEntry{
		{
			DataID:   dataID,
			Metadata: meta,
		},
	}, nil)
	_, err = testState.execute(s.ctx, "head")
	s.Assert().NoError(err)

	// add data
	s.grpc.EXPECT().StoreCredData(s.ctx, gomock.Any(), meta).DoAndReturn(
		func(_ context.Context, in []byte, _ string) (string, string, error) {

			return dataID, meta, nil
		})
	resState, err := testState.execute(s.ctx, "add")
	s.Assert().NoError(err)
	s.Assert().IsType(&stateDataType{}, resState)
	resState, err = resState.execute(s.ctx, dataSlice[0])
	s.Assert().NoError(err)
	s.Assert().IsType(&stateDataType{}, resState)
	resState, err = resState.execute(s.ctx, dataSlice[1])
	s.Assert().NoError(err)
	s.Assert().IsType(&stateDataType{}, resState)
	resState, err = resState.execute(s.ctx, meta)
	s.Assert().NoError(err)
	s.Assert().IsType(&stateDataType{}, resState)

	// err in call
	s.grpc.EXPECT().StoreCredData(s.ctx, gomock.Any(), meta).Return("", "", errCall)
	_, err = testState.execute(s.ctx, "add")
	_, err = testState.execute(s.ctx, dataSlice[0])
	_, err = testState.execute(s.ctx, dataSlice[1])
	resState, err = testState.execute(s.ctx, meta)
	s.Assert().Error(err)
	s.Assert().IsType(&stateDataType{}, resState)

	//load
	s.grpc.EXPECT().LoadCredData(s.ctx, dataID).Return(app.encoder.Encode(data), nil)
	resState, err = testState.execute(s.ctx, "load 0")
	s.Assert().NoError(err)
	s.Assert().IsType(&stateDataType{}, resState)

	// help
	resState, err = testState.execute(s.ctx, "help")
	s.Assert().ErrorIs(err, ErrUnrecognizedCommand)
	s.Assert().IsType(&stateDataType{}, resState)

	// unrecognized command
	resState, err = testState.execute(s.ctx, "wrong")
	s.Assert().ErrorIs(err, ErrUnrecognizedCommand)
	s.Assert().IsType(&stateDataType{}, resState)

	// back
	resState, err = testState.execute(s.ctx, "b")
	s.Assert().NoError(err)
	s.Assert().IsType(&stateMenuType{}, resState)
}

func TestStateMachine(t *testing.T) {
	suite.Run(t, new(StateMachineTestSuite))
}
