package accessibility

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	GetByID(ctx context.Context, id int) (*models.Accessibility, error)
	Update(ctx context.Context, ar *models.Accessibility) error
	Insert(ctx context.Context, a *models.Accessibility) error
	Delete(ctx context.Context, id int, deleted_by string) error
	Count(ctx context.Context) (int, error)
	List(ctx context.Context, limit, offset int) ([]*models.Accessibility, error)
}
