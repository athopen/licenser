package npm

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/athopen/licenser/internal/filesystem"
	"github.com/athopen/licenser/internal/repository"
	"github.com/athopen/licenser/internal/wildecard"
	"github.com/spf13/afero"
)

type Repository struct {
	fs afero.Fs
	wd string
}

func Factory(fs afero.Fs, wd string) repository.Repository {
	return &Repository{
		fs: fs,
		wd: wd,
	}
}

type packageJSON struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	License  string `json:"license"`
	Licenses []struct {
		Type string `json:"type"`
	} `json:"licenses"`
}

func (r Repository) GetPackages(patterns []string) (repository.Packages, error) {
	paths, err := filesystem.Glob(r.fs, filepath.Join(r.wd, "node_modules", "*", "package.json"))
	if err != nil {
		return nil, err
	}

	var pkgs repository.Packages
	for _, path := range paths {
		decoded, err := readPackageJSON(r.fs, path)
		if err != nil {
			return nil, err
		}

		if wildecard.Match(decoded.Name, patterns) {
			continue
		}

		var licenses []string
		if decoded.License != "" {
			licenses = append(licenses, decoded.License)
		} else {
			for _, license := range decoded.Licenses {
				licenses = append(licenses, license.Type)
			}
		}

		pkgs = append(pkgs, repository.Package{
			Name:     decoded.Name,
			Version:  decoded.Version,
			Licenses: licenses,
		})
	}

	return pkgs, nil
}

func readPackageJSON(fs afero.Fs, path string) (*packageJSON, error) {
	contents, err := filesystem.ReadFile(fs, path)
	if err != nil {
		return nil, err
	}

	var decoded packageJSON
	if err = json.Unmarshal(contents, &decoded); err != nil {
		return nil, fmt.Errorf("\"%s\" does not contain valid JSON", path)
	}

	return &decoded, nil
}
