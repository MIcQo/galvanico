package main

import (
	"galvanico/cmd"
	"galvanico/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func main() {
	var cfg, cfgErr = config.Load()
	if cfgErr != nil {
		panic(cfgErr)
	}

	var lvl, err = zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	zerolog.SetGlobalLevel(lvl)

	var formatter = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339, TimeLocation: time.UTC}

	log.Logger = log.Output(formatter).With().Timestamp().Logger()

	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		log.Debug().Msg("Debug logging enabled")
	}

	cmd.Execute()
}
