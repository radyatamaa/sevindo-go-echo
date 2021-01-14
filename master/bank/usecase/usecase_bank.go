package usecase

import (
	"context"
	"math"
	"strconv"
	"time"

	"github.com/auth/user_admin"
	"github.com/master/bank"

	"github.com/models"
)

type bankUsecase struct {
	userAdminUsecase user_admin.Usecase
	bankRepo     bank.Repository
	contextTimeout   time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewbankUsecase(userAdminUsecase user_admin.Usecase, bankRepo bank.Repository, timeout time.Duration) bank.Usecase {
	return &bankUsecase{
		userAdminUsecase: userAdminUsecase,
		bankRepo:     bankRepo,
		contextTimeout:   timeout,
	}
}

func (m bankUsecase) Delete(c context.Context, id int, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser, err := m.userAdminUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return nil, models.ErrUnAuthorize
	}

	err = m.bankRepo.Delete(ctx, id, currentUser.Email)

	result := &models.ResponseDelete{
		Id:      strconv.Itoa(id),
		Message: "Success Delete",
	}

	return result, nil
}

func (m bankUsecase) Update(c context.Context, ar *models.NewCommandBank, token string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser, err := m.userAdminUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return models.ErrUnAuthorize
	}

	getBank, err := m.bankRepo.GetByID(ctx, ar.Id)
	if err != nil {
		return err
	}
	var modifyBy string = currentUser.Email
	now := time.Now()
	getBank.BankName = ar.BankName
	getBank.ModifiedBy = &modifyBy
	getBank.ModifiedDate = &now
	err = m.bankRepo.Update(ctx, getBank)
	if err != nil {
		return err
	}
	return nil
}

func (m bankUsecase) Create(c context.Context, ar *models.NewCommandBank, token string) (*models.NewCommandBank, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser, err := m.userAdminUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return nil, models.ErrUnAuthorize
	}

	insert := models.Bank{
		Id:           0,
		CreatedBy:    currentUser.Email,
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		BankName: ar.BankName,
	}

	err = m.bankRepo.Insert(ctx, &insert)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

func (m bankUsecase) GetById(c context.Context, id int, token string) (*models.BankDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	bank, err := m.bankRepo.GetByID(ctx, id)
	if err != nil {
		return nil, models.ErrNotFound
	}

	result := &models.BankDto{
		Id:           bank.Id,
		BankName: bank.BankName,
	}

	return result, nil
}

func (m bankUsecase) List(ctx context.Context, page, limit, offset int, search string) (*models.BankWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.bankRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	users := make([]*models.BankDto, len(list))
	for i, item := range list {
		users[i] = &models.BankDto{
			Id:          item.Id,
			BankName: item.BankName,
		}
	}
	totalRecords, _ := m.bankRepo.Count(ctx)
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

	response := &models.BankWithPagination{
		Data: users,
		Meta: meta,
	}

	return response, nil
}
