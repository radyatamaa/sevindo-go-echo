package usecase

import (
	"context"
	"time"

	"github.com/master/article_category"
	"github.com/models"
)

type ArticleCategoryUsecase struct {
	articlecategoryRepo article_category.Repository
	contextTimeout      time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewArticleCategoryUsecase(articlecategoryRepo article_category.Repository, timeout time.Duration) article_category.Usecase {
	return &ArticleCategoryUsecase{
		articlecategoryRepo: articlecategoryRepo,
		contextTimeout:      timeout,
	}
}

func (m ArticleCategoryUsecase) Create(c context.Context, ar *models.NewCommandArticleCategory, token string) (*models.NewCommandArticleCategory, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	insert := models.ArticleCategory{
		Id:                  0,
		CreatedBy:           "admin",
		CreatedDate:         time.Now(),
		ModifiedBy:          nil,
		ModifiedDate:        nil,
		DeletedBy:           nil,
		DeletedDate:         nil,
		IsDeleted:           0,
		IsActive:            1,
		ArticleCategoryName: ar.ArticleCategoryName,
	}

	err := m.articlecategoryRepo.Insert(ctx, &insert)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

func (m ArticleCategoryUsecase) GetById(c context.Context, id string, token string) (*models.ArticleCategoryDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	articlecategory, err := m.articlecategoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, models.ErrNotFound
	}

	result := &models.ArticleCategoryDto{
		Id:                  articlecategory.Id,
		ArticleCategoryName: articlecategory.ArticleCategoryName,
	}

	return result, nil
}
