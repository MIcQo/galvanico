package notifications

import (
	"context"
	"galvanico/internal/broker"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

const msgBuffer = 1000

func NewConsumer(ctx context.Context) error {
	var messages = make(chan *nats.Msg, msgBuffer)
	var sub, err = broker.Connection().ChanSubscribe("channels.*", messages)
	if err != nil {
		return err
	}

	defer func(sub *nats.Subscription) {
		subErr := sub.Unsubscribe()
		if subErr != nil {
			panic(subErr)
		}
		close(messages)
	}(sub)

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-messages:
			log.Debug().
				Str("subject", msg.Subject).
				Str("body", string(msg.Data)).Msg("received message")
		}
	}
}
