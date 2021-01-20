package contact

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*models.Contact, error)
	Update(ctx context.Context, ar *models.Contact) error
	Insert(ctx context.Context, a *models.Contact) error
	Delete(ctx context.Context, id string, deleted_by string) error
	Count(ctx context.Context,userId string) (int, error)
	List(ctx context.Context, limit, offset int,userId string) ([]*models.Contact, error)
}
