package application

import (
	"context"
	"fmt"
	"strings"

	"gophKeeper/src/client/cfg"
)

type stateConfigType struct {
	stateGetName
	app    *Application
	config *cfg.ConfigT
}

func newConfigState(app *Application, config *cfg.ConfigT) state {
	return &stateConfigType{
		app:          app,
		config:       config,
		stateGetName: stateGetName{stateName: "Config view"},
	}
}

var (
	commandAbout = []string{"about", "abt", "a"}
	commandMod   = []string{"mod", "m"}
	commandList  = []string{"list", "l"}
)

func (s *stateConfigType) execute(ctx context.Context, command string) (resultState state, err error) {
	arguments := strings.Split(command, " ")
	switch {
	case includes(commandHelp, strings.ToLower(arguments[0])):
		fmt.Printf("This is Login screen. Available commands:\n"+
			"help   - %q - shows available commands (this screen)\n"+
			"back   - %q - return to menu\n"+
			"about  - %q - print build info\n"+
			"list   - %q - lists all configs\n"+
			"modify - %q - modify configs",
			commandHelp,
			commandBack,
			commandAbout,
			commandList,
			commandMod)
		return s, nil

	case includes(commandBack, strings.ToLower(arguments[0])):
		return s.app.cat[stateMenu], nil

	case includes(commandList, strings.ToLower(arguments[0])):
		fmt.Println(
			"\t1 - Server Address: "+s.app.config.ServerAddress,
			"\n\t2 - Secret Path   : "+s.app.config.SecretPath,
		)
		return s, nil

	case includes(commandAbout, strings.ToLower(arguments[0])):
		fmt.Print(colorGreen +
			"\tGophKeeper  v" + s.app.config.Build.BuildVersion +
			colorYellow +
			"\n\tcommit      #" + s.app.config.Build.BuildCommit +
			"\n\tbuild date: " + s.app.config.Build.BuildDate +
			colorNone + "\n",
		)
		return s, nil

	// I realised I can't modify cfg without restart :|
	case includes(commandMod, strings.ToLower(arguments[0])):
		if len(arguments) != 3 {
			return s, ErrUnrecognizedCommand
		}
		switch arguments[1] {
		case "1":
			s.app.config.ServerAddress = arguments[2]
		case "2":
			s.app.config.SecretPath = arguments[2]
		}

		s.app.config.Save()
		return s, nil
	}

	return s, ErrUnrecognizedCommand
}
