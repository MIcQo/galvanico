package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

var FileName = "config.yml"

type Config struct {
	AppName  string   `yaml:"app_name"`
	Debug    bool     `yaml:"debug"`
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
