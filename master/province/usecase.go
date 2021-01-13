package province

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Create(ctx context.Context, ar *models.NewCommandProvince, token string) (*models.NewCommandProvince, error)
	GetById(ctx context.Context, id string, token string) (*models.ProvinceDto, error)
}
