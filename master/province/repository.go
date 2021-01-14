package province

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*models.Province, error)
	Update(ctx context.Context, ar *models.Province) error
	Insert(ctx context.Context, a *models.Province) error
	Delete(ctx context.Context, id string, deleted_by string) error
	Count(ctx context.Context) (int, error)
	List(ctx context.Context, limit, offset int) ([]*models.Province, error)
}
