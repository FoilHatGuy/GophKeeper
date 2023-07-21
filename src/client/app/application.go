package app

import (
	"bufio"
	"fmt"
	"os"

	"gophKeeper/src/client/cfg"
)

func MainLoop(_ *cfg.ConfigT) {
	fmt.Println()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
