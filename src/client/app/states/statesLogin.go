package states

import (
	"gophKeeper/src/client/cfg"
	"strings"
)

type stateLoginType struct {
	app    *Application
	config *cfg.ConfigT
}

func (s *stateLoginType) Execute(command string) (resultState state, err error) {
	arguments := strings.Split(command, " ")
	switch {
	case includes(commandLoad, arguments[0]):
		err = s.Login(arguments[1], arguments[2])
		if err != nil {
			return nil, err
		}
	}
	return s, err
}

func (s *stateLoginType) Login(login, password string) (err error) {

	return err
}

type stateRegisterType struct {
	app    *Application
	config *cfg.ConfigT
}

func (s *stateRegisterType) Execute(command string) (resultState state, err error) {

	return s, err
}
