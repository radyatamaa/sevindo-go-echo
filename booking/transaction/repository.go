package transaction

import (
	"context"
	"github.com/models"
)

type Repository interface {
	Insert(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error)
}