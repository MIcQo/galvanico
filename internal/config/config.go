package config

import (
	"os"
	"sync"

	"github.com/rs/zerolog/log"

	"gopkg.in/yaml.v3"
)

var FileName = "config.yaml"

type Config struct {
	AppName  string   `yaml:"appName"`
	LogLevel string   `yaml:"logLevel"`
	Database Database `yaml:"database"`
}

type Database struct {
	URL string `yaml:"url"`
}

var once sync.Once
var cfg *Config

func Load() (*Config, error) {
	var outputErr error
	once.Do(func() {
		cfg, outputErr = loadFile()
	})

	return cfg, outputErr
}

func newDefaultConfig() *Config {
	log.Warn().Msg("using default config, please define your own")
	return &Config{
		AppName:  "app",
		LogLevel: "info",
		Database: Database{
			// default postgres config, which should fail on other than local env
			URL: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
		},
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return os.IsExist(err)
}

func loadFile() (*Config, error) {
	log.Debug().Str("filename", FileName).Msg("loading config file")

	if !fileExists(FileName) {
		return newDefaultConfig(), nil
	}

	var file, err = os.ReadFile(FileName)
	if err != nil {
		return nil, err
	}

	var config = new(Config)
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
