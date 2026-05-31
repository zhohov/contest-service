package usecase

import (
	"context"
	"fmt"

	"github.com/zhohov/contest-service/backend/internal/domain"
)

func (uc *UseCase) CreateProblem(ctx context.Context, in CreateProblemInput) (*CreateProblemOutput, error) {
	testCases := make([]domain.TestCase, 0, len(in.TestCases))
	for _, testCase := range in.TestCases {
		testCases = append(testCases, *domain.NewTestCase(testCase.Input, testCase.Expected))
	}

	problem, err := domain.NewProblem(in.Title, in.Slug, in.Difficulty, in.Content, testCases)
	if err != nil {
		return nil, fmt.Errorf("error create problem %w", err)
	}

	createdProblem, err := uc.problems.CreateProblem(ctx, problem)
	if err != nil {
		return nil, err
	}

	return &CreateProblemOutput{
		Problem: mapProblemToDTO(createdProblem),
	}, nil
}

func (uc *UseCase) GetProblems(ctx context.Context) (*GetProblemsOutput, error) {
	problems, err := uc.problems.GetProblems(ctx)
	if err != nil {
		return nil, err
	}

	return &GetProblemsOutput{
		Problems: mapProblemsToDTOs(problems),
	}, nil
}

func (uc *UseCase) GetProblemBySlug(ctx context.Context, in GetProblemBySlugInput) (*GetProblemBySlugOutput, error) {
	problem, err := uc.problems.GetProblemBySlug(ctx, in.Slug)
	if err != nil {
		return nil, err
	}

	return &GetProblemBySlugOutput{
		Problem: mapProblemToDTO(problem),
	}, nil
}

// func (uc *UseCase) GetProblemById(in GetProblemById) (*GetProblemBySlugOutput, error) {
// 	problem, err := uc.problems.GetProblemById(in.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &GetProblemBySlugOutput{
// 		Problem: mapProblemToDTO(problem),
// 	}, nil
// }
