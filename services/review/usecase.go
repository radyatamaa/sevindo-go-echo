package review

import (
	"context"
	"github.com/models"
)

type Usecase interface {
	GetAll(ctx context.Context, page, limit, offset int,resortId string) (*models.ReviewDtoWithPagination, error)
}
