package application

import (
	"context"

	"gophKeeper/src/client/cfg"
)

// nolint:unused
type stateTextType struct {
	stateName string
	app       *Application
	config    *cfg.ConfigT
}

// nolint:unused
func (s *stateTextType) execute(ctx context.Context, command string) (resultState state, err error) {
	return s, err
}

// nolint:unused
func (s *stateTextType) show(ctx context.Context) {
}

// nolint:unused
func (s *stateTextType) getName() string {
	return s.stateName
}
