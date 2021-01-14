package resort

import (
	"context"
	"github.com/models"
)

type Usecase interface {
	GetAll(ctx context.Context, page, limit, offset int, capacity int,startDate string, endDate string) (*models.ResortJoinDtoWithPagination, error)
}
