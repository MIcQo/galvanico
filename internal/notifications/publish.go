package notifications

import (
	"galvanico/internal/broker"
)

func NewPublisher(subject string, msg []byte) error {
	return broker.Connection().Publish(subject, msg)
}
