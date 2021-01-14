package province

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Create(ctx context.Context, ar *models.NewCommandProvince, token string) (*models.NewCommandProvince, error)
	GetById(ctx context.Context, id int, token string) (*models.ProvinceDto, error)
	Delete(ctx context.Context, id int, token string) (*models.ResponseDelete, error)
	Update(ctx context.Context, ar *models.NewCommandProvince,  token string) error
	List(ctx context.Context, page, limit, offset int, search string) (*models.ProvinceWithPagination, error)
}
