package app

import (
	"fmt"
	"gophKeeper/src/client/app/states"
	"gophKeeper/src/client/cfg"
)

func MainLoop(config *cfg.ConfigT) {
	fmt.Println("main loop func")
	app := states.New(config)
	app.Run()
}
