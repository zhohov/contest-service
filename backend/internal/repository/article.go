package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zhohov/contest-service/backend/internal/domain"
)

type ArticleRepository struct {
	db *sql.DB
}

func NewArticleRepository(db *sql.DB) *ArticleRepository {
	return &ArticleRepository{
		db: db,
	}
}

func (r *ArticleRepository) CreateArticle(ctx context.Context, article *domain.Article) (*domain.Article, error) {
	query := `
		INSERT INTO articles (title, slug, content, author_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		article.Title,
		article.Slug,
		article.Content,
		article.Author.ID,
	).Scan(&article.ID)
	if err != nil {
		return nil, fmt.Errorf("create article from repo: %w", err)
	}

	return article, nil
}

func (r *ArticleRepository) GetArticles(ctx context.Context) ([]domain.Article, error) {
	query := `
		SELECT a.id, a.author_id, u.username, u.email, u.role, a.title, a.slug, a.content
		FROM articles a
		JOIN users u
			ON u.id = a.author_id
		ORDER BY a.created_at DESC;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repo get articles: %w", err)
	}
	defer rows.Close()

	articles := []domain.Article{}
	for rows.Next() {
		article := domain.Article{}

		err := rows.Scan(
			&article.ID,
			&article.Author.ID,
			&article.Author.Username,
			&article.Author.Email,
			&article.Author.Role,
			&article.Title,
			&article.Slug,
			&article.Content,
		)
		if err != nil {
			return nil, fmt.Errorf("repo scan article: %w", err)
		}

		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("article rows error: %w", err)
	}

	return articles, nil
}

func (r *ArticleRepository) GetArticleById(ctx context.Context, id int) (*domain.Article, error) {
	query := `
		SELECT a.id, a.author_id, u.username, u.email, u.role, a.title, a.slug, a.content
		FROM articles a
		JOIN users u
			ON u.id = a.author_id
		WHERE a.id = $1;
	`

	row := r.db.QueryRowContext(ctx, query, id)

	article, err := r.scanArticle(row)
	if err != nil {
		return nil, fmt.Errorf("repo get article: %w", err)
	}

	return article, nil
}

func (r *ArticleRepository) scanArticle(row *sql.Row) (*domain.Article, error) {
	article := &domain.Article{}

	err := row.Scan(
		&article.ID,
		&article.Author.ID,
		&article.Author.Username,
		&article.Author.Email,
		&article.Author.Role,
		&article.Title,
		&article.Slug,
		&article.Content,
	)
	if err != nil {
		return nil, fmt.Errorf("scan article: %w", err)
	}

	return article, nil
}
