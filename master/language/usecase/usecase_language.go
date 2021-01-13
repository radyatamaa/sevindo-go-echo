package usecase

import (
	"context"
	"time"

	"github.com/master/language"

	"github.com/models"
)

type languageUsecase struct {
	languageRepo   language.Repository
	contextTimeout time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewlanguageUsecase(languageRepo language.Repository, timeout time.Duration) language.Usecase {
	return &languageUsecase{
		languageRepo:   languageRepo,
		contextTimeout: timeout,
	}
}

func (m languageUsecase) Create(c context.Context, ar *models.NewCommandLanguage, token string) (*models.NewCommandLanguage, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	insert := models.Language{
		Id:           0,
		CreatedBy:    "admin",
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		LanguageName: ar.LanguageName,
	}

	err := m.languageRepo.Insert(ctx, &insert)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

func (m languageUsecase) GetById(c context.Context, id string, token string) (*models.LanguageDto, error) {
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
