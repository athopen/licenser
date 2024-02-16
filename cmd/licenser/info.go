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
		Name:   "info",
		Usage:  "N/A",
		Action: withProjectOptions(infoAction),
	}
}

func infoAction(_ *console.Context, opts *cli.ProjectOptions) error {
	repo, err := repository.LoadRepository(opts.Fs, opts.WorkingDir)
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

	for _, pkg := range repo.GetPackages(opts.NoDev) {
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
}
