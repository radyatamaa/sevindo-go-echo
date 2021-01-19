package resort_room_photo

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	GetByResortRoomID(ctx context.Context, resortRoomId string) ([]*models.ResortRoomPhoto, error)
}
