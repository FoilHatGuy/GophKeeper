package states

import "gophKeeper/src/client/cfg"

type stateTextType struct {
	app    *Application
	config *cfg.ConfigT
}

func (s *stateTextType) Execute(command string) (resultState state, err error) {

	return s, err
}
