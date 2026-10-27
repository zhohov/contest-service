package usecase

import (
	"context"
	"fmt"

	"github.com/zhohov/contest-service/backend/internal/domain"
)

func (uc *UseCase) CreateArticle(ctx context.Context, in CreateArticleInput) (*CreateArticleOutput, error) {
	user, err := uc.users.GetUserById(ctx, in.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("get article author from uc: %w", err)
	}

	article := domain.NewArticle(in.Title, in.Slug, in.Content)
	article.Author = *user

	createdArticle, err := uc.articles.CreateArticle(ctx, article)
	if err != nil {
		return nil, err
	}

	return &CreateArticleOutput{
		Article: mapArticleToDTO(createdArticle),
	}, nil
}

func (uc *UseCase) GetArticles(ctx context.Context) (*GetArticlesOutput, error) {
	articles, err := uc.articles.GetArticles(ctx)
	if err != nil {
		return nil, err
	}

	return &GetArticlesOutput{
		Articles: mapArticlesToDTOs(articles),
	}, nil
}

// func (uc *UseCase) GetArticleBySlug(in GetArticleBySlug) (*GetArticleBySlug, error) {
// 	articles, err := uc.articles.GetArticleBySlug(in.Slug)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &GetArticlesOutput{
// 		Articles: mapArticlesToDTOs(articles),
// 	}, nil
// }
