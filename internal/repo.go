package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"

	"github.com/mitchellh/go-homedir"
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

func NewRepo(path string) (*Repo, error) {
	path, err := homedir.Expand(path)
	if err != nil {
		return nil, err
	}

	if path == "" {
		var err error
		if path, err = os.Getwd(); err != nil {
			return nil, err
		}
	}

	reader, err := os.Open(filepath.Join(path, "vendor", "composer", "installed.json"))
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("unable to read installed.json")
	}

	var repo *Repo
	if err = json.Unmarshal(content, &repo); err != nil {
		return nil, fmt.Errorf("installed.json does not contain valid JSON")
	}

	return repo, nil
}

func (r *Repo) GetPackages(noDev bool) []Package {
	pkgs := make([]Package, 0)
	for _, pkg := range r.Packages {
		if noDev == true && slices.Contains(r.DevPackageNames, pkg.Name) {
			continue
		}

		pkgs = append(pkgs, pkg)
	}

	return pkgs
}
