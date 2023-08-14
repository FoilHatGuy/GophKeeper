package application

import (
	"context"
	"fmt"

	"gophKeeper/src/client/cfg"
)

type stateModal struct {
	app       *Application
	config    *cfg.ConfigT
	stateName string
	prompt    string
	execFunc  modalFunction
}

func (s *stateModal) getName() string {
	return s.stateName
}

type modalFunction func(ctx context.Context, command string) (state, error)

func (a *Application) newModal(stateName, prompt string, execFunc modalFunction) (resultState *stateModal) {
	return &stateModal{
		app:       a,
		config:    a.config,
		stateName: stateName,
		prompt:    prompt,
		execFunc:  execFunc,
	}
}

func (s *stateModal) Print() {
	fmt.Println(s.prompt)
}

func (s *stateModal) execute(ctx context.Context, command string) (resultState state, err error) {
	resultState, err = s.execFunc(ctx, command)
	if resultState == nil {
		resultState = s
	}
	if err != nil {
		s.Print()
		return resultState, err
	}
	return resultState, nil
}
