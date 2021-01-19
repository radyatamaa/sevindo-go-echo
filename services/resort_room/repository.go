package resort_room

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetByResortIdJoinWithPayment(ctx context.Context, resortID string) ([]*models.ResortRoom, error)
}
