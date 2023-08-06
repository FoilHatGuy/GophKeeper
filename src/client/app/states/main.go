package states

import (
	"errors"
	"gophKeeper/src/client/cfg"
)

const (
	stateLogin = iota
	stateRegister
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
	Execute(command string) (resultState state, err error)
}
type Application struct {
	state  state
	states map[int]state
}

func New(config *cfg.ConfigT) *Application {
	app := &Application{}
	catalogue := map[int]state{
		stateLogin:    &stateLoginType{app, config},
		stateRegister: &stateRegisterType{app, config},
		stateMenu:     &stateMenuType{app, config},
		stateCreds:    &stateCredsType{app, config},
		stateCard:     &stateCardType{app, config},
		stateText:     &stateTextType{app, config},
		stateFile:     &stateFileType{app, config},
	}
	app.state = catalogue[stateLogin]
	app.states = catalogue
	return app
}

func includes(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (s *Application) Execute(command string) error {
	if includes(commandExit, command) {
		return ErrExit
	}

	newState, err := s.state.Execute(command)
	if errors.Is(err, ErrUnrecognizedCommand) {
		_, err := s.state.Execute(commandHelp[0])
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	s.state = newState
	return nil
}

var (
	commandHelp        = []string{"help", "h"}
	commandLogin       = []string{"login", "l"}
	commandExit        = []string{"exit", "x"}
	commandBack        = []string{"back", "b"}
	commandCredentials = []string{"credentials", "cred", "cr"}
	commandCard        = []string{"card"}
	commandText        = []string{"text"}
	commandFile        = []string{"file"}
	commandLoad        = []string{"load"}
	//commandLoad = "load"
)
