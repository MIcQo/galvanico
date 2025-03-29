package broker

import (
	"errors"
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

		ns.SetErrorHandler(natsErrHandler)

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

func natsErrHandler(_ *nats.Conn, sub *nats.Subscription, natsErr error) {
	log.Err(natsErr)

	if errors.Is(natsErr, nats.ErrSlowConsumer) {
		pendingMsgs, _, err := sub.Pending()
		if err != nil {
			log.Error().Err(err).
				Str("subject", sub.Subject).
				Msg("could not get pending messages")
			return
		}

		log.Warn().
			Str("subject", sub.Subject).
			Int("pending", pendingMsgs).
			Msg("slow consumer, pending messages")

		// TODO: Log error, notify operations...
	}
}
