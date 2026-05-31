package event

import "github.com/zhohov/contest-service/backend/internal/usecase"

func mapSubmissionCompleteToInput(req SubmissionCompleteEventJson) usecase.UpdateSubmissionInput {
	return usecase.UpdateSubmissionInput{
		SubmissionID: req.SubmissionID,
		Status:       req.Status,
		Result:       req.Result,
	}
}
