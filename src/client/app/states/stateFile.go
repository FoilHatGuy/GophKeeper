package states

import (
	"context"

	"gophKeeper/src/client/cfg"
)

type stateFileType struct {
	app    *Application
	config *cfg.ConfigT
}

func (s *stateFileType) Execute(ctx context.Context, command string) (resultState state, err error) {
	return s, err
}

func (s *stateFileType) Show(ctx context.Context) {
}
