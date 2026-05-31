package redis

import "context"

type MessageBroker interface {
	Publish(ctx context.Context, topic string, payload []byte) error
	Subscribe(ctx context.Context, topic string) (<-chan []byte, error)
}

type Consumer struct {
	broker MessageBroker
}

func NewConsumer(broker MessageBroker) *Consumer {
	return &Consumer{
		broker: broker,
	}
}
