package cli

import (
	"fmt"
	"github.com/athopen/licenser/internal/filesystem"
	"github.com/spf13/afero"
	"path/filepath"
)

type ProjectOptions struct {
	Fs afero.Fs

	WorkingDir string
	ConfigFile string
	NoDev      bool
}

type ProjectOptionsFn func(*ProjectOptions) error

func NewProjectOptions(fs afero.Fs, fns ...ProjectOptionsFn) (*ProjectOptions, error) {
	opts := &ProjectOptions{
		Fs: fs,
	}

	for _, fn := range fns {
		if err := fn(opts); err != nil {
			return nil, err
		}
	}

	return opts, nil
}

func WithWorkingDir(path string) ProjectOptionsFn {
	return func(o *ProjectOptions) error {
		if path == "" {
			return nil
		}

		path, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		if found, _ := filesystem.DirExists(o.Fs, path); !found {
			return fmt.Errorf("dir does not exist at \"%s\"", path)
		}

		o.WorkingDir = path

		return nil
	}
}

var defaultFileName = "licenser.yml"

func WithConfigFile(path string) ProjectOptionsFn {
	return func(o *ProjectOptions) error {
		if path == "" {
			path := filepath.Join(o.WorkingDir, defaultFileName)

			if found, _ := filesystem.Exists(o.Fs, path); !found {
				return fmt.Errorf("config file does not exist at \"%s\"", path)
			}

			o.ConfigFile = path

			return nil
		}

		path, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		if found, _ := filesystem.Exists(o.Fs, path); !found {
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
