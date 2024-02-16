package repository

import (
	"encoding/json"
	"fmt"
	"github.com/athopen/licenser/internal/filesystem"
	"github.com/spf13/afero"
	"path/filepath"
	"slices"
)

type Repo struct {
	Packages        []Package `json:"packages"`
	DevMode         bool      `json:"dev"`
	DevPackageNames []string  `json:"dev-package-names"`
}

type Package struct {
	Name     string   `json:"name"`
	Version  string   `json:"version_normalized"`
	Licenses []string `json:"license"`
}

var (
	installedJsonPath = filepath.Join("vendor", "composer", "installed.json")
)

func LoadRepository(fs afero.Fs, wd string) (*Repo, error) {
	contents, err := filesystem.ReadFile(fs, filepath.Join(wd, installedJsonPath))
	if err != nil {
		return nil, err
	}

	var repo Repo
	if err = json.Unmarshal(contents, &repo); err != nil {
		return nil, fmt.Errorf("installed.json does not contain valid JSON")
	}

	return &repo, nil
}

func (r *Repo) GetPackages(noDev bool) []Package {
	packages := make([]Package, 0)
	for _, pkg := range r.Packages {
		if noDev == true && slices.Contains(r.DevPackageNames, pkg.Name) {
			continue
		}

		packages = append(packages, pkg)
	}

	return packages
}
