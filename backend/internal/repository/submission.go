package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zhohov/contest-service/backend/internal/domain"
)

type SubmissionRepository struct {
	db *sql.DB
}

func NewSubmissionRepository(db *sql.DB) *SubmissionRepository {
	return &SubmissionRepository{
		db: db,
	}
}

func (r *SubmissionRepository) CreateSubmission(ctx context.Context, submission *domain.Submission) (*domain.Submission, error) {
	query := `
		INSERT INTO submissions (problem_id, user_id, code, status, result)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	args := []any{submission.Problem.ID, submission.User.ID, submission.Code, submission.Status, submission.Result}

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&submission.ID)
	if err != nil {
		return nil, fmt.Errorf("create submission from repo: %w", err)
	}

	return submission, nil
}

func (r *SubmissionRepository) GetSubmissionByID(ctx context.Context, id int) (*domain.Submission, error) {
	query := `
		SELECT id, problem_id, user_id, code, status, result
		FROM submissions
		WHERE id = $1;
	`

	row := r.db.QueryRowContext(ctx, query, id)

	submission, err := r.scanSubmission(row)
	if err != nil {
		return nil, fmt.Errorf("repo get submission: %w", err)
	}

	return submission, nil
}

func (r *SubmissionRepository) GetSubmissionsByUserID(ctx context.Context, userID int) ([]domain.Submission, error) {
	query := `
		SELECT id, problem_id, user_id, code, status, result
		FROM submissions
		WHERE user_id = $1
		ORDER BY created_at DESC;
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("repo get submissions by user id: %w", err)
	}
	defer rows.Close()

	submissions := []domain.Submission{}
	for rows.Next() {
		submission := domain.Submission{}

		err := rows.Scan(
			&submission.ID,
			&submission.Problem.ID,
			&submission.User.ID,
			&submission.Code,
			&submission.Status,
			&submission.Result,
		)
		if err != nil {
			return nil, fmt.Errorf("scan submission: %w", err)
		}

		submissions = append(submissions, submission)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("submission rows error: %w", err)
	}

	return submissions, nil
}

func (r *SubmissionRepository) GetSubmissionsByContestID(ctx context.Context, contestID int) ([]domain.Submission, error) {
	query := `
		SELECT s.id, s.problem_id, s.user_id, s.code, s.status, s.result
		FROM submissions s
		JOIN contest_problems cp
			ON cp.problem_id = s.problem_id
		WHERE cp.contest_id = $1
		ORDER BY s.created_at DESC;
	`

	rows, err := r.db.QueryContext(ctx, query, contestID)
	if err != nil {
		return nil, fmt.Errorf("repo get submissions by contest id: %w", err)
	}
	defer rows.Close()

	submissions := []domain.Submission{}
	for rows.Next() {
		submission := domain.Submission{}

		err := rows.Scan(
			&submission.ID,
			&submission.Problem.ID,
			&submission.User.ID,
			&submission.Code,
			&submission.Status,
			&submission.Result,
		)
		if err != nil {
			return nil, fmt.Errorf("scan submission: %w", err)
		}

		submissions = append(submissions, submission)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("submission rows error: %w", err)
	}

	return submissions, nil
}

func (r *SubmissionRepository) UpdateSubmission(ctx context.Context, submission *domain.Submission) error {
	query := `
		UPDATE submissions
		SET
			problem_id = $1,
			user_id = $2,
			code = $3,
			status = $4,
			result = $5
		WHERE id = $6;
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		submission.Problem.ID,
		submission.User.ID,
		submission.Code,
		submission.Status,
		submission.Result,
		submission.ID,
	)
	if err != nil {
		return fmt.Errorf("update submission from repo: %w", err)
	}

	return nil
}

func (r *SubmissionRepository) scanSubmission(row *sql.Row) (*domain.Submission, error) {
	submission := &domain.Submission{}
	err := row.Scan(
		&submission.ID,
		&submission.Problem.ID,
		&submission.User.ID,
		&submission.Code,
		&submission.Status,
		&submission.Result,
	)
	if err != nil {
		return nil, fmt.Errorf("scan submission: %w", err)
	}

	return submission, nil
}
