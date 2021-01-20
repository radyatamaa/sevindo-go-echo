package booking

import (
	"context"
	"github.com/models"
)

type Usecase interface {
	Insert(ctx context.Context, booking *models.NewBookingCommand,token string) ([]*models.NewBookingCommand, error, error)
}
