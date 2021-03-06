package accessibility_resort

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetByResortId(ctx context.Context, resortId string) ([]*models.AccessibilityResortJoin, error)
}
