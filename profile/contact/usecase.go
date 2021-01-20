package contact

import (
	"context"
	"github.com/models"
)

type Usecase interface {
	Delete(ctx context.Context, id string, token string) (*models.ResponseDelete, error)
	Update(ctx context.Context, ar *models.NewCommandContact,  token string) error
	List(ctx context.Context, page, limit, offset int, token string) (*models.ContactWithPagination, error)
	Create(ctx context.Context, ar *models.NewCommandContact, token string) (*models.NewCommandContact, error)
	GetById(ctx context.Context, id string, token string) (*models.ContactDto, error)
}
