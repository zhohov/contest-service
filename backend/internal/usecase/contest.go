package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/zhohov/contest-service/backend/internal/domain"
)

var (
	ErrContestNotFound  = errors.New("contest not found")
	ErrContestForbidden = errors.New("contest update forbidden")
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

func (uc *UseCase) GetContestsByUserID(ctx context.Context, in GetContestsByUserIdInput) (*GetContestsByUserIdOutput, error) {
	contests, err := uc.contests.GetContestsByUserID(ctx, in.UserID)
	if err != nil {
		return nil, fmt.Errorf("get contests by user id usecase: %w", err)
	}

	return &GetContestsByUserIdOutput{
		Contests: mapContestsToDTOs(contests),
	}, nil
}

func (uc *UseCase) UpdateContest(ctx context.Context, in UpdateContestInput) (*UpdateContestOutput, error) {
	existingContest, err := uc.contests.GetContestByID(ctx, in.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrContestNotFound, err)
	}

	if existingContest.Creator.ID != in.CreatorID {
		return nil, ErrContestForbidden
	}

	problems := make([]domain.Problem, 0, len(in.Problems))
	for _, problemInput := range in.Problems {
		testCases := make([]domain.TestCase, 0, len(problemInput.TestCases))

		for _, testCaseInput := range problemInput.TestCases {
			testCases = append(
				testCases,
				*domain.NewTestCase(testCaseInput.Input, testCaseInput.Expected),
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
			return nil, fmt.Errorf("error update problem %w", err)
		}

		problem.ID = problemInput.ID
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
	contest.ID = in.ID

	if err := uc.contests.UpdateContest(ctx, contest); err != nil {
		return nil, fmt.Errorf("update contest usecase: %w", err)
	}

	return &UpdateContestOutput{
		Contest: mapContestToDTO(contest),
	}, nil
}
