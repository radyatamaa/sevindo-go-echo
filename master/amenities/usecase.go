package amenities

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Delete(ctx context.Context, id int, token string) (*models.ResponseDelete, error)
	Update(ctx context.Context, ar *models.NewCommandAmenities, token string) error
	List(ctx context.Context, page, limit, offset int, search string) (*models.AmenitiesWithPagination, error)
	Create(ctx context.Context, ar *models.NewCommandAmenities, token string) (*models.NewCommandAmenities, error)
	GetById(ctx context.Context, id int, token string) (*models.AmenitiesDto, error)
}
