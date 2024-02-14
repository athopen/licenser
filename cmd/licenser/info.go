package licenser

import (
	"strings"

	"github.com/athopen/licenser/internal"
	"github.com/olekukonko/tablewriter"
	"github.com/symfony-cli/console"
	"github.com/symfony-cli/terminal"
)

func infoCommand() *console.Command {
	return &console.Command{
		Name:  "info",
		Usage: "N/A",
		Flags: []console.Flag{
			noDevFlag,
		},
		Action: func(ctx *console.Context) error {
			repo, err := internal.NewRepo(ctx.String(dirFlag.Name))
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

			for _, pkg := range repo.GetPackages(ctx.Bool(noDevFlag.Name)) {
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

			return nil
		},
	}
}
