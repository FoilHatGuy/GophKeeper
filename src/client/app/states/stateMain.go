package states

import "gophKeeper/src/client/cfg"

type stateMenuType struct {
	app    *Application
	config *cfg.ConfigT
}

func (s *stateMenuType) Execute(command string) (resultState state, err error) {

	return s, err
}
