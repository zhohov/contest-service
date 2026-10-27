package domain

import "fmt"

type Article struct {
	ID      int
	Author  User
	Title   string
	Slug    string
	Content string
}

func NewArticle(title, slug, content string) *Article {
	return &Article{
		Title:   title,
		Slug:    slug,
		Content: content,
	}
}

func (a *Article) String() string {
	return fmt.Sprintf("Article(id=%d, title=%s)", a.ID, a.Title)
}
