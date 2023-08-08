package states

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"gophKeeper/src/client/cfg"
	"gophKeeper/src/client/encoding"
	GRPCClient "gophKeeper/src/client/grpcClient"
	"os"
)

const (
	stateLogin = iota
	stateMenu
	stateConfig
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

type stateData interface {
	state
	Show(ctx context.Context)
}
type Application struct {
	state     state
	cat       map[int]state
	config    *cfg.ConfigT
	grpc      *GRPCClient.GRPCClient
	closeFunc func() error
	userKey   string
	encoder   *encoding.Encoder
}

func New(config *cfg.ConfigT) *Application {
	app := &Application{
		config: config,
	}
	catalogue := map[int]state{
		stateLogin:  &stateLoginType{app, config},
		stateMenu:   &stateMenuType{app, config},
		stateConfig: &stateConfigType{app, config},
		stateCreds:  &stateCredsType{app, config, nil},
		//stateCard:   &stateCardType{app, config, nil},
		//stateText:   &stateTextType{app, config, nil},
		//stateFile:   &stateFileType{app, config, nil},
	}
	app.state = catalogue[stateLogin]
	app.cat = catalogue
	app.grpc, app.closeFunc = GRPCClient.New(config)
	return app
}

func (a *Application) Run() {
	const colorRed = "\033[0;31m"
	const colorNone = "\033[0m"

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
			fmt.Println(colorRed, err, colorNone)
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

var ErrNoKeyStored = errors.New("no secret found in storage")
var ErrNoFile = errors.New("no file was detected")

func (a *Application) loadSecret() error {
	secretDecoder := encoding.New(a.userKey)

	file, err := os.Open(a.config.SecretPath)
	if err != nil {
		return ErrNoFile
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var newKey string
	for scanner.Scan() {
		line := scanner.Bytes()
		newKey, err = secretDecoder.Decode(line)
		if err == nil {
			break
		}
	}
	a.encoder = encoding.New(newKey)
	return nil
}
func (a *Application) saveSecret(secret string) error {
	secretDecoder := encoding.New(a.userKey)

	file, err := os.OpenFile(a.config.SecretPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return ErrNoFile
	}
	defer file.Close()

	fmt.Printf("resulting secret: % x\n", secretDecoder.Encode(secret))
	_, err2 := file.Write(secretDecoder.Encode(secret))

	if err2 != nil {
		fmt.Println("Could not write text to example.txt")
	}

	a.encoder = encoding.New(secret)
	return nil
}

var (
	commandHelp = []string{"help", "h"}
	commandExit = []string{"exit", "x"}
	commandBack = []string{"back", "b"}
	commandOpen = []string{"open", "o"}
	commandLoad = []string{"load", "l"}
)
