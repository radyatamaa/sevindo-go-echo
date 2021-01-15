package usecase

import (
	"context"
	"github.com/auth/user_admin"
	"github.com/master/province"
	"math"
	"strconv"
	"time"

	"github.com/models"
)

type provinceUsecase struct {
	userAdminUsecase user_admin.Usecase
	provinceRepo    province.Repository
	contextTimeout time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewprovinceUsecase(userAdminUsecase user_admin.Usecase,provinceRepo province.Repository, timeout time.Duration) province.Usecase {
	return &provinceUsecase{
		userAdminUsecase:userAdminUsecase,
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

func (m provinceUsecase) GetById(c context.Context, id int, token string) (*models.ProvinceDto, error) {
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
func (m provinceUsecase) Delete(c context.Context, id int, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return nil,models.ErrUnAuthorize
	}

	err = m.provinceRepo.Delete(ctx,id,currentUser.Email)

	result := &models.ResponseDelete{
		Id:      strconv.Itoa(id),
		Message: "Success Delete",
	}

	return result,nil
}

func (m provinceUsecase) Update(c context.Context, ar *models.NewCommandProvince, token string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return models.ErrUnAuthorize
	}

	getProvince ,err := m.provinceRepo.GetByID(ctx,ar.Id)
	if err != nil{
		return err
	}
	var modifyBy string = currentUser.Email
	now := time.Now()
	getProvince.ProvinceName = ar.ProvinceName
	getProvince.CountryId = ar.CountryId
	getProvince.ModifiedBy = &modifyBy
	getProvince.ModifiedDate = &now
	err = m.provinceRepo.Update(ctx,getProvince)
	if err != nil{
		return err
	}
	return nil
}

func (m provinceUsecase) List(ctx context.Context, page, limit, offset int, search string) (*models.ProvinceWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.provinceRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	users := make([]*models.ProvinceDto, len(list))
	for i, item := range list {
		users[i] = &models.ProvinceDto{
			Id:          item.Id,
			ProvinceName: item.ProvinceName,
			CountryId: item.CountryId,
		}
	}
	totalRecords, _ := m.provinceRepo.Count(ctx)
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

	response := &models.ProvinceWithPagination{
		Data: users,
		Meta: meta,
	}

	return response, nil
}

