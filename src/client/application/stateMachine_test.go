//go:build unit

package application

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"gophKeeper/src/client/cfg"
	"gophKeeper/src/client/encoding"
	grpcclient "gophKeeper/src/client/grpcClient"
)

var errCall = errors.New("error in grpc")

type StateMachineTestSuite struct {
	suite.Suite
	config  *cfg.ConfigT
	ctx     context.Context
	secrets []string
}

func (s *StateMachineTestSuite) SetupSuite() {
	testConf := "./test_config_delete_later"
	s.secrets = append(s.secrets, testConf)
	s.T().Setenv("GKEEPER_CONFIG", testConf)
	s.config = cfg.New(cfg.FromDefaults())
	s.config.SecretPath = "./TEST_FILE_DELETE_LATER"
	s.ctx = context.Background()
}

func (s *StateMachineTestSuite) TearDownSuite() {
	for _, fName := range s.secrets {
		os.Remove(fName)
	}
	os.Remove(s.config.SecretPath)
}

func (s *StateMachineTestSuite) TestLoginStateLogin() {
	grpc := NewMockGRPCWrapper(gomock.NewController(s.T()))
	const (
		login = "login"
		pass  = "pass"
	)
	config := cfg.New(cfg.FromDefaults())
	config.SecretPath = "./TestLoginStateLogin"
	s.secrets = append(s.secrets, config.SecretPath)
	app := newApplication(config, grpc, func() error { return nil })
	testState := newLoginState(app, s.config)

	// wrong login command
	resState, err := testState.execute(s.ctx, "login")
	s.Assert().ErrorIs(err, ErrUnrecognizedCommand)
	s.Assert().IsType(&stateLoginType{}, resState)

	// normal login
	grpc.EXPECT().Login(s.ctx, login, pass).Return(nil)
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

	// repeat login with present secret
	app.encoder = nil
	grpc.EXPECT().Login(s.ctx, login, pass).Return(nil)
	resState, err = testState.execute(s.ctx,
		fmt.Sprintf("l %s %s", login, pass),
	)
	s.Assert().NoError(err)
	s.Assert().IsType(&stateMenuType{}, resState)

	// login grpc error
	config.SecretPath = "./TestLoginStateLogin_2"
	s.secrets = append(s.secrets, config.SecretPath)
	grpc.EXPECT().Login(s.ctx, login, pass).Return(errCall)
	resState, err = testState.execute(s.ctx,
		fmt.Sprintf("l %s %s", login, pass),
	)
	s.Assert().Error(err)
	s.Assert().IsType(&stateLoginType{}, resState)

	// login grpc already logged in
	grpc.EXPECT().Login(s.ctx, login, pass).Return(grpcclient.ErrAlreadyLoggedIn)
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
	grpc.EXPECT().KickOtherSession(s.ctx, login, pass).Return(nil)
	resState, err = modal.execute(s.ctx, "y")
	s.Assert().NoError(err)
	s.Assert().IsType(&stateMenuType{}, resState)

	// error in grpc
	grpc.EXPECT().KickOtherSession(s.ctx, login, pass).Return(errCall)
	resState, err = modal.execute(s.ctx, "y")
	s.Assert().Error(err)
	s.Assert().IsType(&stateModal{}, resState)

	// unexpected command
	resState, err = modal.execute(s.ctx, "bruh")
	s.Assert().NoError(err)
	s.Assert().IsType(&stateModal{}, resState)
}

func (s *StateMachineTestSuite) TestLoginStateOther() {
	grpc := NewMockGRPCWrapper(gomock.NewController(s.T()))
	const (
		login = "login"
		pass  = "pass"
	)

	config := cfg.New(cfg.FromDefaults())
	config.SecretPath = "./TestLoginStateOther"
	s.secrets = append(s.secrets, config.SecretPath)
	app := newApplication(config, grpc, func() error { return nil })
	testState := newLoginState(app, s.config)

	// normal register
	grpc.EXPECT().Register(s.ctx, login, pass).Return(nil)
	grpc.EXPECT().Login(s.ctx, login, pass).Return(nil)
	resState, err := testState.execute(s.ctx,
		fmt.Sprintf("r %s %s", login, pass),
	)
	s.Assert().NoError(err)
	s.Assert().IsType(&stateModal{}, resState)

	// register argument count
	resState, err = testState.execute(s.ctx, "reg")
	s.Assert().Error(err)
	s.Assert().IsType(&stateLoginType{}, resState)

	// normal register
	grpc.EXPECT().Register(s.ctx, login, pass).Return(errCall)
	resState, err = testState.execute(s.ctx,
		fmt.Sprintf("r %s %s", login, pass),
	)
	s.Assert().Error(err)
	s.Assert().IsType(&stateLoginType{}, resState)

	resState, err = testState.execute(s.ctx, "?")
	s.Assert().NoError(err)
	s.Assert().IsType(&stateLoginType{}, resState)

	resState, err = testState.execute(s.ctx, "wrong")
	s.Assert().ErrorIs(err, ErrUnrecognizedCommand)
	s.Assert().IsType(&stateLoginType{}, resState)
}

