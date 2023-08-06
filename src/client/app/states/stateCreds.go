package states

import "gophKeeper/src/client/cfg"

type stateCredsType struct {
	app    *Application
	config *cfg.ConfigT
}

func (s *stateCredsType) Execute(command string) (resultState state, err error) {

	return s, err
}
