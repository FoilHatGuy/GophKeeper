package states

import (
	"context"

	"gophKeeper/src/client/cfg"
)

type stateCardType struct {
	app    *Application
	config *cfg.ConfigT
}

func (s *stateCardType) Execute(ctx context.Context, command string) (resultState state, err error) {
	return s, err
}

func (s *stateCardType) Show(ctx context.Context) {
}