func (s *StateMachineTestSuite) TestMenuState() {
	grpc := NewMockGRPCWrapper(gomock.NewController(s.T()))
	config := cfg.New(cfg.FromDefaults())
	config.SecretPath = "./TestMenuState"
	s.secrets = append(s.secrets, config.SecretPath)
	app := newApplication(config, grpc, func() error { return nil })
	testState := newMenuState(app, s.config)
	app.encoder = encoding.New("secret")

	// enter config
	resState, err := testState.execute(s.ctx, "cfg")
	s.Assert().NoError(err)
	s.Assert().IsType(&stateConfigType{}, resState)

	// grpc err
	grpc.EXPECT().GetCategoryHead(s.ctx, grpcclient.CategoryCred).Return(nil, errCall)
	resState, err = testState.execute(s.ctx, "open cr")
	s.Assert().Error(err)
	s.Assert().IsType(&stateDataType{}, resState)

	// enter Cred
	grpc.EXPECT().GetCategoryHead(s.ctx, grpcclient.CategoryCred).Return([]*grpcclient.CategoryEntry{
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
	grpc.EXPECT().GetCategoryHead(s.ctx, grpcclient.CategoryText).Return([]*grpcclient.CategoryEntry{
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
	grpc.EXPECT().GetCategoryHead(s.ctx, grpcclient.CategoryCard).Return([]*grpcclient.CategoryEntry{
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
	//	grpc.EXPECT().GetCategoryHead(s.ctx, grpcclient.CategoryFile).Return([]*grpcclient.CategoryEntry{
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
	s.Assert().NoError(err)
	s.Assert().IsType(&stateMenuType{}, resState)

	// unrecognized command
	resState, err = testState.execute(s.ctx, "wrong")
	s.Assert().ErrorIs(err, ErrUnrecognizedCommand)
	s.Assert().IsType(&stateMenuType{}, resState)
}

func (s *StateMachineTestSuite) TestDataState() {
	grpc := NewMockGRPCWrapper(gomock.NewController(s.T()))
	config := cfg.New(cfg.FromDefaults())
	config.SecretPath = "./TestDataState"
	s.secrets = append(s.secrets, config.SecretPath)
	app := newApplication(config, grpc, func() error { return nil })
	testState := newCredsState(app, s.config)
	app.encoder = encoding.New("secret")
	data := "login\x00pass"
	dataSlice := strings.Split(data, "\x00")
	meta := "metadata"
	dataID := uuid.NewString()

	// head
	grpc.EXPECT().GetCategoryHead(s.ctx, grpcclient.CategoryCred).Return([]*grpcclient.CategoryEntry{}, nil)
	_, err := testState.execute(s.ctx, "head")
	s.Assert().Equal("empty list", err.Error())

	grpc.EXPECT().GetCategoryHead(s.ctx, grpcclient.CategoryCred).Return([]*grpcclient.CategoryEntry{}, errCall)
	_, err = testState.execute(s.ctx, "head")
	s.Assert().Error(err)

	grpc.EXPECT().GetCategoryHead(s.ctx, grpcclient.CategoryCred).Return([]*grpcclient.CategoryEntry{
		{
			DataID:   dataID,
			Metadata: meta,
		},
	}, nil)
	_, err = testState.execute(s.ctx, "head")
	s.Assert().NoError(err)

	// add data
	grpc.EXPECT().StoreCredData(s.ctx, gomock.Any(), meta).DoAndReturn(
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
	grpc.EXPECT().StoreCredData(s.ctx, gomock.Any(), meta).Return("", "", errCall)
	_, err = testState.execute(s.ctx, "add")
	_, err = testState.execute(s.ctx, dataSlice[0])
	_, err = testState.execute(s.ctx, dataSlice[1])
	resState, err = testState.execute(s.ctx, meta)
	s.Assert().Error(err)
	s.Assert().IsType(&stateDataType{}, resState)

	// load
	grpc.EXPECT().LoadCredData(s.ctx, dataID).Return(app.encoder.Encode(data), nil)
	resState, err = testState.execute(s.ctx, "load 0")
	s.Assert().NoError(err)
	s.Assert().IsType(&stateDataType{}, resState)

	// help
	resState, err = testState.execute(s.ctx, "help")
	s.Assert().NoError(err)
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

func (s *StateMachineTestSuite) TestConfigState() {
	grpc := NewMockGRPCWrapper(gomock.NewController(s.T()))
	config := cfg.New(cfg.FromDefaults())
	config.SecretPath = "./TestConfigState"
	s.secrets = append(s.secrets, config.SecretPath)
	app := newApplication(config, grpc, func() error { return nil })
	testState := newConfigState(app, s.config)

	// list
	resState, err := testState.execute(s.ctx, "list")
	s.Assert().NoError(err)
	s.Assert().IsType(&stateConfigType{}, resState)

	// about
	resState, err = testState.execute(s.ctx, "about")
	s.Assert().NoError(err)
	s.Assert().IsType(&stateConfigType{}, resState)

	// modify
	// wrong num of arguments
	resState, err = testState.execute(s.ctx, "mod")
	s.Assert().ErrorIs(err, ErrUnrecognizedCommand)
	s.Assert().IsType(&stateConfigType{}, resState)

	newAddr := "localhost:1234"
	resState, err = testState.execute(s.ctx, "mod 1 "+newAddr)
	s.Assert().NoError(err)
	s.Assert().IsType(&stateConfigType{}, resState)
	s.Assert().Equal(newAddr, config.ServerAddress)

	newSecret := "./.keys"
	resState, err = testState.execute(s.ctx, "mod 2 "+newSecret)
	s.Assert().NoError(err)
	s.Assert().IsType(&stateConfigType{}, resState)
	s.Assert().Equal(newSecret, config.SecretPath)

	// help
	resState, err = testState.execute(s.ctx, "help")
	s.Assert().NoError(err)
	s.Assert().IsType(&stateConfigType{}, resState)

	// unrecognized command
	resState, err = testState.execute(s.ctx, "wrong")
	s.Assert().ErrorIs(err, ErrUnrecognizedCommand)
	s.Assert().IsType(&stateConfigType{}, resState)

	// back
	resState, err = testState.execute(s.ctx, "b")
	s.Assert().NoError(err)
	s.Assert().IsType(&stateMenuType{}, resState)
}

func (s *StateMachineTestSuite) TestApplication() {
	grpc := NewMockGRPCWrapper(gomock.NewController(s.T()))
	config := cfg.New(cfg.FromDefaults())
	config.SecretPath = "./TestApplication"
	s.secrets = append(s.secrets, config.SecretPath)
	app := newApplication(config, grpc, func() error { return nil })

	do := make(chan string)
	res := make(chan error)
	var err error
	go func() {
		for cmd := range do {
			if cmd == "" {
				break
			}
			res <- app.Execute(context.Background(), cmd)
		}
	}()
	s.Assert().IsType(&stateLoginType{}, app.state)

	grpc.EXPECT().Ping(gomock.Any()).Return(nil)
	fmt.Println("command: ping")
	do <- "ping"
	err = <-res
	s.Assert().NoError(err)

	grpc.EXPECT().Ping(gomock.Any()).Return(errCall)
	fmt.Println("command: ping")
	do <- "ping"
	err = <-res
	s.Assert().ErrorIs(err, errCall)

	fmt.Println("command: help")
	do <- "help"
	err = <-res
	s.Assert().NoError(err)

	grpc.EXPECT().Login(gomock.Any(), "login", "pass")
	fmt.Println("command: l login pass")
	do <- "l login pass"
	err = <-res
	s.Assert().NoError(err)

	fmt.Println("command: exit")
	do <- "exit"
	err = <-res
	s.Assert().ErrorIs(err, ErrExit)

	do <- ""
}

func TestStateMachine(t *testing.T) {
	suite.Run(t, new(StateMachineTestSuite))
}
