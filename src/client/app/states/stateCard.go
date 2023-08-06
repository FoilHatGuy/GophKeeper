package states

import "gophKeeper/src/client/cfg"

type stateCardType struct {
	app    *Application
	config *cfg.ConfigT
}

func (s *stateCardType) Execute(command string) (resultState state, err error) {

	return s, err
}
