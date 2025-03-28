package config

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var FileName = "config.yaml"

type Config struct {
	AppName string `yaml:"appName"`

	LogLevel string `yaml:"logLevel"`

	Database Database `yaml:"database"`
}

type Database struct {
	URL string `yaml:"url"`
}

var once sync.Once
var cfg Config

func Load() (*Config, error) {
	var outputErr error
	once.Do(func() {
		var file, err = os.ReadFile(FileName)
		if err != nil {
			outputErr = err
			return
		}

		err = yaml.Unmarshal(file, &cfg)
		if err != nil {
			outputErr = err
			return
		}
	})

	return &cfg, outputErr
}
