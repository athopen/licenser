package licenser

import (
	"fmt"
	"time"

	"github.com/athopen/licenser/internal/repository"
	"github.com/athopen/licenser/internal/repository/composer"

	"github.com/spf13/afero"
	"github.com/symfony-cli/console"
)

var (
	// version is overridden at linking time
	version = "dev"
	// overridden at linking time
	buildDate string
)

var (
	fs = afero.NewOsFs()
)

var (
	managerArg = &console.Arg{Name: "manager", Description: "The package manager", Default: "composer"}
)

var (
	fileFlag  = &console.StringFlag{Name: "file", Usage: "Config file"}
	dirFlag   = &console.StringFlag{Name: "dir", Usage: "Working directory"}
	noDevFlag = &console.BoolFlag{Name: "no-dev", Usage: "Excluded require-dev packages"}
)

var (
	factories = map[string]repository.Factory{
		"composer": composer.Factory,
	}
)

func resolveFactory(manager string) (repository.Factory, error) {
	factory, exists := factories[manager]
	if !exists {
		return nil, fmt.Errorf("package manager \"%s\" is not supported", manager)
	}

	return factory, nil
}

var (
	helpTemplate = `<info>
        _____ _______ _______ __   _ _______ _______  ______
 |        |   |       |______ | \  | |______ |______ |_____/
 |_____ __|__ |_____  |______ |  \_| ______| |______ |    \_
</>

<info>{{.Name}}</>{{if .Version}} version <comment>{{.Version}}</>{{end}}{{if .Copyright}} {{.Copyright}}{{end}}

{{.Usage}}

<comment>Usage</>:
  {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} <command> [command options]{{end}} [arguments...]{{if .Description}}

{{.Description}}{{end}}{{if .VisibleFlags}}

<comment>Global options:</>
  {{range $index, $option := .VisibleFlags}}{{if $index}}
  {{end}}{{$option}}{{end}}{{end}}{{if .VisibleCommands}}

<comment>Available commands:</>{{range .VisibleCategories}}{{if .Name}}
 <comment>{{.Name}}</>{{"\t"}}{{end}}{{range .VisibleCommands}}
  <info>{{join .Names ", "}}</>{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}
`
)

func Application() *console.Application {
	return &console.Application{
		Name:      "licenser",
		Copyright: fmt.Sprintf("(c) %d <info>Andreas Penz <andreas.penz.1989@gmail.com></>", time.Now().Year()),
		Usage:     "Licenser is a tool designed to check and report on the licenses used by a package and its dependencies.",
		Action: func(ctx *console.Context) error {
			console.HelpPrinter(ctx.App.Writer, helpTemplate, ctx.App)
			return nil
		},
		Commands: []*console.Command{
			infoCommand(),
			checkCommand(),
		},
		Version:   version,
		BuildDate: buildDate,
	}
}
