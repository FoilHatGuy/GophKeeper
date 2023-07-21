package main

import (
	"fmt"

	"gophKeeper/src/client/app"
)

const (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("buildVersion\t= %q\n", buildVersion)
	fmt.Printf("buildDate\t= %q\n", buildDate)
	fmt.Printf("buildCommit\t= %q\n", buildCommit)
	fmt.Print("Wow, client is running!")
	app.MainLoop()
}
