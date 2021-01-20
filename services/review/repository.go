package review

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetByResortIdJoinWithPayment(ctx context.Context, resortID string,limit int ,offset int) ([]*models.ReviewJoin, error)
	Count(ctx context.Context, resortID string) (int, error)
}
