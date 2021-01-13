package usecase

import (
	"context"
	"github.com/master/province"
	"time"

	"github.com/models"
)

type provinceUsecase struct {
	provinceRepo    province.Repository
	contextTimeout time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewprovinceUsecase(provinceRepo province.Repository, timeout time.Duration) province.Usecase {
	return &provinceUsecase{
		provinceRepo:    provinceRepo,
		contextTimeout: timeout,
	}
}

func (m provinceUsecase) Create(c context.Context, ar *models.NewCommandProvince, token string) (*models.NewCommandProvince, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	insert := models.Province{
		Id:           0,
		CreatedBy:    "admin",
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		ProvinceName:  ar.ProvinceName,
		CountryId:  ar.CountryId,
	}

	err := m.provinceRepo.Insert(ctx, &insert)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

func (m provinceUsecase) GetById(c context.Context, id string, token string) (*models.ProvinceDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	province, err := m.provinceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, models.ErrNotFound
	}

	result := &models.ProvinceDto{
		Id:          province.Id,
		ProvinceName: province.ProvinceName,
		CountryId: province.CountryId,
	}

	return result, nil
}
