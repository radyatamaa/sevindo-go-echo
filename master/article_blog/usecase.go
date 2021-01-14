package article_blog

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Delete(ctx context.Context, id int, token string) (*models.ResponseDelete, error)
	Update(ctx context.Context, ar *models.NewCommandArticleBlog, token string) error
	List(ctx context.Context, page, limit, offset int, search string) (*models.ArticleBlogWithPagination, error)
	Create(ctx context.Context, ar *models.NewCommandArticleBlog, token string) (*models.NewCommandArticleBlog, error)
	GetById(ctx context.Context, id int, token string) (*models.ArticleBlogDto, error)
}
