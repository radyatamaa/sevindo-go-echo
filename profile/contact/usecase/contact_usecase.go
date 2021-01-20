package usecase

import (
	"context"
	"github.com/auth/user"
	guuid "github.com/google/uuid"
	"github.com/models"
	"github.com/profile/contact"
	"math"
	"time"
)

type contactUsecase struct {
	userUsecase user.Usecase
	contactRepo    contact.Repository
	contextTimeout time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewcontactUsecase(	userUsecase user.Usecase,contactRepo    contact.Repository, timeout time.Duration) contact.Usecase {
	return &contactUsecase{
		userUsecase:userUsecase,
		contactRepo:    contactRepo,
		contextTimeout: timeout,
	}
}


func (m contactUsecase) Delete(c context.Context, id string, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return nil,models.ErrUnAuthorize
	}

	err = m.contactRepo.Delete(ctx,id,currentUser.UserEmail)

	result := &models.ResponseDelete{
		Id:      id,
		Message: "Success Delete",
	}

	return result,nil
}

func (m contactUsecase) Update(c context.Context, ar *models.NewCommandContact, token string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return models.ErrUnAuthorize
	}

	getAccessibility ,err := m.contactRepo.GetByID(ctx,ar.Id)
	if err != nil{
		return err
	}
	var modifyBy string = currentUser.UserEmail
	now := time.Now()
	getAccessibility.FullName = ar.FullName
	getAccessibility.TypeAs = ar.TypeAs
	getAccessibility.PhoneNumber = ar.PhoneNumber
	getAccessibility.ModifiedBy = &modifyBy
	getAccessibility.ModifiedDate = &now
	err = m.contactRepo.Update(ctx,getAccessibility)
	if err != nil{
		return err
	}
	return nil
}

func (m contactUsecase) List(ctx context.Context, page, limit, offset int, token string) (*models.ContactWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return nil,models.ErrUnAuthorize
	}

	list, err := m.contactRepo.List(ctx, limit, offset,currentUser.Id)
	if err != nil {
		return nil, err
	}

	users := make([]*models.ContactDto, len(list))
	for i, item := range list {
		users[i] = &models.ContactDto{
			Id:          item.Id,
			FullName:    item.FullName,
			TypeAs:      item.TypeAs,
			PhoneNumber: item.PhoneNumber,
			UserId:      item.UserId,
		}
	}
	totalRecords, _ := m.contactRepo.Count(ctx,currentUser.Id)
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

	response := &models.ContactWithPagination{
		Data: users,
		Meta: meta,
	}

	return response, nil
}

func (m contactUsecase) Create(c context.Context, ar *models.NewCommandContact, token string) (*models.NewCommandContact, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return nil,models.ErrUnAuthorize
	}

	insert := models.Contact{
		Id:           guuid.New().String(),
		CreatedBy:    currentUser.UserEmail,
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		FullName:     ar.FullName,
		TypeAs:       ar.TypeAs,
		PhoneNumber:  ar.PhoneNumber,
		UserId:       currentUser.Id,
	}

	err = m.contactRepo.Insert(ctx, &insert)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

func (m contactUsecase) GetById(c context.Context, id string, token string) (*models.ContactDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	accessibility, err := m.contactRepo.GetByID(ctx, id)
	if err != nil {
		return nil, models.ErrNotFound
	}

	result := &models.ContactDto{
		Id:          accessibility.Id,
		FullName:    accessibility.FullName,
		TypeAs:      accessibility.TypeAs,
		PhoneNumber: accessibility.PhoneNumber,
		UserId:      accessibility.UserId,
	}

	return result, nil
}

