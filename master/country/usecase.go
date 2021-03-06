package country

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Delete(ctx context.Context, id int, token string) (*models.ResponseDelete, error)
	Update(ctx context.Context, ar *models.NewCommandCountry,  token string) error
	List(ctx context.Context, page, limit, offset int, search string) (*models.CountryWithPagination, error)
	Create(ctx context.Context, ar *models.NewCommandCountry, token string) (*models.NewCommandCountry, error)
	GetById(ctx context.Context, id int, token string) (*models.CountryDto, error)
}
