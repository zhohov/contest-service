package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/zhohov/contest-service/backend/internal/usecase"
)

type MessageBroker interface {
	Publish(ctx context.Context, topic string, message any) error
	Subscribe(ctx context.Context, topic string) (<-chan []byte, error)
}

type UseCase interface {
	UpdateSubmission(ctx context.Context, in usecase.UpdateSubmissionInput) (*usecase.UpdateSubmissionOutput, error)
}

type Consumer struct {
	broker MessageBroker
	uc     UseCase
}

func NewConsumer(broker MessageBroker, usecase UseCase) *Consumer {
	return &Consumer{
		broker: broker,
		uc:     usecase,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		if err := c.consumeSubmissionComplete(ctx); err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("consumer run: %w", err)
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Consumer) consumeSubmissionComplete(ctx context.Context) error {
	out, err := c.broker.Subscribe(ctx, "submissions.completed")
	if err != nil {
		return fmt.Errorf("subscribe submissions.completed: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case msg, ok := <-out:
			if !ok {
				return nil
			}

			var event SubmissionCompleteEventJson
			if err := json.Unmarshal(msg, &event); err != nil {
				log.Println(err)
				continue
			}

			_, err = c.uc.UpdateSubmission(
				ctx,
				mapSubmissionCompleteToInput(event),
			)
			if err != nil {
				log.Println("update submission:", err)
				continue
			}
		}
	}
}
