package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zhohov/contest-service/backend/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
		INSERT INTO users (username, email, role, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, email, role, password_hash;
	`

	args := []any{user.Username, user.Email, user.Role, user.PasswordHash}

	row := r.db.QueryRowContext(ctx, query, args...)

	createdUser, err := r.scanUser(row)
	if err != nil {
		return nil, fmt.Errorf("repo create user: %w", err)
	}

	return createdUser, nil
}

func (r *UserRepository) GetUsers(ctx context.Context) ([]domain.User, error) {
	query := `
		SELECT id, username, email, role, password_hash
		FROM users;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repo get users: %w", err)
	}
	defer rows.Close()

	users := []domain.User{}
	for rows.Next() {
		user := domain.User{}

		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.PasswordHash)
		if err != nil {
			return nil, fmt.Errorf("repo scan user: %w", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repo rows error: %w", err)
	}

	return users, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, id int) (*domain.User, error) {
	query := `
		SELECT id, username, email, role, password_hash
		FROM users
		WHERE id = $1;
	`

	row := r.db.QueryRowContext(ctx, query, id)

	user, err := r.scanUser(row)
	if err != nil {
		return nil, fmt.Errorf("repo get user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, username, email, role, password_hash
		FROM users
		WHERE email = $1;
	`

	row := r.db.QueryRowContext(ctx, query, email)

	user, err := r.scanUser(row)
	if err != nil {
		return nil, fmt.Errorf("repo get user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) scanUser(row *sql.Row) (*domain.User, error) {
	user := &domain.User{}

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("repo scan user: %w", err)
	}

	return user, nil
}
