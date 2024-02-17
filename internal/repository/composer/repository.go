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

var (
	installedJSONPath = filepath.Join("vendor", "composer", "installed.json")
)

func (r Repository) GetPackages(noDev bool, patterns []string) (repository.Packages, error) {
	contents, err := filesystem.ReadFile(r.fs, filepath.Join(r.wd, installedJSONPath))
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

	var packages repository.Packages
	for _, p := range repo.Packages {
		isDev := slices.Contains(repo.DevPackageNames, p.Name)
		if noDev && isDev {
			continue
		}

		if wildecard.Match(p.Name, patterns) {
			continue
		}

		pkg := repository.Package{
			Dev:      isDev,
			Name:     p.Name,
			Version:  p.Version,
			Licenses: p.Licenses,
		}

		packages = append(packages, pkg)
	}

	return packages, nil
}
