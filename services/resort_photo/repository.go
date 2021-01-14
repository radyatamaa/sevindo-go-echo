package resort_photo

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetByResortID(ctx context.Context, resortId string) ([]*models.ResortPhoto, error)
}
