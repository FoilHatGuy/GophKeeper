package main

import (
	"fmt"

	"gophKeeper/src/server/cfg"
	app "gophKeeper/src/server/server"
)

const (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	config := cfg.New(
		cfg.FromDefaults(),
		cfg.FromEnv(),
		cfg.FromJSON(),
	)

	fmt.Printf("buildVersion\t= %q\n", buildVersion)
	fmt.Printf("buildDate\t= %q\n", buildDate)
	fmt.Printf("buildCommit\t= %q\n", buildCommit)
	fmt.Print("Wow, sever is running!")
	app.RunHTTPServer(config)
	app.RunGRPCServer(config)
}
