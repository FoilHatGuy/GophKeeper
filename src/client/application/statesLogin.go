package application

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"

	"gophKeeper/src/client/cfg"
	GRPCClient "gophKeeper/src/client/grpcClient"
)

type stateLoginType struct {
	stateGetName
	app    *Application
	config *cfg.ConfigT
}

func newLoginState(app *Application, config *cfg.ConfigT) state {
	return &stateLoginType{
		app:          app,
		config:       config,
		stateGetName: stateGetName{stateName: "Login view"},
	}
}

var (
	commandLogin    = []string{"login", "l"}
	commandRegister = []string{"register", "reg", "r"}
)

func (s *stateLoginType) execute(ctx context.Context, command string) (resultState state, err error) {
	arguments := strings.Split(command, " ")
	switch {
	case includes(commandHelp, strings.ToLower(arguments[0])):
		fmt.Printf("This is Login screen. Available commands:\n"+
			"help     - %q - shows available commands (this screen)\n"+
			"login    - %q $login$ $password$ - tries to log in with entered credentials\n"+
			"register - %q $login$ $password$ - registers on the server with entered credentials\n",
			commandHelp, commandLogin, commandRegister)
		return s, nil

	case includes(commandLogin, strings.ToLower(arguments[0])):
		if len(arguments) != 3 {
			return s, ErrUnrecognizedCommand
		}
		return s.login(ctx, arguments[1], arguments[2])

	case includes(commandRegister, strings.ToLower(arguments[0])):
		if len(arguments) != 3 {
			return s, ErrUnrecognizedCommand
		}
		err = s.app.grpc.Register(ctx, arguments[1], arguments[2])
		if err != nil {
			return s, fmt.Errorf("error occured during Register attempt. details: %w", err)
		}
		return s.login(ctx, arguments[1], arguments[2])
	}
	return s, ErrUnrecognizedCommand
}

func (s *stateLoginType) login(ctx context.Context, login, password string) (state, error) {
	modalKick := s.app.newModal(
		"Kick other session?",
		"This user is already logged in.\n"+
			"You can kick other device or don't do anything.\n"+
			"Kick other session? [y/n]",

		func(c context.Context, command string) (state, error) {
			switch command {
			case "y":
				err := s.app.grpc.KickOtherSession(c, login, password)
				if err != nil {
					return nil, fmt.Errorf("error occured during kicking session. details: %w", err)
				}
				return s.saveSecret(login, password)
			case "n":
				return s.app.cat[stateLogin], nil
			}
			return nil, nil
		},
	)

	err := s.app.grpc.Login(ctx, login, password)
	if errors.Is(err, GRPCClient.ErrAlreadyLoggedIn) {
		modalKick.Print()
		return modalKick, nil
	}
	if err != nil {
		return s, fmt.Errorf("error occured during Login attempt. details: %w", err)
	}
	return s.saveSecret(login, password)
}

func (s *stateLoginType) saveSecret(login, password string) (outState state, err error) {
	keyLength := 5
	modalSecret := s.app.newModal(
		"Enter your secret",
		"Secret key for this user is not saved on the device.\n"+
			"Please enter a secret which will be used to encode your data:",

		func(c context.Context, command string) (state, error) {
			if len(command) < keyLength {
				return nil, fmt.Errorf("please use a key at least %d long", keyLength)
			}
			err = s.app.saveSecret(command)
			if err != nil {
				return nil, fmt.Errorf("error occured during saving secret. details: %w", err)
			}
			return s.app.cat[stateMenu], nil
		},
	)

	key := strings.Join([]string{
		"USER_KEY",
		login,
		password,
	}, ":")
	arr := sha256.Sum256([]byte(key))
	s.app.userKey = string(arr[:])

	if s.app.encoder == nil {
		err = s.app.loadSecret()
		if err != nil {
			modalSecret.Print()
			return modalSecret, nil
		}
	}

	return s.app.cat[stateMenu], nil
}
