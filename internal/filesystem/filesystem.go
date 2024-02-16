package filesystem

import "github.com/spf13/afero"

func DirExists(fs afero.Fs, path string) (bool, error) {
	return afero.DirExists(fs, path)
}

func Exists(fs afero.Fs, path string) (bool, error) {
	return afero.Exists(fs, path)
}

func ReadFile(fs afero.Fs, path string) ([]byte, error) {
	return afero.ReadFile(fs, path)
}
