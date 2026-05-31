package infra

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisBroker struct {
	Client *redis.Client
}

func NewRedisBroker(client *redis.Client) *RedisBroker {
	return &RedisBroker{
		Client: client,
	}
}

func (r *RedisBroker) Publish(ctx context.Context, topic string, message any) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("redis publish marshal error %w", err)
	}

	if err := r.Client.Publish(ctx, topic, payload).Err(); err != nil {
		return fmt.Errorf("redis publish error %w", err)
	}

	return nil
}

func (r *RedisBroker) Subscribe(ctx context.Context, topic string) (<-chan []byte, error) {
	out := make(chan []byte)
	pubsub := r.Client.Subscribe(ctx, topic)
	go func() {
		defer close(out)
		defer pubsub.Close()
		for msg := range pubsub.Channel() {
			out <- []byte(msg.Payload)
		}
	}()

	return out, nil
}
