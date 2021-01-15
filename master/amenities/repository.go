package amenities

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	GetByID(ctx context.Context, id int) (*models.Amenities, error)
	Update(ctx context.Context, ar *models.Amenities) error
	Insert(ctx context.Context, a *models.Amenities) error
	Delete(ctx context.Context, id int, deleted_by string) error
	Count(ctx context.Context) (int, error)
	List(ctx context.Context, limit, offset int) ([]*models.Amenities, error)
}
