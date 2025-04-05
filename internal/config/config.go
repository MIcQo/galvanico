package config

import (
	"galvanico/internal/utils"
	"os"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var FileName = "config.yaml"

const (
	AuthKeyLength = 32
)

type Config struct {
	AppName  string   `yaml:"appName"`
	LogLevel string   `yaml:"logLevel"`
	Database Database `yaml:"database"`
	Broker   Broker   `yaml:"broker"`
	Auth     *Auth    `yaml:"auth"`
}

type Auth struct {
	Provider string            `yaml:"provider"`
	Settings map[string]string `yaml:"settings"`
}

func (a *Auth) GetJWTKey() []byte {
	var key, ok = a.Settings["key"]

	if !ok || key == "" {
		key = utils.RandomString(AuthKeyLength)
		a.Settings["key"] = key
		log.Warn().
			Str("key", a.Settings["key"]).
			Msg(
				"field 'key' in 'auth.settings' was not found or is empty, defaulting to random key, if application " +
					"would be restarted, then every issued token would be invalid and every user should re-authorize. " +
					"Please set static key to 'auth.settings.key'",
			)
	}

	return []byte(key)
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

func NewDefaultConfig() *Config {
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
		Auth: &Auth{
			Provider: "jwt",
			Settings: map[string]string{
				"key": utils.RandomString(AuthKeyLength),
			},
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

		return NewDefaultConfig(), nil
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
