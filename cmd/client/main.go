package main

import (
	"fmt"

	"gophKeeper/src/client/application"
	"gophKeeper/src/client/cfg"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	config := cfg.New(
		cfg.FromDefaults(),
		cfg.FromJSON(),
		cfg.FromFlags(),
		cfg.WithBuild(&cfg.BuildT{
			BuildVersion: buildVersion,
			BuildDate:    buildDate,
			BuildCommit:  buildCommit,
		}),
	)
	config.Save()
	fmt.Printf("buildVersion\t= %q\n", buildVersion)
	fmt.Printf("buildDate\t= %q\n", buildDate)
	fmt.Printf("buildCommit\t= %q\n", buildCommit)
	fmt.Print("Wow, client is running!")
	application.New(config)
}
