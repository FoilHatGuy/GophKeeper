package states

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"gophKeeper/src/client/cfg"
	GRPCClient "gophKeeper/src/client/grpcClient"
	"os"
)

const (
	stateLogin = iota
	stateMenu
	stateCreds
	stateCard
	stateText
	stateFile
)

var (
	ErrExit                = errors.New("application exit")
	ErrUnrecognizedCommand = errors.New("incorrect command")
)

type state interface {
	Execute(ctx context.Context, command string) (resultState state, err error)
}
type Application struct {
	state     state
	cat       map[int]state
	config    *cfg.ConfigT
	grpc      *GRPCClient.GRPCClient
	closeFunc func() error
}

func New(config *cfg.ConfigT) *Application {
	app := &Application{
		config: config,
	}
	catalogue := map[int]state{
		stateLogin: &stateLoginType{app, config},
		stateMenu:  &stateMenuType{app, config},
		stateCreds: &stateCredsType{app, config},
		stateCard:  &stateCardType{app, config},
		stateText:  &stateTextType{app, config},
		stateFile:  &stateFileType{app, config},
	}
	app.state = catalogue[stateLogin]
	app.cat = catalogue
	app.grpc, app.closeFunc = GRPCClient.New(config)
	return app
}

func (a *Application) Run() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("app is ready to accept commands!")
	for {
		fmt.Printf("%T: ", a.state)
		scanner.Scan()
		ctx := context.Background()
		err := a.Execute(ctx, scanner.Text())
		if errors.Is(err, ErrExit) {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
	}
	err := a.closeFunc()
	panic(err)
}

func includes(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (a *Application) Execute(ctx context.Context, command string) error {
	if includes(commandExit, command) {
		return ErrExit
	}

	newState, err := a.state.Execute(ctx, command)
	if errors.Is(err, ErrUnrecognizedCommand) {
		_, err := a.state.Execute(ctx, commandHelp[0])
		if err != nil {
			return err
		}
	}
	a.state = newState
	if err != nil {
		return err
	}
	return nil
}

var (
	commandHelp        = []string{"help", "h"}
	commandLogin       = []string{"login", "l"}
	commandRegister    = []string{"register", "reg", "r"}
	commandExit        = []string{"exit", "x"}
	commandBack        = []string{"back", "b"}
	commandCredentials = []string{"credentials", "cred", "cr"}
	commandCard        = []string{"card"}
	commandText        = []string{"text"}
	commandFile        = []string{"file"}
	commandLoad        = []string{"load"}
	//commandLoad = "load"
)
