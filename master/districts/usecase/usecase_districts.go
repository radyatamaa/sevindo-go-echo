package usecase

import (
	"context"
	"github.com/auth/user_admin"
	"github.com/master/districts"
	"math"
	"strconv"
	"time"

	"github.com/models"
)

type districtsUsecase struct {
	userAdminUsecase user_admin.Usecase
	districtsRepo    districts.Repository
	contextTimeout time.Duration
}



// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewdistrictsUsecase(	userAdminUsecase user_admin.Usecase,districtsRepo districts.Repository, timeout time.Duration) districts.Usecase {
	return &districtsUsecase{
		userAdminUsecase:userAdminUsecase,
		districtsRepo:    districtsRepo,
		contextTimeout: timeout,
	}
}
func (m districtsUsecase) Delete(c context.Context, id int, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return nil,models.ErrUnAuthorize
	}

	err = m.districtsRepo.Delete(ctx,id,currentUser.Email)

	result := &models.ResponseDelete{
		Id:      strconv.Itoa(id),
		Message: "Success Delete",
	}

	return result,nil
}

func (m districtsUsecase) Update(c context.Context, ar *models.NewCommandDistricts, token string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return models.ErrUnAuthorize
	}

	getDistricts ,err := m.districtsRepo.GetByID(ctx,ar.Id)
	if err != nil{
		return err
	}
	var modifyBy string = currentUser.Email
	now := time.Now()
	getDistricts.DistrictsName = ar.DistrictsName
	getDistricts.CityId = ar.CityId
	getDistricts.ModifiedBy = &modifyBy
	getDistricts.ModifiedDate = &now
	err = m.districtsRepo.Update(ctx,getDistricts)
	if err != nil{
		return err
	}
	return nil
}

func (m districtsUsecase) List(ctx context.Context, page, limit, offset int, search string) (*models.DistrictsWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.districtsRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	users := make([]*models.DistrictsDto, len(list))
	for i, item := range list {
		users[i] = &models.DistrictsDto{
			Id:          item.Id,
			DistrictsName: item.DistrictsName,
			CityId: item.CityId,
		}
	}
	totalRecords, _ := m.districtsRepo.Count(ctx)
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

	response := &models.DistrictsWithPagination{
		Data: users,
		Meta: meta,
	}

	return response, nil
}

func (m districtsUsecase) Create(c context.Context, ar *models.NewCommandDistricts, token string) (*models.NewCommandDistricts, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return nil,models.ErrUnAuthorize
	}

	insert := models.Districts{
		Id:           0,
		CreatedBy:    currentUser.Email,
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		DistrictsName:  ar.DistrictsName,
		CityId:  ar.CityId,
	}

	err = m.districtsRepo.Insert(ctx, &insert)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

func (m districtsUsecase) GetById(c context.Context, id int, token string) (*models.DistrictsDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	districts, err := m.districtsRepo.GetByID(ctx, id)
	if err != nil {
		return nil, models.ErrNotFound
	}

	result := &models.DistrictsDto{
		Id:          districts.Id,
		DistrictsName: districts.DistrictsName,
	}

	return result, nil
}
