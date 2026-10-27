package http

import (
	"context"

	"github.com/zhohov/contest-service/backend/internal/domain"
	"github.com/zhohov/contest-service/backend/internal/usecase"
)

// type ProblemUseCase interface {
// }

// type JudgeUseCase interface {
// }

// type ArticleUseCase interface {
// 	CreateArticle()
// 	GetArticles() []*domain.Article
// }

// type UseCase struct {
// 	problem *ProblemUseCase
// 	judge   *JudgeUseCase
// 	article *ArticleUseCase
// }

type UseCase interface {
	CreateArticle(ctx context.Context, in usecase.CreateArticleInput) (*usecase.CreateArticleOutput, error)
	GetArticles(ctx context.Context) (*usecase.GetArticlesOutput, error)

	CreateProblem(ctx context.Context, in usecase.CreateProblemInput) (*usecase.CreateProblemOutput, error)
	GetProblems(ctx context.Context) (*usecase.GetProblemsOutput, error)
	GetProblemBySlug(ctx context.Context, in usecase.GetProblemBySlugInput) (*usecase.GetProblemBySlugOutput, error)

	Login(ctx context.Context, in usecase.LoginInput) (*usecase.LoginOutput, error)
	Register(ctx context.Context, in usecase.RegisterInput) (*usecase.RegisterOutput, error)

	CreateSubmission(ctx context.Context, in usecase.CreateSubmissionInput) (*usecase.CreateSubmissionOutput, error)
	GetSubmissionByID(ctx context.Context, in usecase.GetSubmissionByIdInput) (*usecase.GetSubmissionByIdOutput, error)
	GetSubmissionsByUserID(ctx context.Context, in usecase.GetSubmissionsByUserIdInput) (*usecase.GetSubmissionsByUserIdOutput, error)
	GetSubmissionsByContestID(ctx context.Context, in usecase.GetSubmissionsByContestIdInput) (*usecase.GetSubmissionsByContestIdOutput, error)

	CreateContest(ctx context.Context, in usecase.CreateContestInput) (*usecase.CreateContestOutput, error)
	UpdateContest(ctx context.Context, in usecase.UpdateContestInput) (*usecase.UpdateContestOutput, error)
	GetContestsByUserID(ctx context.Context, in usecase.GetContestsByUserIdInput) (*usecase.GetContestsByUserIdOutput, error)
	GetContestBySlug(ctx context.Context, in usecase.GetContestBySlugInput) (*usecase.GetContestBySlugOutput, error)
}

type JwtService interface {
	ParseToken(token string) (*domain.AuthUser, error)
}

type HTTPHandler struct {
	uc  UseCase
	jwt JwtService
}

func NewHTTPHandler(uc UseCase, jwt JwtService) *HTTPHandler {
	return &HTTPHandler{
		uc:  uc,
		jwt: jwt,
	}
}
