package repository

import "github.com/spf13/afero"

type Packages []Package

type Package struct {
	Name     string
	Version  string
	Licenses []string
}

type Repository interface {
	GetPackages(patterns []string) (Packages, error)
}

type Factory func(fs afero.Fs, wd string) Repository
