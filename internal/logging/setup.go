package logging

import (
	"galvanico/internal/config"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetupLevel(cfg *config.Config) error {
	var lvl, err = zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		return err
	}

	zerolog.SetGlobalLevel(lvl)

	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		log.Debug().Msg("Debug logging enabled")
	}

	return nil
}

func Setup() error {
	var formatter = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339, TimeLocation: time.UTC}

	log.Logger = log.Output(formatter). //nolint:reassign // because we need to set default log for app
						With().
						Timestamp().
						Logger()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	return nil
}
