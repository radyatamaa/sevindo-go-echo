package bank

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Delete(ctx context.Context, id int, token string) (*models.ResponseDelete, error)
	Update(ctx context.Context, ar *models.NewCommandBank, token string) error
	List(ctx context.Context, page, limit, offset int, search string) (*models.BankWithPagination, error)
	Create(ctx context.Context, ar *models.NewCommandBank, token string) (*models.NewCommandBank, error)
	GetById(ctx context.Context, id int, token string) (*models.BankDto, error)
}
