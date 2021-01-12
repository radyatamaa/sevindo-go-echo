package country

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Create(ctx context.Context, ar *models.NewCommandCountry, token string) (*models.NewCommandCountry, error)
	GetById(ctx context.Context, id string, token string) (*models.CountryDto, error)
}
