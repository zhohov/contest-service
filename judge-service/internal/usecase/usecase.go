package usecase

import (
	"context"

	"github.com/zhohov/contest-service/judge-service/internal/domain"
)

type Judger interface {
	Submit(ctx context.Context, code string, testCases []domain.TestCase) (*domain.Submission, error)
}

type UseCase struct {
}

func NewUseCase() *UseCase {
	return &UseCase{}
}
