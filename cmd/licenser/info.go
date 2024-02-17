package licenser

import (
	"github.com/athopen/licenser/internal/cli"
	"github.com/athopen/licenser/internal/repository"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/symfony-cli/console"
	"github.com/symfony-cli/terminal"
)

func infoCommand() *console.Command {
	return &console.Command{
		Name:  "info",
		Usage: "N/A",
		Flags: []console.Flag{
			dirFlag,
			noDevFlag,
		},
		Action: infoAction,
	}
}

func infoAction(ctx *console.Context) error {
	opts, err := cli.NewProjectOptions(
		cli.WithWorkingDir(fs, ctx.String(dirFlag.Name)),
		cli.WithNoDev(ctx.Bool(noDevFlag.Name)),
	)

	if err != nil {
		return err
	}

	pkgs, err := repository.LoadPackages(fs, opts.WorkingDir, opts.NoDev, []string{})
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(terminal.Stdout)
	table.SetAutoFormatHeaders(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{
		terminal.Format("<header>Name</>"),
		terminal.Format("<header>Version</>"),
		terminal.Format("<header>License</>"),
	})

	for _, pkg := range pkgs {
		licenseStr := "none"
		if len(pkg.Licenses) != 0 {
			licenseStr = strings.Join(pkg.Licenses, ", ")
		}

		table.Append([]string{
			pkg.Name,
			pkg.Version,
			licenseStr,
		})
	}

	table.Render()

	terminal.Printfln("<comment>%d packages found</>", len(pkgs))

	return nil
}
