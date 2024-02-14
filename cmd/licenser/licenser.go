package licenser

import (
	"github.com/athopen/licenser/internal"
	"github.com/symfony-cli/console"
)

var (
	dirFlag   = &console.StringFlag{Name: "dir", Usage: "Project directory"}
	noDevFlag = &console.BoolFlag{Name: "no-dev", Usage: "Exclude require-dev packages"}
)

func Application() *console.Application {
	return &console.Application{
		Name:      "license-checker",
		Copyright: "(c) 2024 Andreas Penz",
		Usage:     "N/A",
		Flags: []console.Flag{
			dirFlag,
		},
		Commands: []*console.Command{
			infoCommand(),
			checkCommand(),
		},
		Version:   internal.Version,
		Channel:   internal.Channel,
		BuildDate: internal.BuildDate,
	}
}
