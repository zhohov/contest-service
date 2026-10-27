package usecase

import (
	"context"

	"github.com/zhohov/contest-service/backend/internal/domain"
)

type ArticleRepository interface {
	CreateArticle(ctx context.Context, article *domain.Article) (*domain.Article, error)
	GetArticles(ctx context.Context) ([]domain.Article, error)
	GetArticleById(ctx context.Context, id int) (*domain.Article, error)
}

type ProblemRepository interface {
	CreateProblem(ctx context.Context, problem *domain.Problem) (*domain.Problem, error)
	GetProblems(ctx context.Context) ([]domain.Problem, error)
	GetProblemById(ctx context.Context, id int) (*domain.Problem, error)
	GetProblemBySlug(ctx context.Context, slug string) (*domain.Problem, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUsers(ctx context.Context) ([]domain.User, error)
	GetUserById(ctx context.Context, id int) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type SubmissionRepository interface {
	CreateSubmission(ctx context.Context, submission *domain.Submission) (*domain.Submission, error)
	GetSubmissionByID(ctx context.Context, id int) (*domain.Submission, error)
	GetSubmissionsByUserID(ctx context.Context, userID int) ([]domain.Submission, error)
	GetSubmissionsByContestID(ctx context.Context, contestID int) ([]domain.Submission, error)
	UpdateSubmission(ctx context.Context, submission *domain.Submission) error
}

type ContestRepository interface {
	CreateContest(ctx context.Context, contest *domain.Contest) (*domain.Contest, error)
	GetContestByID(ctx context.Context, id int) (*domain.Contest, error)
	GetContestsByUserID(ctx context.Context, userID int) ([]domain.Contest, error)
	GetContestBySlug(ctx context.Context, slug string) (*domain.Contest, error)
	UpdateContest(ctx context.Context, contest *domain.Contest) error
}

type JwtService interface {
	GenerateToken(userID int, role string) (string, error)
	ParseToken(token string) (*domain.AuthUser, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password, hash string) bool
}

type MessageBroker interface {
	Publish(ctx context.Context, topic string, message any) error
	Subscribe(ctx context.Context, topic string) (<-chan []byte, error)
}

type UseCase struct {
	articles    ArticleRepository
	problems    ProblemRepository
	users       UserRepository
	submissions SubmissionRepository
	contests    ContestRepository
	jwt         JwtService
	hasher      PasswordHasher
	broker      MessageBroker
}

func NewUseCase(
	articles ArticleRepository,
	problems ProblemRepository,
	users UserRepository,
	submissions SubmissionRepository,
	contests ContestRepository,
	jwt JwtService,
	hasher PasswordHasher,
	broker MessageBroker,
) *UseCase {
	return &UseCase{
		articles:    articles,
		problems:    problems,
		users:       users,
		submissions: submissions,
		contests:    contests,
		jwt:         jwt,
		hasher:      hasher,
		broker:      broker,
	}
}
