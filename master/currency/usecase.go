package currency

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Create(ctx context.Context, ar *models.NewCommandCurrency, token string) (*models.NewCommandCurrency, error)
	GetById(ctx context.Context, id string, token string) (*models.CurrencyDto, error)
}
