package application

import (
	"context"

	"gophKeeper/src/client/cfg"
)

// nolint:unused
type stateFileType struct {
	stateName string
	app       *Application
	config    *cfg.ConfigT
}

// nolint:unused
func (s *stateFileType) execute(ctx context.Context, command string) (resultState state, err error) {
	return s, err
}

// nolint:unused
func (s *stateFileType) show(ctx context.Context) {
}

// nolint:unused
func (s *stateFileType) getName() string {
	return s.stateName
}
