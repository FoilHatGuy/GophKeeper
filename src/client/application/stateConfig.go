package application

import (
	"context"
	"fmt"
	"strings"

	"gophKeeper/src/client/cfg"
)

type stateConfigType struct {
	stateName string
	app       *Application
	config    *cfg.ConfigT
}

func newConfigState(app *Application, config *cfg.ConfigT) state {
	return &stateConfigType{
		app:       app,
		config:    config,
		stateName: "Config view",
	}
}

func (s *stateConfigType) getName() string {
	return s.stateName
}

// var (
// commandCred   = []string{"credentials", "cred", "cr"}
// commandCard   = []string{"card", "c"}
// commandText   = []string{"text", "t"}
// commandFile   = []string{"file", "f"}
// )

func (s *stateConfigType) execute(ctx context.Context, command string) (resultState state, err error) {
	arguments := strings.Split(command, " ")
	switch {
	case includes(commandHelp, strings.ToLower(arguments[0])):
		fmt.Printf("This is Login screen. Available commands:\n"+
			"help   - %q - shows available commands (this screen)\n"+
			"back   - %q - return to menu:\n"+
			"config - %q - opens config page",
			commandHelp,
			commandBack,
			commandConfig)

	case includes(commandBack, strings.ToLower(arguments[0])):
		return s.app.cat[stateMenu], nil

	case includes(commandOpen, strings.ToLower(arguments[0])):
		switch {
		case includes(commandCred, strings.ToLower(arguments[1])):
			return s.app.cat[stateCreds], nil
		case includes(commandCard, strings.ToLower(arguments[1])):
			return s.app.cat[stateCard], nil
		case includes(commandText, strings.ToLower(arguments[1])):
			return s.app.cat[stateText], nil
		case includes(commandFile, strings.ToLower(arguments[1])):
			return s.app.cat[stateFile], nil
		}

	default:
		return s, ErrUnrecognizedCommand
	}

	if err != nil {
		return s, err
	}
	return s, err
}
