package licenser

import (
	"fmt"
	"github.com/athopen/licenser/internal"
	"github.com/github/go-spdx/expression"
	"github.com/symfony-cli/terminal"
	"path/filepath"
	"strings"

	"github.com/symfony-cli/console"
)

func checkCommand() *console.Command {
	return &console.Command{
		Name:  "check",
		Usage: "N/A",
		Action: func(ctx *console.Context) error {
			config, err := internal.LoadConfig(ctx.String(fileFlag.Name))
			if err != nil {
				return err
			}

			repo, err := internal.NewRepo(ctx.String(dirFlag.Name))
			if err != nil {
				return err
			}

			var violations []string
			for _, pkg := range repo.GetPackages(ctx.Bool(noDevFlag.Name)) {
				skip := false
				for _, pattern := range config.Packages {
					match, err := filepath.Match(pattern, pkg.Name)
					if err != nil {
						return err
					}

					if match {
						skip = true
						continue
					}
				}

				if skip {
					continue
				}

				if len(pkg.Licenses) == 0 {
					violations = append(violations, fmt.Sprintf("No license found for package %s.", pkg.Name))

					continue
				}

				valid := false
				for _, license := range pkg.Licenses {
					license = strings.ReplaceAll(license, "or", "OR")
					license = strings.ReplaceAll(license, "and", "AND")

					satisfies, err := expression.Satisfies(license, config.Licenses)
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
		},
	}
}
