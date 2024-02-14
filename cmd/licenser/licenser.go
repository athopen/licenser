package licenser

import (
	"github.com/symfony-cli/console"
)

var (
	// version is overridden at linking time
	version = "dev"
)

var (
	fileFlag  = &console.StringFlag{Name: "file", Usage: "Config file", DefaultValue: "licenser.yml"}
	dirFlag   = &console.StringFlag{Name: "dir", Usage: "Project directory"}
	noDevFlag = &console.BoolFlag{Name: "no-dev", Usage: "Exclude require-dev packages"}
)

func Application() *console.Application {
	return &console.Application{
		Name:      "licenser",
		Copyright: "(c) 2024 Andreas Penz",
		Usage:     "Licenser is a tool designed to check and report on the licenses used by a package and its dependencies.",
		Flags: []console.Flag{
			fileFlag,
			dirFlag,
		},
		Commands: []*console.Command{
			infoCommand(),
			checkCommand(),
		},
		Version: version,
	}
}
