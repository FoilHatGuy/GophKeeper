package app

import (
	"bufio"
	"fmt"
	"os"
)

func MainLoop() {
	fmt.Println()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
