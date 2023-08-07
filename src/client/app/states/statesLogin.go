package states

import (
	"context"
	"errors"
	"fmt"
	"gophKeeper/src/client/cfg"
	GRPCClient "gophKeeper/src/client/grpcClient"
	"strings"
)

type stateLoginType struct {
	app    *Application
	config *cfg.ConfigT
}

func (s *stateLoginType) Execute(ctx context.Context, command string) (resultState state, err error) {
	arguments := strings.Split(command, " ")
	switch {
	case includes(commandHelp, strings.ToLower(arguments[0])):
		fmt.Printf("This is Login screen. Available commands:\n"+
			"help     - %q - shows available commands (this screen)\n"+
			"login    - %q $login$ $password$ - tries to log in with entered credentials\n"+
			"register - %q $login$ $password$ - registers on the server with entered credentials\n",
			commandHelp, commandLogin, commandRegister)

	case includes(commandLogin, strings.ToLower(arguments[0])):
		if len(arguments) != 3 {
			return s, ErrUnrecognizedCommand
		}
		return s.Login(ctx, arguments[1], arguments[2])

	case includes(commandRegister, strings.ToLower(arguments[0])):
		if len(arguments) != 3 {
			return s, ErrUnrecognizedCommand
		}
		err = s.app.grpc.Register(ctx, arguments[1], arguments[2])
		if err != nil {
			return s, fmt.Errorf("error occured during Register attempt. details: %w", err)
		}
		return s.Login(ctx, arguments[1], arguments[2])
	}
	if err != nil {
		return s, err
	}
	return s, err
}

func (s *stateLoginType) Login(ctx context.Context, login, password string) (state, error) {
	err := s.app.grpc.Login(ctx, login, password)
	if errors.Is(err, GRPCClient.ErrAlreadyLoggedIn) {
		return s.app.newModal(
			"This user is already logged in.\n"+
				"You can kick other device or don't do anything.\n"+
				"Kick other session? [y/n]",

			func(c context.Context) (state, error) {
				err := s.app.grpc.KickOtherSession(c, login, password)
				if err != nil {
					return nil, fmt.Errorf("error occured during kicking session. details: %w", err)
				}
				return s.app.cat[stateMenu], nil
			},

			func(c context.Context) (state, error) {
				return s, nil
			},
		), nil
	}
	if err != nil {
		return s, fmt.Errorf("error occured during Login attempt. details: %w", err)
	}
	return s.app.cat[stateMenu], nil
}
