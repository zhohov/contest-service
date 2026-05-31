package usecase

import (
	"time"

	"github.com/zhohov/contest-service/backend/internal/domain"
)

type ArticleDTO struct {
	ID      int
	Title   string
	Slug    string
	Content string
}

type CreateArticleInput struct {
	Title   string
	Slug    string
	Content string
}

type CreateArticleOutput struct {
	Article ArticleDTO
}

type GetArticlesOutput struct {
	Articles []ArticleDTO
}

type TestCaseDTO struct {
	Input    string
	Expected string
}

type ProblemDTO struct {
	ID         int
	Title      string
	Slug       string
	Difficulty string
	Content    string
	TestCases  []TestCaseDTO
}

type CreateProblemInput struct {
	Title      string
	Slug       string
	Difficulty string
	Content    string
	TestCases  []TestCaseDTO
}

type CreateProblemOutput struct {
	Problem ProblemDTO
}

type GetProblemsOutput struct {
	Problems []ProblemDTO
}

type GetProblemBySlugInput struct {
	Slug string
}

type GetProblemBySlugOutput struct {
	Problem ProblemDTO
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	Token string
}

type RegisterInput struct {
	Username string
	Email    string
	Password string
}

type RegisterOutput struct {
	Token string
}

type UserDTO struct {
	ID       int
	Username string
	Email    string
	Role     string
}

type CreateSubmissionInput struct {
	UserID    int
	ProblemID int
	Code      string
}

type CreateSubmissionOutput struct {
	SubmissionID int
}

type UpdateSubmissionInput struct {
	SubmissionID int
	Status       domain.SubmissionStatus
	Result       domain.SubmissionResult
}

type UpdateSubmissionOutput struct{}

type GetSubmissionByIdInput struct {
	ID int
}

type SubmissionDTO struct {
	ID     int
	Code   string
	Status string
	Result string
}

type GetSubmissionByIdOutput struct {
	Submission SubmissionDTO
}

type CreateContestInput struct {
	CreatorID   int
	Title       string
	Slug        string
	Description string
	StartsAt    time.Time
	EndsAt      time.Time
	Problems    []CreateContestProblemInput
}

type CreateContestProblemInput struct {
	Title      string
	Slug       string
	Difficulty string
	Content    string
	TestCases  []CreateContestTestCaseInput
}

type CreateContestTestCaseInput struct {
	Input    string
	Expected string
}

type ContestDTO struct {
	ID          int
	CreatorID   int
	Title       string
	Slug        string
	Description string
	StartsAt    time.Time
	EndsAt      time.Time
	Problems    []ProblemDTO
}

type CreateContestOutput struct {
	Contest ContestDTO
}

type GetContestBySlugInput struct {
	Slug string
}

type GetContestBySlugOutput struct {
	Contest ContestDTO
}
