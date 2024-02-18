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

type PackageLockJSONData struct {
	Packages map[string]struct {
		Dev bool `json:"dev"`
	} `json:"packages"`
}

type PackageJSONData struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	License  string `json:"license"`
	Licenses []struct {
		Type string `json:"type"`
	} `json:"licenses"`
}

func (r Repository) GetPackages(noDev bool, patterns []string) (repository.Packages, error) {
	packageLockJSON, err := readPackageLockJSON(r.fs, filepath.Join(r.wd, "node_modules", ".package-lock.json"))
	if err != nil {
		return nil, err
	}

	paths, err := filesystem.Glob(r.fs, filepath.Join(r.wd, "node_modules", "*", "package.json"))
	if err != nil {
		return nil, err
	}

	var pkgs repository.Packages
	for _, path := range paths {
		rel, err := filepath.Rel(r.wd, path)
		if err != nil {
			return nil, err
		}

		pkgLockData, found := packageLockJSON.Packages[filepath.Dir(rel)]
		if !found {
			return nil, fmt.Errorf("\"%s\" unable to determine if package is installed as dev", path)
		}

		if noDev && pkgLockData.Dev {
			continue
		}

		packageJSONData, err := readPackageJSON(r.fs, path)
		if err != nil {
			return nil, err
		}

		if wildecard.Match(packageJSONData.Name, patterns) {
			continue
		}

		var licenses []string
		if packageJSONData.License != "" {
			licenses = append(licenses, packageJSONData.License)
		} else {
			for _, license := range packageJSONData.Licenses {
				licenses = append(licenses, license.Type)
			}
		}

		pkgs = append(pkgs, repository.Package{
			Name:     packageJSONData.Name,
			Version:  packageJSONData.Version,
			Licenses: licenses,
		})
	}

	return pkgs, nil
}

func readPackageLockJSON(fs afero.Fs, path string) (*PackageLockJSONData, error) {
	contents, err := filesystem.ReadFile(fs, path)
	if err != nil {
		return nil, err
	}

	var decoded PackageLockJSONData
	if err = json.Unmarshal(contents, &decoded); err != nil {
		return nil, fmt.Errorf("\"%s\" does not contain valid JSON", path)
	}

	return &decoded, nil
}

func readPackageJSON(fs afero.Fs, path string) (*PackageJSONData, error) {
	contents, err := filesystem.ReadFile(fs, path)
	if err != nil {
		return nil, err
	}

	var decoded PackageJSONData
	if err = json.Unmarshal(contents, &decoded); err != nil {
		return nil, fmt.Errorf("\"%s\" does not contain valid JSON", path)
	}

	return &decoded, nil
}
