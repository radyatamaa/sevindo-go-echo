package usecase

import (
	"context"
	"math"
	"strconv"
	"time"

	"github.com/auth/user_admin"
	"github.com/master/currency"

	"github.com/models"
)

type currencyUsecase struct {
	userAdminUsecase user_admin.Usecase
	currencyRepo     currency.Repository
	contextTimeout   time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewcurrencyUsecase(userAdminUsecase user_admin.Usecase, currencyRepo currency.Repository, timeout time.Duration) currency.Usecase {
	return &currencyUsecase{
		userAdminUsecase: userAdminUsecase,
		currencyRepo:     currencyRepo,
		contextTimeout:   timeout,
	}
}

func (m currencyUsecase) Delete(c context.Context, id int, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser, err := m.userAdminUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return nil, models.ErrUnAuthorize
	}

	err = m.currencyRepo.Delete(ctx, id, currentUser.Email)

	result := &models.ResponseDelete{
		Id:      strconv.Itoa(id),
		Message: "Success Delete",
	}

	return result, nil
}

func (m currencyUsecase) Update(c context.Context, ar *models.NewCommandCurrency, token string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser, err := m.userAdminUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return models.ErrUnAuthorize
	}

	getCurrency, err := m.currencyRepo.GetByID(ctx, ar.Id)
	if err != nil {
		return err
	}
	var modifyBy string = currentUser.Email
	now := time.Now()
	getCurrency.CurrencyName = ar.CurrencyName
	getCurrency.ModifiedBy = &modifyBy
	getCurrency.ModifiedDate = &now
	err = m.currencyRepo.Update(ctx, getCurrency)
	if err != nil {
		return err
	}
	return nil
}

func (m currencyUsecase) Create(c context.Context, ar *models.NewCommandCurrency, token string) (*models.NewCommandCurrency, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser, err := m.userAdminUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return nil, models.ErrUnAuthorize
	}

	insert := models.Currency{
		Id:           0,
		CreatedBy:    currentUser.Email,
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		CurrencyName: ar.CurrencyName,
	}

	err = m.currencyRepo.Insert(ctx, &insert)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

func (m currencyUsecase) GetById(c context.Context, id int, token string) (*models.CurrencyDto, error) {
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

func (m currencyUsecase) List(ctx context.Context, page, limit, offset int, search string) (*models.CurrencyWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.currencyRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	users := make([]*models.CurrencyDto, len(list))
	for i, item := range list {
		users[i] = &models.CurrencyDto{
			Id:          item.Id,
			CurrencyName: item.CurrencyName,
		}
	}
	totalRecords, _ := m.currencyRepo.Count(ctx)
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

	response := &models.CurrencyWithPagination{
		Data: users,
		Meta: meta,
	}

	return response, nil
}
