package promo

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Create(ctx context.Context, ar *models.NewCommandPromo, token string) (*models.NewCommandPromo, error)
	GetById(ctx context.Context, id string, token string) (*models.PromoDto, error)
	Delete(ctx context.Context, id string, token string) (*models.ResponseDelete, error)
	Update(ctx context.Context, ar *models.NewCommandPromo,  token string) error
	List(ctx context.Context, page, limit, offset int, search string) (*models.PromoWithPagination, error)
}
