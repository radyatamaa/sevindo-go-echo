package article_category

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Delete(ctx context.Context, id int, token string) (*models.ResponseDelete, error)
	Update(ctx context.Context, ar *models.NewCommandArticleCategory, token string) error
	List(ctx context.Context, page, limit, offset int, search string) (*models.ArticleCategoryWithPagination, error)
	Create(ctx context.Context, ar *models.NewCommandArticleCategory, token string) (*models.NewCommandArticleCategory, error)
	GetById(ctx context.Context, id int, token string) (*models.ArticleCategoryDto, error)
}
