package repository

import (
	"encoding/json"
	"fmt"
	"github.com/athopen/licenser/internal/filesystem"
	"github.com/athopen/licenser/internal/wildecard"
	"github.com/spf13/afero"
	"path/filepath"
	"slices"
)

type Packages []Package

type Package struct {
	Dev      bool
	Name     string   `json:"name"`
	Version  string   `json:"version_normalized"`
	Licenses []string `json:"license"`
}

var (
	installedJsonPath = filepath.Join("vendor", "composer", "installed.json")
)

func LoadPackages(fs afero.Fs, wd string, noDev bool, patterns []string) (Packages, error) {
	contents, err := filesystem.ReadFile(fs, filepath.Join(wd, installedJsonPath))
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

	matcher := wildecard.NewMatcher(patterns)

	var packages Packages
	for _, p := range repo.Packages {
		if noDev && slices.Contains(repo.DevPackageNames, p.Name) {
			continue
		}

		matches, err := matcher.Match(p.Name)
		if err != nil {
			return nil, err
		}

		if matches {
			continue
		}

		pkg := Package{
			Dev:      slices.Contains(repo.DevPackageNames, p.Name),
			Name:     p.Name,
			Version:  p.Version,
			Licenses: p.Licenses,
		}

		packages = append(packages, pkg)
	}

	return packages, nil
}
