package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zhohov/contest-service/backend/internal/domain"
)

type ContestRepository struct {
	db *sql.DB
}

func NewContestRepository(db *sql.DB) *ContestRepository {
	return &ContestRepository{
		db: db,
	}
}

func (r *ContestRepository) CreateContest(ctx context.Context, contest *domain.Contest) (*domain.Contest, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("create contest repo start tx: %w", err)
	}
	defer tx.Rollback()

	queryContest := `
		INSERT INTO contests (creator_id, title, slug, description, starts_at, ends_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`

	err = tx.QueryRowContext(
		ctx,
		queryContest,
		contest.Creator.ID,
		contest.Title,
		contest.Slug,
		contest.Description,
		contest.StartsAt,
		contest.EndsAt,
	).Scan(&contest.ID)
	if err != nil {
		return nil, fmt.Errorf("create contest from repo: %w", err)
	}

	queryProblem := `
		INSERT INTO problems (title, slug, difficulty, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	queryTestCase := `
		INSERT INTO test_cases (problem_id, input, expected)
		VALUES ($1, $2, $3);
	`

	queryContestProblem := `
		INSERT INTO contest_problems (contest_id, problem_id)
		VALUES ($1, $2);
	`

	for i := range contest.Problems {
		problem := &contest.Problems[i]

		err := tx.QueryRowContext(
			ctx,
			queryProblem,
			problem.Title,
			problem.Slug,
			problem.Difficulty,
			problem.Content,
		).Scan(&problem.ID)
		if err != nil {
			return nil, fmt.Errorf("create contest problem from repo: %w", err)
		}

		for _, testCase := range problem.TestCases {
			_, err := tx.ExecContext(
				ctx,
				queryTestCase,
				problem.ID,
				testCase.Input,
				testCase.Expected,
			)
			if err != nil {
				return nil, fmt.Errorf("create contest problem test case from repo: %w", err)
			}
		}

		_, err = tx.ExecContext(
			ctx,
			queryContestProblem,
			contest.ID,
			problem.ID,
		)
		if err != nil {
			return nil, fmt.Errorf("create contest problem relation from repo: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("create contest repo commit tx: %w", err)
	}

	return contest, nil
}

func (r *ContestRepository) GetContestByID(ctx context.Context, id int) (*domain.Contest, error) {
	query := `
		SELECT id, creator_id, title, slug, description, starts_at, ends_at, created_at
		FROM contests
		WHERE id = $1;
	`

	row := r.db.QueryRowContext(ctx, query, id)

	contest, err := r.scanContest(row)
	if err != nil {
		return nil, fmt.Errorf("repo get contest: %w", err)
	}

	return contest, nil
}

func (r *ContestRepository) GetContests(ctx context.Context) ([]domain.Contest, error) {
	query := `
		SELECT id, creator_id, title, slug, description, starts_at, ends_at, created_at
		FROM contests
		ORDER BY created_at DESC;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repo get contests: %w", err)
	}
	defer rows.Close()

	contests := []domain.Contest{}
	for rows.Next() {
		contest := domain.Contest{}

		err := rows.Scan(
			&contest.ID,
			&contest.Creator.ID,
			&contest.Title,
			&contest.Slug,
			&contest.Description,
			&contest.StartsAt,
			&contest.EndsAt,
			&contest.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan contest: %w", err)
		}

		contests = append(contests, contest)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("contest rows error: %w", err)
	}

	return contests, nil
}

func (r *ContestRepository) GetContestBySlug(ctx context.Context, slug string) (*domain.Contest, error) {
	query := `
		SELECT id, creator_id, title, slug, description, starts_at, ends_at, created_at
		FROM contests
		WHERE slug = $1;
	`

	row := r.db.QueryRowContext(ctx, query, slug)

	contest, err := r.scanContest(row)
	if err != nil {
		return nil, fmt.Errorf("repo get contest by slug: %w", err)
	}

	problems, err := r.getContestProblems(ctx, contest.ID)
	if err != nil {
		return nil, fmt.Errorf("repo get contest problems: %w", err)
	}

	contest.Problems = problems

	return contest, nil
}

func (r *ContestRepository) getContestProblems(ctx context.Context, contestID int) ([]domain.Problem, error) {
	query := `
		SELECT p.id, p.title, p.slug, p.difficulty, p.content
		FROM problems p
		JOIN contest_problems cp
			ON cp.problem_id = p.id
		WHERE cp.contest_id = $1
		ORDER BY p.id;
	`

	rows, err := r.db.QueryContext(ctx, query, contestID)
	if err != nil {
		return nil, fmt.Errorf("get contest problems: %w", err)
	}
	defer rows.Close()

	problems := []domain.Problem{}
	for rows.Next() {
		problem := domain.Problem{}

		err := rows.Scan(
			&problem.ID,
			&problem.Title,
			&problem.Slug,
			&problem.Difficulty,
			&problem.Content,
		)
		if err != nil {
			return nil, fmt.Errorf("scan contest problem: %w", err)
		}

		problems = append(problems, problem)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("contest problems rows error: %w", err)
	}

	return problems, nil
}

func (r *ContestRepository) UpdateContest(ctx context.Context, contest *domain.Contest) error {
	query := `
		UPDATE contests
		SET
			creator_id = $1,
			title = $2,
			slug = $3,
			description = $4,
			starts_at = $5,
			ends_at = $6
		WHERE id = $7;
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		contest.Creator.ID,
		contest.Title,
		contest.Slug,
		contest.Description,
		contest.StartsAt,
		contest.EndsAt,
		contest.ID,
	)
	if err != nil {
		return fmt.Errorf("update contest from repo: %w", err)
	}

	return nil
}

func (r *ContestRepository) scanContest(row *sql.Row) (*domain.Contest, error) {
	contest := &domain.Contest{}

	err := row.Scan(
		&contest.ID,
		&contest.Creator.ID,
		&contest.Title,
		&contest.Slug,
		&contest.Description,
		&contest.StartsAt,
		&contest.EndsAt,
		&contest.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan contest: %w", err)
	}

	return contest, nil
}
