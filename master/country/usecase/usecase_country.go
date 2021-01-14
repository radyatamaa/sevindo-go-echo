package usecase

import (
	"context"
	"github.com/auth/user_admin"
	"github.com/master/country"
	"math"
	"strconv"
	"time"

	"github.com/models"
)

type countryUsecase struct {
	userAdminUsecase user_admin.Usecase
	countryRepo    country.Repository
	contextTimeout time.Duration
}



// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewcountryUsecase(	userAdminUsecase user_admin.Usecase,countryRepo country.Repository, timeout time.Duration) country.Usecase {
	return &countryUsecase{
		userAdminUsecase:userAdminUsecase,
		countryRepo:    countryRepo,
		contextTimeout: timeout,
	}
}
func (m countryUsecase) Delete(c context.Context, id int, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return nil,models.ErrUnAuthorize
	}

	err = m.countryRepo.Delete(ctx,id,currentUser.Email)

	result := &models.ResponseDelete{
		Id:      strconv.Itoa(id),
		Message: "Success Delete",
	}

	return result,nil
}

func (m countryUsecase) Update(c context.Context, ar *models.NewCommandCountry, token string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return models.ErrUnAuthorize
	}

	getCountry ,err := m.countryRepo.GetByID(ctx,ar.Id)
	if err != nil{
		return err
	}
	var modifyBy string = currentUser.Email
	now := time.Now()
	getCountry.CountryName = ar.CountryName
	getCountry.ModifiedBy = &modifyBy
	getCountry.ModifiedDate = &now
	err = m.countryRepo.Update(ctx,getCountry)
	if err != nil{
		return err
	}
	return nil
}

func (m countryUsecase) List(ctx context.Context, page, limit, offset int, search string) (*models.CountryWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.countryRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	users := make([]*models.CountryDto, len(list))
	for i, item := range list {
		users[i] = &models.CountryDto{
			Id:          item.Id,
			CountryName: item.CountryName,
		}
	}
	totalRecords, _ := m.countryRepo.Count(ctx)
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

	response := &models.CountryWithPagination{
		Data: users,
		Meta: meta,
	}

	return response, nil
}

func (m countryUsecase) Create(c context.Context, ar *models.NewCommandCountry, token string) (*models.NewCommandCountry, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return nil,models.ErrUnAuthorize
	}

	insert := models.Country{
		Id:           0,
		CreatedBy:    currentUser.Email,
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		CountryName:  ar.CountryName,
	}

	err = m.countryRepo.Insert(ctx, &insert)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

func (m countryUsecase) GetById(c context.Context, id int, token string) (*models.CountryDto, error) {
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
