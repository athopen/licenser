package config

import (
	"github.com/athopen/licenser/internal/config/schema"
	"github.com/athopen/licenser/internal/filesystem"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

type Project struct {
	Licenses []string `yaml:"licenses"`
	Packages []string `yaml:"packages"`
}

func LoadProject(fs afero.Fs, path string) (*Project, error) {
	contents, err := filesystem.ReadFile(fs, path)
	if err != nil {
		return nil, err
	}

	var dict map[string]interface{}
	if err := yaml.Unmarshal(contents, &dict); err != nil {
		return nil, err
	}

	if err := schema.Validate(dict); err != nil {
		return nil, err
	}

	project := &Project{}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:   project,
		Metadata: &mapstructure.Metadata{},
	})

	if err != nil {
		return nil, err
	}

	if err := decoder.Decode(dict); err != nil {
		return nil, err
	}

	return project, nil

}
