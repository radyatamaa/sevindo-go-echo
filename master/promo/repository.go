package promo

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*models.Promo, error)
	Update(ctx context.Context, ar *models.Promo) error
	Insert(ctx context.Context, a *models.Promo) error
	Delete(ctx context.Context, id string, deleted_by string) error
	Count(ctx context.Context) (int, error)
	List(ctx context.Context, limit, offset int) ([]*models.Promo, error)
}
