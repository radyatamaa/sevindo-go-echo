package usecase

import (
	"context"
	"github.com/auth/user_admin"
	"github.com/master/city"
	"math"
	"strconv"
	"time"

	"github.com/models"
)

type cityUsecase struct {
	userAdminUsecase user_admin.Usecase
	cityRepo   city.Repository
	contextTimeout time.Duration
}



// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewcityUsecase(	userAdminUsecase user_admin.Usecase,cityRepo city.Repository, timeout time.Duration) city.Usecase {
	return &cityUsecase{
		userAdminUsecase:userAdminUsecase,
		cityRepo:    cityRepo,
		contextTimeout: timeout,
	}
}
func (m cityUsecase) Delete(c context.Context, id int, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return nil,models.ErrUnAuthorize
	}

	err = m.cityRepo.Delete(ctx,id,currentUser.Email)

	result := &models.ResponseDelete{
		Id:      strconv.Itoa(id),
		Message: "Success Delete",
	}

	return result,nil
}

func (m cityUsecase) Update(c context.Context, ar *models.NewCommandCity, token string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return models.ErrUnAuthorize
	}

	getCity ,err := m.cityRepo.GetByID(ctx,ar.Id)
	if err != nil{
		return err
	}
	var modifyBy string = currentUser.Email
	now := time.Now()
	getCity.CityName = ar.CityName
	getCity.ProvinceId = ar.ProvinceId
	getCity.ModifiedBy = &modifyBy
	getCity.ModifiedDate = &now
	err = m.cityRepo.Update(ctx,getCity)
	if err != nil{
		return err
	}
	return nil
}

func (m cityUsecase) List(ctx context.Context, page, limit, offset int, search string) (*models.CityWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.cityRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	users := make([]*models.CityDto, len(list))
	for i, item := range list {
		users[i] = &models.CityDto{
			Id:          item.Id,
			CityName: item.CityName,
			ProvinceId: item.ProvinceId,
		}
	}
	totalRecords, _ := m.cityRepo.Count(ctx)
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

	response := &models.CityWithPagination{
		Data: users,
		Meta: meta,
	}

	return response, nil
}

func (m cityUsecase) Create(c context.Context, ar *models.NewCommandCity, token string) (*models.NewCommandCity, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return nil,models.ErrUnAuthorize
	}

	insert := models.City{
		Id:           0,
		CreatedBy:    currentUser.Email,
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		CityName:  ar.CityName,
		ProvinceId:  ar.ProvinceId,
	}

	err = m.cityRepo.Insert(ctx, &insert)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

func (m cityUsecase) GetById(c context.Context, id int, token string) (*models.CityDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	city, err := m.cityRepo.GetByID(ctx, id)
	if err != nil {
		return nil, models.ErrNotFound
	}

	result := &models.CityDto{
		Id:          city.Id,
		CityName: city.CityName,
	}

	return result, nil
}
