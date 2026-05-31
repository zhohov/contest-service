package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/moby/moby/client"
	"github.com/redis/go-redis/v9"
	"github.com/zhohov/contest-service/judge-service/config"
	"github.com/zhohov/contest-service/judge-service/internal/domain"
	"github.com/zhohov/contest-service/judge-service/internal/infra"
)

type TestCaseEventDTO struct {
	Input    string
	Expected string
}

type SubmissionCreateEvent struct {
	SubmissionID int
	Status       domain.SubmissionStatus
	Result       domain.SubmissionResult
	Code         string
	TestCases    []TestCaseEventDTO
}

type SubmissionCompletedEvent struct {
	SubmissionID int                     `json:"submission_id"`
	Status       domain.SubmissionStatus `json:"status"`
	Result       domain.SubmissionResult `json:"result"`
}

func Run(ctx context.Context, config *config.Config) error {
	log.Println("Judge")

	dockerCli, err := client.New(client.FromEnv)
	if err != nil {
		log.Println(err)
	}
	judger := infra.NewDockerJudger(dockerCli)

	client := redis.NewClient(&redis.Options{Addr: "redis:6379", Password: "redis"})
	broker := infra.NewRedisBroker(client)

	// server := event.NewEventServer(broker)
	// if err := server.Run(); err != nil {
	//     panic()
	// }

	out, err := broker.Subscribe(ctx, "submissions.created")
	if err != nil {
		return fmt.Errorf("redis subscribe error %w", err)
	}

	for msg := range out {
		var event SubmissionCreateEvent
		if err := json.Unmarshal(msg, &event); err != nil {
			log.Println(err)
			continue
		}
		log.Println(event)
		testCases := make([]infra.TestCase, 0, len(event.TestCases))
		for _, tc := range event.TestCases {
			testCases = append(testCases, infra.TestCase{
				Input:    tc.Input,
				Expected: tc.Expected,
			})
		}
		result, err := judger.Submit(ctx, event.Code, testCases)
		if err != nil {
			log.Println("submit error:", err)
			continue
		}
		log.Printf("submit result: passed=%v", result)

		completedEvent := SubmissionCompletedEvent{
			SubmissionID: event.SubmissionID,
			Status:       domain.SubmissionStatsDone,
			Result:       result.Result,
		}

		log.Printf("completedEvent: %v", completedEvent)

		if err := broker.Publish(ctx, "submissions.completed", completedEvent); err != nil {
			log.Println("publish submissions.completed error:", err)
		}
	}

	return nil

	// pubsub := client.Subscribe(ctx, "test-channel")
	// defer pubsub.Close()
	// ch := pubsub.Channel()
	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		return ctx.Err()
	// 	case msg := <-ch:
	// 		var event Event
	// 		err := json.Unmarshal(
	// 			[]byte(msg.Payload),
	// 			&event,
	// 		)
	// 		if err != nil {
	// 			log.Println(err)
	// 			continue
	// 		}
	// 		log.Println("received event:", event.Message)
	// 	}
	// }
}
