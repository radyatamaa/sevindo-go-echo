package booking

import (
	"context"
	"github.com/models"
)

type Repository interface {
	Insert(ctx context.Context, booking *models.Booking) (*models.Booking, error)
	CheckBookingCode(ctx context.Context, bookingCode string) bool
}
