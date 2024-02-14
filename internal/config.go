package internal

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

type Config struct {
	Licenses []string `yaml:"licenses"`
	Packages []string `yaml:"packages"`
}

func LoadConfig(path string) (*Config, error) {
	path, err := homedir.Expand(path)
	if err != nil {
		return nil, err
	}

	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("unable to read licenser.yaml")
	}

	var config Config
	if err = yaml.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf("unable to read licenser.yaml")
	}

	return &config, nil
}
