package app

import (
	"bufio"
	"fmt"
	"gophKeeper/src/client/app/states"
	"gophKeeper/src/client/cfg"
	"os"
)

func MainLoop(config *cfg.ConfigT) {
	fmt.Println()
	app := states.New(config)
	Run(config, app)
}

func Run(_ *cfg.ConfigT, app *states.Application) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		app.Execute(scanner.Text())
	}
}
