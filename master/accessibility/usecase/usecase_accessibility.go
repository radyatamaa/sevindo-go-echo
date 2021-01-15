package usecase

import (
	"context"
	"github.com/auth/user_admin"
	"github.com/master/accessibility"
	"github.com/models"
	"math"
	"strconv"
	"time"
)

type accessibilityUsecase struct {
	userAdminUsecase user_admin.Usecase
	accessibilityRepo    accessibility.Repository
	contextTimeout time.Duration
}



// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewaccessibilityUsecase(	userAdminUsecase user_admin.Usecase,accessibilityRepo accessibility.Repository, timeout time.Duration) accessibility.Usecase {
	return &accessibilityUsecase{
		userAdminUsecase:userAdminUsecase,
		accessibilityRepo:    accessibilityRepo,
		contextTimeout: timeout,
	}
}
func (m accessibilityUsecase) Delete(c context.Context, id int, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return nil,models.ErrUnAuthorize
	}

	err = m.accessibilityRepo.Delete(ctx,id,currentUser.Email)

	result := &models.ResponseDelete{
		Id:      strconv.Itoa(id),
		Message: "Success Delete",
	}

	return result,nil
}

func (m accessibilityUsecase) Update(c context.Context, ar *models.NewCommandAccessibility, token string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return models.ErrUnAuthorize
	}

	getAccessibility ,err := m.accessibilityRepo.GetByID(ctx,ar.Id)
	if err != nil{
		return err
	}
	var modifyBy string = currentUser.Email
	now := time.Now()
	getAccessibility.Name = ar.Name
	getAccessibility.ModifiedBy = &modifyBy
	getAccessibility.ModifiedDate = &now
	err = m.accessibilityRepo.Update(ctx,getAccessibility)
	if err != nil{
		return err
	}
	return nil
}

func (m accessibilityUsecase) List(ctx context.Context, page, limit, offset int, search string) (*models.AccessibilityWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.accessibilityRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	users := make([]*models.AccessibilityDto, len(list))
	for i, item := range list {
		users[i] = &models.AccessibilityDto{
			Id:          item.Id,
			Name: item.Name,
		}
	}
	totalRecords, _ := m.accessibilityRepo.Count(ctx)
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

	response := &models.AccessibilityWithPagination{
		Data: users,
		Meta: meta,
	}

	return response, nil
}

func (m accessibilityUsecase) Create(c context.Context, ar *models.NewCommandAccessibility, token string) (*models.NewCommandAccessibility, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return nil,models.ErrUnAuthorize
	}

	insert := models.Accessibility{
		Id:           0,
		CreatedBy:    currentUser.Email,
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		Name:  ar.Name,
	}

	err = m.accessibilityRepo.Insert(ctx, &insert)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

func (m accessibilityUsecase) GetById(c context.Context, id int, token string) (*models.AccessibilityDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	accessibility, err := m.accessibilityRepo.GetByID(ctx, id)
	if err != nil {
		return nil, models.ErrNotFound
	}

	result := &models.AccessibilityDto{
		Id:          accessibility.Id,
		Name: accessibility.Name,
	}

	return result, nil
}
