package usecase

import (
	"context"
	"github.com/auth/user"
	"github.com/models"
	"github.com/services/review"
	"math"
	"time"
)

type reviewUsecase struct {
	userUsecase user.Usecase
	reviewRepo    review.Repository
	contextTimeout time.Duration
}



// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewreviewUsecase(	userUsecase user.Usecase,reviewRepo    review.Repository, timeout time.Duration) review.Usecase {
	return &reviewUsecase{
		userUsecase:userUsecase,
		reviewRepo:    reviewRepo,
		contextTimeout: timeout,
	}
}

func (m reviewUsecase) GetAll(ctx context.Context, page, limit, offset int,resortId string) (*models.ReviewDtoWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.reviewRepo.GetByResortIdJoinWithPayment(ctx,resortId, limit, offset)
	if err != nil {
		return nil, err
	}

	users := make([]*models.ReviewDto, len(list))
	for i, item := range list {
		users[i] = &models.ReviewDto{
			Id:            item.Id,
			CreatedDate:   item.CreatedDate,
			Values:        item.Values,
			Desc:          item.Desc,
			UserId:        item.UserId,
			TransactionId: item.TransactionId,
			Name:          item.Name,
		}
	}
	totalRecords, _ := m.reviewRepo.Count(ctx,resortId)
	totalPage := int(math.Ceil(float64(totalRecords) / float64(limit)))
	prev := page
	next := page
	if page != 1 {
		prev = page - 1
	}

	if page != totalPage {
		next = page + 1
	}
	meta := &models.MetaPagination{
		Page:          page,
		Total:         totalPage,
		TotalRecords:  totalRecords,
		Prev:          prev,
		Next:          next,
		RecordPerPage: len(list),
	}

	response := &models.ReviewDtoWithPagination{
		Data: users,
		Meta: meta,
	}

	return response, nil
}