package states

import "gophKeeper/src/client/cfg"

type stateFileType struct {
	app    *Application
	config *cfg.ConfigT
}

func (s *stateFileType) Execute(command string) (resultState state, err error) {

	return s, err
}
