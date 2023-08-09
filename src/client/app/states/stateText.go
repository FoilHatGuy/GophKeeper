package states

import (
	"context"

	"gophKeeper/src/client/cfg"
)

type stateTextType struct {
	app    *Application
	config *cfg.ConfigT
}

func (s *stateTextType) Execute(ctx context.Context, command string) (resultState state, err error) {
	return s, err
}

func (s *stateTextType) Show(ctx context.Context) {
}
