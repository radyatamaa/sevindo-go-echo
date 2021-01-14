package usecase

import (
	"context"
	"github.com/auth/user_admin"
	"math"
	"strconv"
	"time"

	"github.com/master/language"

	"github.com/models"
)

type languageUsecase struct {
	userAdminUsecase user_admin.Usecase
	languageRepo   language.Repository
	contextTimeout time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewlanguageUsecase(userAdminUsecase user_admin.Usecase, languageRepo language.Repository, timeout time.Duration) language.Usecase {
	return &languageUsecase{
		userAdminUsecase: userAdminUsecase,
		languageRepo:   languageRepo,
		contextTimeout: timeout,
	}
}

func (m languageUsecase) Create(c context.Context, ar *models.NewCommandLanguage, token string) (*models.NewCommandLanguage, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser, err := m.userAdminUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return nil, models.ErrUnAuthorize
	}

	insert := models.Language{
		Id:           0,
		CreatedBy:    currentUser.Email,
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		LanguageName: ar.LanguageName,
	}

	err = m.languageRepo.Insert(ctx, &insert)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

func (m languageUsecase) GetById(c context.Context, id int, token string) (*models.LanguageDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()



	language, err := m.languageRepo.GetByID(ctx, id)
	if err != nil {
		return nil, models.ErrNotFound
	}

	result := &models.LanguageDto{
		Id:            language.Id,
		LanguageName: language.LanguageName,
	}

	return result, nil
}

func (m languageUsecase) Delete(c context.Context, id int, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser, err := m.userAdminUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return nil, models.ErrUnAuthorize
	}

	err = m.languageRepo.Delete(ctx, id, currentUser.Email)

	result := &models.ResponseDelete{
		Id:      strconv.Itoa(id),
		Message: "Success Delete",
	}

	return result, nil
}

func (m languageUsecase) Update(c context.Context, ar *models.NewCommandLanguage, token string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser, err := m.userAdminUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return models.ErrUnAuthorize
	}

	getLanguage, err := m.languageRepo.GetByID(ctx, ar.Id)
	if err != nil {
		return err
	}
	var modifyBy string = currentUser.Email
	now := time.Now()
	getLanguage.LanguageName = ar.LanguageName
	getLanguage.ModifiedBy = &modifyBy
	getLanguage.ModifiedDate = &now
	err = m.languageRepo.Update(ctx, getLanguage)
	if err != nil {
		return err
	}
	return nil
}

func (m languageUsecase) List(ctx context.Context, page, limit, offset int, search string) (*models.LanguageWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.languageRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	users := make([]*models.LanguageDto, len(list))
	for i, item := range list {
		users[i] = &models.LanguageDto{
			Id:          item.Id,
			LanguageName: item.LanguageName,
		}
	}
	totalRecords, _ := m.languageRepo.Count(ctx)
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

	response := &models.LanguageWithPagination{
		Data: users,
		Meta: meta,
	}

	return response, nil
}