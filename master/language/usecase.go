package language

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Delete(ctx context.Context, id int, token string) (*models.ResponseDelete, error)
	Update(ctx context.Context, ar *models.NewCommandLanguage, token string) error
	List(ctx context.Context, page, limit, offset int, search string) (*models.LanguageWithPagination, error)
	Create(ctx context.Context, ar *models.NewCommandLanguage, token string) (*models.NewCommandLanguage, error)
	GetById(ctx context.Context, id int, token string) (*models.LanguageDto, error)
}
