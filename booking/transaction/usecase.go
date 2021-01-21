package transaction

import (
	"context"
	"github.com/models"
)

type Usecase interface {
	Insert(ctx context.Context, transaction *models.Transaction,token string) (*models.Transaction, error)
}
