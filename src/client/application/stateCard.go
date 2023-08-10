package application

import (
	"context"

	"gophKeeper/src/client/cfg"
)

// nolint:unused
type stateCardType struct {
	stateName string
	app       *Application
	config    *cfg.ConfigT
}

// nolint:unused
func (s *stateCardType) execute(ctx context.Context, command string) (resultState state, err error) {
	return s, err
}

// nolint:unused
func (s *stateCardType) show(ctx context.Context) {
}

// nolint:unused
func (s *stateCardType) getName() string {
	return s.stateName
}
