package city

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Create(ctx context.Context, ar *models.NewCommandCity, token string) (*models.NewCommandCity, error)
	GetById(ctx context.Context, id int, token string) (*models.CityDto, error)
	Delete(ctx context.Context, id int, token string) (*models.ResponseDelete, error)
	Update(ctx context.Context, ar *models.NewCommandCity,  token string) error
	List(ctx context.Context, page, limit, offset int, search string) (*models.CityWithPagination, error)
}
