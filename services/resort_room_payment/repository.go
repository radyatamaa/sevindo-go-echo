package resort_room_payment

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetByResortRoomID(ctx context.Context, resortRoomId string) ([]*models.ResortRoomPayment, error)
}
