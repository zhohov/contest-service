package usecase

import (
	"github.com/zhohov/contest-service/backend/internal/domain"
)

func (uc *UseCase) CreateArticle(in CreateArticleInput) (*CreateArticleOutput, error) {
	article := domain.NewArticle(in.Title, in.Slug, in.Content)
	createdArticle, err := uc.articles.CreateArticle(article)
	if err != nil {
		return nil, err
	}

	return &CreateArticleOutput{
		Article: mapArticleToDTO(createdArticle),
	}, nil
}

func (uc *UseCase) GetArticles() (*GetArticlesOutput, error) {
	articles, err := uc.articles.GetArticles()
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
