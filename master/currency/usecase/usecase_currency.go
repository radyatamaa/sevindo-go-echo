package usecase

import (
	"context"
	"time"

	"github.com/master/currency"

	"github.com/models"
)

type currencyUsecase struct {
	currencyRepo   currency.Repository
	contextTimeout time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewcurrencyUsecase(currencyRepo currency.Repository, timeout time.Duration) currency.Usecase {
	return &currencyUsecase{
		currencyRepo:   currencyRepo,
		contextTimeout: timeout,
	}
}

func (m currencyUsecase) Create(c context.Context, ar *models.NewCommandCurrency, token string) (*models.NewCommandCurrency, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	insert := models.Currency{
		Id:           0,
		CreatedBy:    "admin",
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		CurrencyName: ar.CurrencyName,
	}

	err := m.currencyRepo.Insert(ctx, &insert)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

func (m currencyUsecase) GetById(c context.Context, id string, token string) (*models.CurrencyDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currency, err := m.currencyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, models.ErrNotFound
	}

	result := &models.CurrencyDto{
		Id:           currency.Id,
		CurrencyName: currency.CurrencyName,
	}

	return result, nil
}
