package licenser

import (
	"fmt"
	"github.com/athopen/licenser/internal/cli"
	"github.com/athopen/licenser/internal/config"
	"github.com/athopen/licenser/internal/repository"
	"github.com/github/go-spdx/expression"
	"github.com/symfony-cli/console"
	"github.com/symfony-cli/terminal"
	"strings"
)

func checkCommand() *console.Command {
	return &console.Command{
		Name:  "check",
		Usage: "N/A",
		Flags: []console.Flag{
			dirFlag,
			fileFlag,
			noDevFlag,
		},
		Action: checkAction,
	}
}

func checkAction(ctx *console.Context) error {
	opts, err := cli.NewProjectOptions(
		cli.WithWorkingDir(fs, ctx.String(dirFlag.Name)),
		cli.WithConfigFile(fs, ctx.String(fileFlag.Name)),
		cli.WithNoDev(ctx.Bool(noDevFlag.Name)),
	)

	if err != nil {
		return err
	}

	project, err := config.LoadProject(fs, opts.ConfigFile)
	if err != nil {
		return err
	}

	pkgs, err := repository.LoadPackages(fs, opts.WorkingDir, opts.NoDev, project.Packages)
	if err != nil {
		return err
	}

	var violations []string
	for _, pkg := range pkgs {
		if len(pkg.Licenses) == 0 {
			violations = append(violations, fmt.Sprintf("No license found for package %s.", pkg.Name))

			continue
		}

		valid := false
		for _, license := range pkg.Licenses {
			license = strings.ReplaceAll(license, "or", "OR")
			license = strings.ReplaceAll(license, "and", "AND")

			satisfies, err := expression.Satisfies(license, project.Licenses)
			if satisfies && err == nil {
				valid = true
			}
		}

		if !valid {
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
