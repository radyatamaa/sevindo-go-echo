package districts

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Create(ctx context.Context, ar *models.NewCommandDistricts, token string) (*models.NewCommandDistricts, error)
	GetById(ctx context.Context, id int, token string) (*models.DistrictsDto, error)
	Delete(ctx context.Context, id int, token string) (*models.ResponseDelete, error)
	Update(ctx context.Context, ar *models.NewCommandDistricts,  token string) error
	List(ctx context.Context, page, limit, offset int, search string) (*models.DistrictsWithPagination, error)
}
