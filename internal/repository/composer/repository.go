package composer

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"slices"

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

type InstalledJSONData struct {
	Packages []struct {
		Name     string   `json:"name"`
		Version  string   `json:"version_normalized"`
		Licenses []string `json:"license"`
	} `json:"packages"`
	DevPackageNames []string `json:"dev-package-names"`
}

func (r Repository) GetPackages(noDev bool, patterns []string) (repository.Packages, error) {
	installedJSON, err := readInstalledJSON(r.fs, filepath.Join(r.wd, "vendor", "composer", "installed.json"))
	if err != nil {
		return nil, err
	}

	var pkgs repository.Packages
	for _, pkg := range installedJSON.Packages {
		if noDev && slices.Contains(installedJSON.DevPackageNames, pkg.Name) {
			continue
		}

		if wildecard.Match(pkg.Name, patterns) {
			continue
		}

		pkgs = append(pkgs, repository.Package{
			Name:     pkg.Name,
			Version:  pkg.Version,
			Licenses: pkg.Licenses,
		})
	}

	return pkgs, nil
}

func readInstalledJSON(fs afero.Fs, path string) (*InstalledJSONData, error) {
	contents, err := filesystem.ReadFile(fs, path)
	if err != nil {
		return nil, err
	}

	var decoded InstalledJSONData
	if err = json.Unmarshal(contents, &decoded); err != nil {
		return nil, fmt.Errorf("installed.json does not contain valid JSON")
	}

	return &decoded, nil
}
