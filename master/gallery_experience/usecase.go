package gallery_experience

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Create(ctx context.Context, ar *models.NewCommandGalleryExperience, token string) (*models.NewCommandGalleryExperience, error)
	GetById(ctx context.Context, id string, token string) (*models.GalleryExperienceDto, error)
	Delete(ctx context.Context, id string, token string) (*models.ResponseDelete, error)
	Update(ctx context.Context, ar *models.NewCommandGalleryExperience,  token string) error
	List(ctx context.Context, page, limit, offset int, search string) (*models.GalleryExperienceWithPagination, error)
}
