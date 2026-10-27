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

func (r *ContestRepository) GetContestsByUserID(ctx context.Context, userID int) ([]domain.Contest, error) {
	query := `
		SELECT id, creator_id, title, slug, description, starts_at, ends_at, created_at
		FROM contests
		WHERE creator_id = $1
		ORDER BY created_at DESC;
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("repo get contests by user id: %w", err)
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
	rows.Close()

	for i := range problems {
		testCases, err := r.getContestProblemTestCases(ctx, problems[i].ID)
		if err != nil {
			return nil, fmt.Errorf("get contest problem test cases: %w", err)
		}

		problems[i].TestCases = testCases
	}

	return problems, nil
}

func (r *ContestRepository) getContestProblemTestCases(ctx context.Context, problemID int) ([]domain.TestCase, error) {
	query := `
		SELECT id, input, expected
		FROM test_cases
		WHERE problem_id = $1
		ORDER BY id;
	`

	rows, err := r.db.QueryContext(ctx, query, problemID)
	if err != nil {
		return nil, fmt.Errorf("get contest problem test cases: %w", err)
	}
	defer rows.Close()

	testCases := []domain.TestCase{}
	for rows.Next() {
		testCase := domain.TestCase{}
		if err := rows.Scan(&testCase.ID, &testCase.Input, &testCase.Expected); err != nil {
			return nil, fmt.Errorf("scan contest problem test case: %w", err)
		}

		testCases = append(testCases, testCase)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("contest problem test cases rows error: %w", err)
	}

	return testCases, nil
}

func (r *ContestRepository) UpdateContest(ctx context.Context, contest *domain.Contest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("update contest repo start tx: %w", err)
	}
	defer tx.Rollback()

	queryContest := `
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

	result, err := tx.ExecContext(
		ctx,
		queryContest,
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

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update contest rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("update contest from repo: contest not found")
	}

	queryCurrentProblems := `
		SELECT problem_id
		FROM contest_problems
		WHERE contest_id = $1
		FOR UPDATE;
	`

	rows, err := tx.QueryContext(ctx, queryCurrentProblems, contest.ID)
	if err != nil {
		return fmt.Errorf("get contest problems for update: %w", err)
	}

	currentProblemIDs := make(map[int]struct{})
	for rows.Next() {
		var problemID int
		if err := rows.Scan(&problemID); err != nil {
			rows.Close()
			return fmt.Errorf("scan contest problem for update: %w", err)
		}
		currentProblemIDs[problemID] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("contest problems for update rows error: %w", err)
	}
	rows.Close()

	queryUpdateProblem := `
		UPDATE problems
		SET title = $1, slug = $2, difficulty = $3, content = $4
		WHERE id = $5;
	`
	queryCreateProblem := `
		INSERT INTO problems (title, slug, difficulty, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	queryCreateRelation := `
		INSERT INTO contest_problems (contest_id, problem_id)
		VALUES ($1, $2);
	`
	queryDeleteTestCases := `
		DELETE FROM test_cases
		WHERE problem_id = $1;
	`
	queryCreateTestCase := `
		INSERT INTO test_cases (problem_id, input, expected)
		VALUES ($1, $2, $3);
	`

	keptProblemIDs := make(map[int]struct{}, len(contest.Problems))
	for i := range contest.Problems {
		problem := &contest.Problems[i]

		if problem.ID == 0 {
			err := tx.QueryRowContext(
				ctx,
				queryCreateProblem,
				problem.Title,
				problem.Slug,
				problem.Difficulty,
				problem.Content,
			).Scan(&problem.ID)
			if err != nil {
				return fmt.Errorf("create contest problem on update: %w", err)
			}

			if _, err := tx.ExecContext(ctx, queryCreateRelation, contest.ID, problem.ID); err != nil {
				return fmt.Errorf("create contest problem relation on update: %w", err)
			}
		} else {
			if _, ok := currentProblemIDs[problem.ID]; !ok {
				return fmt.Errorf("update contest problem: problem %d does not belong to contest", problem.ID)
			}

			if _, exists := keptProblemIDs[problem.ID]; exists {
				return fmt.Errorf("update contest problem: duplicate problem id %d", problem.ID)
			}

			if _, err := tx.ExecContext(
				ctx,
				queryUpdateProblem,
				problem.Title,
				problem.Slug,
				problem.Difficulty,
				problem.Content,
				problem.ID,
			); err != nil {
				return fmt.Errorf("update contest problem: %w", err)
			}

			if _, err := tx.ExecContext(ctx, queryDeleteTestCases, problem.ID); err != nil {
				return fmt.Errorf("delete contest problem test cases on update: %w", err)
			}
		}

		for _, testCase := range problem.TestCases {
			if _, err := tx.ExecContext(
				ctx,
				queryCreateTestCase,
				problem.ID,
				testCase.Input,
				testCase.Expected,
			); err != nil {
				return fmt.Errorf("create contest problem test case on update: %w", err)
			}
		}

		keptProblemIDs[problem.ID] = struct{}{}
	}

	queryDeleteRelation := `
		DELETE FROM contest_problems
		WHERE contest_id = $1 AND problem_id = $2;
	`
	for problemID := range currentProblemIDs {
		if _, ok := keptProblemIDs[problemID]; ok {
			continue
		}

		if _, err := tx.ExecContext(ctx, queryDeleteRelation, contest.ID, problemID); err != nil {
			return fmt.Errorf("delete contest problem relation on update: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("update contest repo commit tx: %w", err)
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
