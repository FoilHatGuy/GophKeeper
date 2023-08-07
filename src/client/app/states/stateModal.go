package states

import (
	"context"
	"fmt"
	"gophKeeper/src/client/cfg"
)

type stateModal struct {
	app                        *Application
	config                     *cfg.ConfigT
	prompt                     string
	positiveFunc, negativeFunc modalFunction
}

type modalFunction func(ctx context.Context) (state, error)

func (a *Application) newModal(prompt string, positiveFunc, negativeFunc modalFunction) (resultState state) {
	fmt.Println(prompt)
	return &stateModal{
		app:          a,
		config:       a.config,
		prompt:       prompt,
		positiveFunc: positiveFunc,
		negativeFunc: negativeFunc,
	}
}

func (s *stateModal) Execute(ctx context.Context, command string) (resultState state, err error) {
	fmt.Println(s.prompt)
	switch command {
	case "y":
		resultState, err = s.positiveFunc(ctx)
	case "n":
		resultState, err = s.negativeFunc(ctx)
	}
	if resultState == nil {
		resultState = s
	}
	if err != nil {
		return resultState, err
	}
	return resultState, nil
}
