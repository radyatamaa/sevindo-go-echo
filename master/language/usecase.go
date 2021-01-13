package language

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Create(ctx context.Context, ar *models.NewCommandLanguage, token string) (*models.NewCommandLanguage, error)
	GetById(ctx context.Context, id string, token string) (*models.LanguageDto, error)
}
