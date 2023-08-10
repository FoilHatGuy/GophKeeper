package application

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"

	"gophKeeper/src/client/cfg"
	"gophKeeper/src/client/encoding"
	GRPCClient "gophKeeper/src/client/grpcClient"
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
	execute(ctx context.Context, command string) (resultState state, err error)
	getName() string
}

type stateData interface {
	state
	show(ctx context.Context)
	fetch(ctx context.Context) (err error)
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

func newApplication(config *cfg.ConfigT) *Application {
	app := &Application{
		config: config,
	}
	catalogue := map[int]state{
		stateLogin:  newLoginState(app, config),
		stateMenu:   newMenuState(app, config),
		stateConfig: newConfigState(app, config),
		stateCreds:  newCredsState(app, config),
		// stateCard:   &stateCardType{app, config, nil},
		// stateText:   &stateTextType{app, config, nil},
		// stateFile:   &stateFileType{app, config, nil},
	}
	app.state = catalogue[stateLogin]
	app.cat = catalogue
	app.grpc, app.closeFunc = GRPCClient.New(config)
	return app
}

func (a *Application) Run() {
	const colorRed = "\033[0;31m"
	const colorBlue = "\033[0;34m"
	const colorNone = "\033[0m"

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("app is ready to accept commands!")
	for {
		fmt.Printf("%s%s:%s ", colorBlue, a.state.getName(), colorNone)
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

func New(config *cfg.ConfigT) {
	app := newApplication(config)
	app.Run()
}

func (a *Application) Execute(ctx context.Context, command string) error {
	if includes(commandExit, command) {
		return ErrExit
	}

	if includes(commandPing, command) {
		err := a.grpc.Ping(ctx)
		if err != nil {
			return fmt.Errorf("ping unsucessful: %w", err)
		}
	}

	newState, err := a.state.execute(ctx, command)
	if errors.Is(err, ErrUnrecognizedCommand) {
		newState, err = a.state.execute(ctx, commandHelp[0])
		if err != nil {
			return fmt.Errorf("application error: %w", err)
		}
	}
	a.state = newState
	if err != nil {
		return fmt.Errorf("application error: %w", err)
	}
	return nil
}

var (
	ErrNoKeyStored = errors.New("no secret found in storage")
	ErrNoFile      = errors.New("no file was detected")
)

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
	if newKey == "" {
		return ErrNoKeyStored
	}

	a.encoder = encoding.New(newKey)
	return nil
}

func (a *Application) saveSecret(secret string) error {
	secretDecoder := encoding.New(a.userKey)

	file, err := os.OpenFile(a.config.SecretPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		return ErrNoFile
	}
	defer file.Close()

	fmt.Printf("resulting secret: % x\n", secretDecoder.Encode(secret))
	_, err2 := file.Write(secretDecoder.Encode(secret))

	if err2 != nil {
		return errors.New("wile is not writeable")
	}

	a.encoder = encoding.New(secret)
	return nil
}

var (
	commandHelp = []string{"help", "?"}
	commandPing = []string{"ping", "p"}
	commandExit = []string{"exit", "x"}
	commandBack = []string{"back", "b"}
	commandOpen = []string{"open", "o"}
	commandLoad = []string{"load", "l"}
)
