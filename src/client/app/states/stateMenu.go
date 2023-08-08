package states

import (
	"context"
	"fmt"
	"gophKeeper/src/client/cfg"
	"strings"
)

type stateMenuType struct {
	app    *Application
	config *cfg.ConfigT
}

var (
	commandCred   = []string{"credentials", "cred", "cr"}
	commandCard   = []string{"card", "c"}
	commandText   = []string{"text", "t"}
	commandFile   = []string{"file", "f"}
	commandConfig = []string{"config", "cfg", "c"}
)

func (s *stateMenuType) Execute(ctx context.Context, command string) (resultState state, err error) {

	arguments := strings.Split(command, " ")
	switch {
	case includes(commandHelp, strings.ToLower(arguments[0])):
		fmt.Printf("This is Login screen. Available commands:\n"+
			"help   - %q - shows available commands (this screen)\n"+
			"open   - %q $category$ - opens one of the categories:\n"+
			"\t%q - Card information\n"+
			"\t%q - Text information\n"+
			"\t%q - File information\n"+
			"\t%q - Credentials information\n"+
			"config - %q - opens config page\n",
			commandHelp,
			commandOpen, commandCred, commandCard, commandText, commandFile,
			commandConfig)

	case includes(commandOpen, strings.ToLower(arguments[0])):
		if len(arguments) != 2 {
			return s, ErrUnrecognizedCommand
		}
		switch {
		case includes(commandCred, strings.ToLower(arguments[1])):
			newState := s.app.cat[stateCreds].(stateData)
			newState.Show(ctx)
			return newState, nil
		case includes(commandCard, strings.ToLower(arguments[1])):
			newState := s.app.cat[stateCard].(stateData)
			newState.Show(ctx)
			return newState, nil
		case includes(commandText, strings.ToLower(arguments[1])):
			newState := s.app.cat[stateText].(stateData)
			newState.Show(ctx)
			return newState, nil
		case includes(commandFile, strings.ToLower(arguments[1])):
			newState := s.app.cat[stateFile].(stateData)
			newState.Show(ctx)
			return newState, nil
		}

	case includes(commandConfig, strings.ToLower(arguments[0])):
		return s.app.cat[stateConfig], nil

	default:
		return s, ErrUnrecognizedCommand
	}

	if err != nil {
		return s, err
	}
	return s, err
}
