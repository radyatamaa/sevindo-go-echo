package usecase

import (
	"context"
	"github.com/master/country"
	"time"

	"github.com/models"
)

type countryUsecase struct {
	countryRepo    country.Repository
	contextTimeout time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewcountryUsecase(countryRepo country.Repository, timeout time.Duration) country.Usecase {
	return &countryUsecase{
		countryRepo:    countryRepo,
		contextTimeout: timeout,
	}
}

func (m countryUsecase) Create(c context.Context, ar *models.NewCommandCountry, token string) (*models.NewCommandCountry, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	insert := models.Country{
		Id:           0,
		CreatedBy:    "admin",
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		CountryName:  ar.CountryName,
	}

	err := m.countryRepo.Insert(ctx, &insert)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

func (m countryUsecase) GetById(c context.Context, id string, token string) (*models.CountryDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	country, err := m.countryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, models.ErrNotFound
	}

	result := &models.CountryDto{
		Id:          country.Id,
		CountryName: country.CountryName,
	}

	return result, nil
}
