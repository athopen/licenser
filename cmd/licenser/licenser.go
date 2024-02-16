package licenser

import (
	"github.com/athopen/licenser/internal/cli"
	"github.com/spf13/afero"
	"github.com/symfony-cli/console"
)

var (
	// version is overridden at linking time
	version = "dev"
)

var (
	osFs = afero.NewOsFs()
)

var (
	fileFlag  = &console.StringFlag{Name: "file", Usage: "Config file"}
	dirFlag   = &console.StringFlag{Name: "dir", Usage: "Working directory"}
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
			noDevFlag,
		},
		Commands: []*console.Command{
			infoCommand(),
			checkCommand(),
		},
		Version: version,
	}
}

type projectOptionsFunc func(ctx *console.Context, opts *cli.ProjectOptions) error

func withProjectOptions(fn projectOptionsFunc) func(ctx *console.Context) error {
	return func(ctx *console.Context) error {
		opts, err := cli.NewProjectOptions(
			osFs,
			cli.WithWorkingDir(ctx.String(dirFlag.Name)),
			cli.WithConfigFile(ctx.String(fileFlag.Name)),
			cli.WithNoDev(ctx.Bool(noDevFlag.Name)),
		)

		if err != nil {
			return err
		}

		if err != nil {
			return err
		}

		return fn(ctx, opts)
	}
}
