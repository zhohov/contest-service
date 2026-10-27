package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/zhohov/contest-service/backend/internal/domain"
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

func (uc *UseCase) CreateSubmission(ctx context.Context, in CreateSubmissionInput) (*CreateSubmissionOutput, error) {
	problem, err := uc.problems.GetProblemById(ctx, in.ProblemID)
	if err != nil {
		log.Println(fmt.Errorf("get problem by id usecase %w", err))
		return nil, fmt.Errorf("get problem by id usecase %w", err)
	}

	user, err := uc.users.GetUserById(ctx, in.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user from uc: %w", err)
	}

	submission := domain.NewSubmission(in.Code, *user, *problem)
	createdSubmission, err := uc.submissions.CreateSubmission(ctx, submission)
	if err != nil {
		return nil, fmt.Errorf("create submission: %w", err)
	}

	testCases := []TestCaseEventDTO{}
	for _, testCase := range submission.Problem.TestCases {
		testCases = append(testCases, TestCaseEventDTO{
			Input:    testCase.Input,
			Expected: testCase.Expected,
		})
	}

	event := SubmissionCreateEvent{
		SubmissionID: createdSubmission.ID,
		Status:       submission.Status,
		Result:       submission.Result,
		Code:         submission.Code,
		TestCases:    testCases,
	}

	if err := uc.broker.Publish(ctx, "submissions.created", event); err != nil {
		log.Println(fmt.Errorf("publish submission: %w", err))
		return nil, fmt.Errorf("publish submission: %w", err)
	}

	return &CreateSubmissionOutput{
		SubmissionID: createdSubmission.ID,
	}, nil
}

func (uc *UseCase) GetSubmissionByID(ctx context.Context, in GetSubmissionByIdInput) (*GetSubmissionByIdOutput, error) {
	submission, err := uc.submissions.GetSubmissionByID(ctx, in.ID)
	if err != nil {
		return nil, fmt.Errorf(": %w", err)
	}

	return &GetSubmissionByIdOutput{
		Submission: mapSubmissionToDTO(submission),
	}, nil
}

func (uc *UseCase) GetSubmissionsByUserID(ctx context.Context, in GetSubmissionsByUserIdInput) (*GetSubmissionsByUserIdOutput, error) {
	submissions, err := uc.submissions.GetSubmissionsByUserID(ctx, in.UserID)
	if err != nil {
		return nil, fmt.Errorf("get submissions by user id usecase: %w", err)
	}

	return &GetSubmissionsByUserIdOutput{
		Submissions: mapSubmissionsToDTOs(submissions),
	}, nil
}

func (uc *UseCase) GetSubmissionsByContestID(ctx context.Context, in GetSubmissionsByContestIdInput) (*GetSubmissionsByContestIdOutput, error) {
	submissions, err := uc.submissions.GetSubmissionsByContestID(ctx, in.ContestID)
	if err != nil {
		return nil, fmt.Errorf("get submissions by contest id usecase: %w", err)
	}

	return &GetSubmissionsByContestIdOutput{
		Submissions: mapSubmissionsToDTOs(submissions),
	}, nil
}

func (uc *UseCase) UpdateSubmission(ctx context.Context, in UpdateSubmissionInput) (*UpdateSubmissionOutput, error) {
	submission, err := uc.submissions.GetSubmissionByID(ctx, in.SubmissionID)
	if err != nil {
		return nil, fmt.Errorf(": %w", err)
	}

	submission.Status = in.Status
	submission.Result = in.Result

	err = uc.submissions.UpdateSubmission(ctx, submission)
	if err != nil {
		return nil, fmt.Errorf(": %w", err)
	}

	return &UpdateSubmissionOutput{}, nil
}
