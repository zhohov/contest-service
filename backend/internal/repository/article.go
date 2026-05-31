package repository

// import (
// 	"context"
// 	"database/sql"
// 	"log"

// 	_ "github.com/lib/pq"

// 	"github.com/zhohov/judge-service/internal/domain"
// )

// type ArticleRepository struct {
// 	db *sql.DB
// }

// func NewArticleRepository(db *sql.DB) *ArticleRepository {
// 	return &ArticleRepository{
// 		db: db,
// 	}
// }

// func (r *ArticleRepository) Create(ctx context.Context) (*domain.Article, error) {
// 	query := `
// 		INSERT INTO users (username, email, password_hash, role)
// 		VALUES ($1, $2, $3, $4)
// 		RETURNING id;
// 	`

// 	res, err := r.db.ExecContext(ctx, query)
// 	if err != nil {
// 		return nil, err
// 	}

// 	log.Println(res)

// 	return nil, nil
// }

// func (r *ArticleRepository) List() ([]domain.Article, error) {

// }

// func (r *ArticleRepository) GetByID(id int) (*domain.Article, error) {

// }

// func (r *ArticleRepository) GetBySlug(slug string) (*domain.Article, error) {

// }
