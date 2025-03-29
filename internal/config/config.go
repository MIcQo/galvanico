package config

import (
	"github.com/nats-io/nats.go"
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
	Broker   Broker   `yaml:"broker"`
}

type Database struct {
	URL string `yaml:"url"`
}

type Broker struct {
	URL string `yaml:"url"`
}

var once sync.Once
var cfg *Config

func Load() (*Config, error) {
	var outputErr error
	once.Do(func() {
		cfg, outputErr = loadFile(FileName)
	})

	return cfg, outputErr
}

func newDefaultConfig() *Config {
	return &Config{
		AppName:  "app",
		LogLevel: "info",
		Database: Database{
			// default postgres config, which should fail on other than local env
			URL: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
		},
		Broker: Broker{
			URL: nats.DefaultURL,
		},
	}
}

func loadFile(filename string) (*Config, error) {
	log.Info().Str("filename", filename).Msg("loading config file")

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Warn().
			Str("filename", filename).
			Bool("exists", false).
			Msg("using default config, please define your own")

		return newDefaultConfig(), nil
	}

	var file, err = os.ReadFile(filename)
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
