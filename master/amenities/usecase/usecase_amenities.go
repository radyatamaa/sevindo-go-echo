package usecase

import (
	"context"

	"math"
	"strconv"
	"time"

	"github.com/auth/user_admin"
	"github.com/master/amenities"

	"github.com/models"
)

type amenitiesUsecase struct {
	userAdminUsecase user_admin.Usecase
	amenitiesRepo     amenities.Repository
	contextTimeout   time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewAmenitiesUsecase(userAdminUsecase user_admin.Usecase, amenitiesRepo amenities.Repository, timeout time.Duration) amenities.Usecase {
	return &amenitiesUsecase{
		userAdminUsecase: userAdminUsecase,
		amenitiesRepo:     amenitiesRepo,
		contextTimeout:   timeout,
	}
}

func (m amenitiesUsecase) Delete(c context.Context, id int, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser, err := m.userAdminUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return nil, models.ErrUnAuthorize
	}

	err = m.amenitiesRepo.Delete(ctx, id, currentUser.Email)

	result := &models.ResponseDelete{
		Id:      strconv.Itoa(id),
		Message: "Success Delete",
	}

	return result, nil
}

func (m amenitiesUsecase) Update(c context.Context, ar *models.NewCommandAmenities, token string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser, err := m.userAdminUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return models.ErrUnAuthorize
	}

	getCurrency, err := m.amenitiesRepo.GetByID(ctx, ar.Id)
	if err != nil {
		return err
	}
	var modifyBy string = currentUser.Email
	now := time.Now()
	getCurrency.Name = ar.Name
	getCurrency.ModifiedBy = &modifyBy
	getCurrency.ModifiedDate = &now
	err = m.amenitiesRepo.Update(ctx, getCurrency)
	if err != nil {
		return err
	}
	return nil
}

func (m amenitiesUsecase) Create(c context.Context, ar *models.NewCommandAmenities, token string) (*models.NewCommandAmenities, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser, err := m.userAdminUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return nil, models.ErrUnAuthorize
	}

	insert := models.Amenities{
		Id:           0,
		CreatedBy:    currentUser.Email,
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		Name: ar.Name,
	}

	err = m.amenitiesRepo.Insert(ctx, &insert)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

func (m amenitiesUsecase) GetById(c context.Context, id int, token string) (*models.AmenitiesDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currency, err := m.amenitiesRepo.GetByID(ctx, id)
	if err != nil {
		return nil, models.ErrNotFound
	}

	result := &models.AmenitiesDto{
		Id:           currency.Id,
		Name: currency.Name,
	}

	return result, nil
}

func (m amenitiesUsecase) List(ctx context.Context, page, limit, offset int, search string) (*models.AmenitiesWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.amenitiesRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	users := make([]*models.AmenitiesDto, len(list))
	for i, item := range list {
		users[i] = &models.AmenitiesDto{
			Id:          item.Id,
			Name: item.Name,
		}
	}
	totalRecords, _ := m.amenitiesRepo.Count(ctx)
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

	response := &models.AmenitiesWithPagination{
		Data: users,
		Meta: meta,
	}

	return response, nil
}
