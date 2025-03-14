package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ExitOnCompletion bool `yaml:"exit_on_completion"`
	Commands         []struct {
		Name          string `yaml:"name"`
		Command       string `yaml:"command"`
		Dir           string `yaml:"dir,omitempty"`
		ExitImportant bool   `yaml:"exit_important,omitempty"`
	} `yaml:"commands"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	return &cfg, err
}
