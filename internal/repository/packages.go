package repository

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"slices"

	"github.com/athopen/licenser/internal/filesystem"
	"github.com/athopen/licenser/internal/wildecard"
	"github.com/spf13/afero"
)

type Packages []Package

type Package struct {
	Dev      bool
	Name     string   `json:"name"`
	Version  string   `json:"version_normalized"`
	Licenses []string `json:"license"`
}

var (
	installedJSONPath = filepath.Join("vendor", "composer", "installed.json")
)

func LoadPackages(fs afero.Fs, wd string, noDev bool, patterns []string) (Packages, error) {
	contents, err := filesystem.ReadFile(fs, filepath.Join(wd, installedJSONPath))
	if err != nil {
		return nil, err
	}

	var repo struct {
		Packages []struct {
			Name     string   `json:"name"`
			Version  string   `json:"version_normalized"`
			Licenses []string `json:"license"`
		} `json:"packages"`
		DevPackageNames []string `json:"dev-package-names"`
	}

	if err = json.Unmarshal(contents, &repo); err != nil {
		return nil, fmt.Errorf("installed.json does not contain valid JSON")
	}

	var packages Packages
	for _, p := range repo.Packages {
		isDev := slices.Contains(repo.DevPackageNames, p.Name)
		if noDev && isDev {
			continue
		}

		if wildecard.Match(p.Name, patterns) {
			continue
		}

		pkg := Package{
			Dev:      isDev,
			Name:     p.Name,
			Version:  p.Version,
			Licenses: p.Licenses,
		}

		packages = append(packages, pkg)
	}

	return packages, nil
}
