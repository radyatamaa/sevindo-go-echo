package article_category

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Create(ctx context.Context, ar *models.NewCommandArticleCategory, token string) (*models.NewCommandArticleCategory, error)
	GetById(ctx context.Context, id string, token string) (*models.ArticleCategoryDto, error)
}
