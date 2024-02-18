package licenser

import (
	"strings"

	"github.com/athopen/licenser/internal/cli"
	"github.com/athopen/licenser/internal/repository"

	"github.com/olekukonko/tablewriter"
	"github.com/symfony-cli/console"
	"github.com/symfony-cli/terminal"
)

func infoCommand() *console.Command {
	return &console.Command{
		Name:  "info",
		Usage: "List licenses of dependencies",
		Args: []*console.Arg{
			managerArg,
		},
		Flags: []console.Flag{
			dirFlag,
			noDevFlag,
		},
		Action: func(ctx *console.Context) error {
			opts, err := cli.NewProjectOptions(
				cli.WithWorkingDir(fs, ctx.String(dirFlag.Name)),
			)

			if err != nil {
				return err
			}

			factory, err := resolveFactory(ctx.Args().Get(managerArg.Name))
			if err != nil {
				return err
			}

			return infoAction(ctx, factory(fs, opts.WorkingDir))
		},
	}
}

func infoAction(ctx *console.Context, repo repository.Repository) error {
	pkgs, err := repo.GetPackages(ctx.Bool(noDevFlag.Name), []string{})
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
