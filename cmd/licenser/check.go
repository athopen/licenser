package licenser

import (
	"fmt"

	"github.com/symfony-cli/console"
)

func checkCommand() *console.Command {
	return &console.Command{
		Name:  "check",
		Usage: "N/A",

		Action: func(ctx *console.Context) error {
			fmt.Println("check something")

			return nil
		},
	}
}
