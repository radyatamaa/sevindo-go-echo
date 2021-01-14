package currency

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Delete(ctx context.Context, id int, token string) (*models.ResponseDelete, error)
	Update(ctx context.Context, ar *models.NewCommandCurrency, token string) error
	List(ctx context.Context, page, limit, offset int, search string) (*models.CurrencyWithPagination, error)
	Create(ctx context.Context, ar *models.NewCommandCurrency, token string) (*models.NewCommandCurrency, error)
	GetById(ctx context.Context, id int, token string) (*models.CurrencyDto, error)
}
