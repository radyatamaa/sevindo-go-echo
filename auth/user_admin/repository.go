package user_admin

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*models.UserAdmin, error)
	GetByUserEmail(ctx context.Context, userEmail string,isAdmin bool) (*models.UserAdmin, error)
	Update(ctx context.Context, ar *models.UserAdmin) error
	Insert(ctx context.Context, a *models.UserAdmin) error
	Delete(ctx context.Context, id string, deleted_by string) error
	Count(ctx context.Context) (int, error)
	List(ctx context.Context, limit, offset int, search string) ([]*models.UserAdmin, error)
}
