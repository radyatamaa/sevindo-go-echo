package role

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Delete(ctx context.Context, id int, token string) (*models.ResponseDelete, error)
	Update(ctx context.Context, ar *models.NewCommandRole,  token string) error
	List(ctx context.Context, page, limit, offset int, search string) (*models.RoleWithPagination, error)
	Create(ctx context.Context, ar *models.NewCommandRole, token string) (*models.NewCommandRole, error)
	GetById(ctx context.Context, id int, token string) (*models.RoleDto, error)
}
