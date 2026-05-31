package usecase

import (
	"context"
	"fmt"

	"github.com/zhohov/contest-service/backend/internal/domain"
)

func (uc *UseCase) CreateContest(
	ctx context.Context,
	in CreateContestInput,
) (*CreateContestOutput, error) {
	problems := make([]domain.Problem, 0, len(in.Problems))

	for _, problemInput := range in.Problems {
		testCases := make([]domain.TestCase, 0, len(problemInput.TestCases))

		for _, testCaseInput := range problemInput.TestCases {
			testCases = append(
				testCases,
				*domain.NewTestCase(
					testCaseInput.Input,
					testCaseInput.Expected,
				),
			)
		}

		problem, err := domain.NewProblem(
			problemInput.Title,
			problemInput.Slug,
			problemInput.Difficulty,
			problemInput.Content,
			testCases,
		)
		if err != nil {
			return nil, fmt.Errorf("error create problem %w", err)
		}

		problems = append(problems, *problem)
	}

	contest := domain.NewContest(
		domain.User{ID: in.CreatorID},
		in.Title,
		in.Slug,
		in.Description,
		problems,
		in.StartsAt,
		in.EndsAt,
	)

	createdContest, err := uc.contests.CreateContest(ctx, contest)
	if err != nil {
		return nil, err
	}

	return &CreateContestOutput{
		Contest: mapContestToDTO(createdContest),
	}, nil
}

func (uc *UseCase) GetContestBySlug(ctx context.Context, in GetContestBySlugInput) (*GetContestBySlugOutput, error) {
	contest, err := uc.contests.GetContestBySlug(ctx, in.Slug)
	if err != nil {
		return nil, err
	}

	return &GetContestBySlugOutput{
		Contest: mapContestToDTO(contest),
	}, nil
}
