package main

import (
	"os"

	"github.com/athopen/licenser/cmd/licenser"
)

func main() {
	if err := licenser.Application().Run(os.Args); err != nil {
		os.Exit(-1)
	}
}
