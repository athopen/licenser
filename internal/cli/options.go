package cli

import (
	"fmt"
	"path/filepath"

	"github.com/athopen/licenser/internal/filesystem"
	"github.com/spf13/afero"
)

type ProjectOptions struct {
	WorkingDir string
	ConfigFile string
	NoDev      bool
}

type ProjectOptionsFn func(*ProjectOptions) error

func NewProjectOptions(fns ...ProjectOptionsFn) (*ProjectOptions, error) {
	opts := &ProjectOptions{}

	for _, fn := range fns {
		if err := fn(opts); err != nil {
			return nil, err
		}
	}

	return opts, nil
}

func WithWorkingDir(fs afero.Fs, path string) ProjectOptionsFn {
	return func(o *ProjectOptions) error {
		if path == "" {
			return nil
		}

		path, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		if found, _ := filesystem.DirExists(fs, path); !found {
			return fmt.Errorf("dir does not exist at \"%s\"", path)
		}

		o.WorkingDir = path

		return nil
	}
}

var defaultFileName = "licenser.yml"

func WithConfigFile(fs afero.Fs, path string) ProjectOptionsFn {
	return func(o *ProjectOptions) error {
		if path == "" {
			path := filepath.Join(o.WorkingDir, defaultFileName)

			if found, _ := filesystem.Exists(fs, path); !found {
				return fmt.Errorf("config file does not exist at \"%s\"", path)
			}

			o.ConfigFile = path

			return nil
		}

		path, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		if found, _ := filesystem.Exists(fs, path); !found {
			return fmt.Errorf("config file does not exist at \"%s\"", path)
		}

		o.ConfigFile = path

		return nil
	}
}

func WithNoDev(noDev bool) ProjectOptionsFn {
	return func(o *ProjectOptions) error {
		o.NoDev = noDev

		return nil
	}
}
