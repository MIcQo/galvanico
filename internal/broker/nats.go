package broker

import (
	"galvanico/internal/config"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

var broker *nats.Conn
var once sync.Once

func Connection() *nats.Conn {
	once.Do(func() {
		var cfg, cfgErr = config.Load()
		if cfgErr != nil {
			panic(cfgErr)
		}
		var ns, err = nats.Connect(cfg.Broker.URL)
		if err != nil {
			panic(err)
		}

		broker = ns
		log.Debug().Str("uri", cfg.Broker.URL).Msg("connected to broker")
	})

	return broker
}

func Close() error {
	if err := broker.Drain(); err != nil {
		return err
	}

	broker.Close()
	return nil
}
