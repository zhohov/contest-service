package repository

import (
	"context"
	"errors"

	"github.com/zhohov/contest-service/backend/internal/domain"
)

var ErrNotFound = errors.New("Not Found")

type ArticleInMemoryRepo struct {
	storage []domain.Article
}

func NewArticleInMemoryRepo(storage []domain.Article) *ArticleInMemoryRepo {
	return &ArticleInMemoryRepo{
		storage: storage,
	}
}

func (r *ArticleInMemoryRepo) CreateArticle(ctx context.Context, article *domain.Article) (*domain.Article, error) {
	r.storage = append(r.storage, *article)
	return article, nil
}

func (r *ArticleInMemoryRepo) GetArticles(ctx context.Context) ([]domain.Article, error) {
	return r.storage, nil
}

func (r *ArticleInMemoryRepo) GetArticleById(ctx context.Context, id int) (*domain.Article, error) {
	for i := range r.storage {
		if r.storage[i].ID == id {
			return &r.storage[i], nil
		}
	}

	return nil, ErrNotFound
}
