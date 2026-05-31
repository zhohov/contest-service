package event

import "github.com/zhohov/contest-service/backend/internal/domain"

type SubmissionCompleteEventJson struct {
	SubmissionID int                     `json:"submission_id"`
	Status       domain.SubmissionStatus `json:"status"`
	Result       domain.SubmissionResult `json:"result"`
}
