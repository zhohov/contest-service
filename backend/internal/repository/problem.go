package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zhohov/contest-service/backend/internal/domain"
)

type ProblemRepository struct {
	db *sql.DB
}

func NewProblemRepository(db *sql.DB) *ProblemRepository {
	return &ProblemRepository{
		db: db,
	}
}

func (r *ProblemRepository) CreateProblem(ctx context.Context, problem *domain.Problem) (*domain.Problem, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("create problem repo start tx: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO problems (title, slug, difficulty, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	args := []any{problem.Title, problem.Slug, problem.Difficulty, problem.Content}

	var problemID int
	err = tx.QueryRowContext(ctx, query, args...).Scan(&problemID)
	if err != nil {
		return nil, fmt.Errorf("create problem from repo: %w", err)
	}

	queryTestCase := `
		INSERT INTO test_cases (problem_id, input, expected)
		VALUES ($1, $2, $3);
	`
	for _, testCase := range problem.TestCases {
		args = []any{problemID, testCase.Input, testCase.Expected}

		_, err := tx.ExecContext(ctx, queryTestCase, args...)
		if err != nil {
			return nil, fmt.Errorf("create test case from repo: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("create problem repo commit tx: %w", err)
	}

	problem.ID = problemID

	return problem, nil
}

func (r *ProblemRepository) GetProblems(ctx context.Context) ([]domain.Problem, error) {
	query := `
		SELECT id, title, slug, difficulty, content
		FROM problems;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repo get problems: %w", err)
	}
	defer rows.Close()

	problems := []domain.Problem{}
	for rows.Next() {
		problem := domain.Problem{}

		err := rows.Scan(&problem.ID, &problem.Title, &problem.Slug, &problem.Difficulty, &problem.Content)
		if err != nil {
			return nil, fmt.Errorf("repo scan problem: %w", err)
		}

		problems = append(problems, problem)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repo rows error: %w", err)
	}

	return problems, nil
}

func (r *ProblemRepository) GetProblemById(ctx context.Context, id int) (*domain.Problem, error) {
	query := `
		SELECT id, title, slug, difficulty, content
		FROM problems
		WHERE id = $1;
	`

	problem, err := r.getProblem(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("get problem from repo: %w", err)
	}

	testCases, err := r.getTestCases(ctx, problem.ID)
	if err != nil {
		return nil, fmt.Errorf("get test cases from repo: %w", err)
	}

	problem.TestCases = testCases

	return problem, nil
}

func (r *ProblemRepository) GetProblemBySlug(ctx context.Context, slug string) (*domain.Problem, error) {
	query := `
		SELECT id, title, slug, difficulty, content
		FROM problems
		WHERE slug = $1;
	`

	problem, err := r.getProblem(ctx, query, slug)
	if err != nil {
		return nil, fmt.Errorf("get problem from repo: %w", err)
	}

	testCases, err := r.getTestCases(ctx, problem.ID)
	if err != nil {
		return nil, fmt.Errorf("get test cases from repo: %w", err)
	}

	problem.TestCases = testCases

	return problem, nil
}

func (r *ProblemRepository) getProblem(ctx context.Context, query string, args ...any) (*domain.Problem, error) {
	row := r.db.QueryRowContext(ctx, query, args...)

	problem, err := r.scanProblem(row)
	if err != nil {
		return nil, fmt.Errorf("scan problem from repo: %w", err)
	}

	return problem, nil
}

func (r *ProblemRepository) getTestCases(ctx context.Context, problemID int) ([]domain.TestCase, error) {
	query := `
		SELECT id, input, expected
		FROM test_cases
		WHERE problem_id = $1;
	`

	rows, err := r.db.QueryContext(ctx, query, problemID)
	if err != nil {
		return nil, fmt.Errorf("get test cases from repo: %w", err)
	}
	defer rows.Close()

	testCases := []domain.TestCase{}
	for rows.Next() {
		testCase := domain.TestCase{}
		err := rows.Scan(&testCase.ID, &testCase.Input, &testCase.Expected)
		if err != nil {
			return nil, fmt.Errorf("scan test cases from repo: %w", err)
		}

		testCases = append(testCases, testCase)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repo rows error: %w", err)
	}

	return testCases, nil
}

func (r *ProblemRepository) scanProblem(row *sql.Row) (*domain.Problem, error) {
	problem := &domain.Problem{}
	err := row.Scan(
		&problem.ID,
		&problem.Title,
		&problem.Slug,
		&problem.Difficulty,
		&problem.Content,
	)
	if err != nil {
		return nil, fmt.Errorf("scan problem from repo: %w", err)
	}

	return problem, nil
}
