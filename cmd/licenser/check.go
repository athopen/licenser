package licenser

import (
	"fmt"
	"strings"

	"github.com/athopen/licenser/internal/repository"

	"github.com/athopen/licenser/internal/license"

	"github.com/athopen/licenser/internal/cli"
	"github.com/athopen/licenser/internal/config"
	"github.com/symfony-cli/console"
	"github.com/symfony-cli/terminal"
)

func checkCommand() *console.Command {
	return &console.Command{
		Name:  "check",
		Usage: "Check licenses of dependencies",
		Args: []*console.Arg{
			managerArg,
		},
		Flags: []console.Flag{
			dirFlag,
			fileFlag,
		},
		Action: func(ctx *console.Context) error {
			opts, err := cli.NewProjectOptions(
				cli.WithWorkingDir(fs, ctx.String(dirFlag.Name)),
				cli.WithConfigFile(fs, ctx.String(fileFlag.Name)),
			)

			if err != nil {
				return err
			}

			factory, err := resolveFactory(ctx.Args().Get(managerArg.Name))
			if err != nil {
				return err
			}

			project, err := config.LoadProject(fs, opts.ConfigFile)
			if err != nil {
				return err
			}

			return checkAction(factory(fs, opts.WorkingDir), project)
		},
	}
}

func checkAction(repo repository.Repository, project *config.Project) error {
	pkgs, err := repo.GetPackages(project.Excluded)
	if err != nil {
		return err
	}

	var violations []string
	for _, pkg := range pkgs {
		if len(pkg.Licenses) == 0 {
			violations = append(violations, fmt.Sprintf("No license found for package %s.", pkg.Name))

			continue
		}

		if !license.Satisfies(pkg.Licenses, project.Licenses) {
			violations = append(violations, fmt.Sprintf("License \"%s\" of package %s is not allowed.", strings.Join(pkg.Licenses, ", "), pkg.Name))
		}
	}

	if len(violations) == 0 {
		terminal.Println("<info>License check passed!</>")

		return nil
	}

	for _, violation := range violations {
		terminal.Printfln("<comment>%s</>", violation)
	}

	terminal.Println("<error>License check failed!</>")

	return console.Exit("", 1)
}
