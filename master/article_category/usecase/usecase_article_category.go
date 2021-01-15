package usecase

import (
	"context"
	"math"
	"strconv"
	"time"
	"github.com/auth/user_admin"
	"github.com/master/article_category"
	"github.com/models"
)

type ArticleCategoryUsecase struct {
	userAdminUsecase user_admin.Usecase
	articlecategoryRepo article_category.Repository
	contextTimeout      time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewArticleCategoryUsecase(userAdminUsecase user_admin.Usecase, articlecategoryRepo article_category.Repository, timeout time.Duration) article_category.Usecase {
	return &ArticleCategoryUsecase{
		userAdminUsecase: userAdminUsecase,
		articlecategoryRepo: articlecategoryRepo,
		contextTimeout:      timeout,
	}
}

func (m ArticleCategoryUsecase) Delete(c context.Context, id int, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser, err := m.userAdminUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return nil, models.ErrUnAuthorize
	}

	err = m.articlecategoryRepo.Delete(ctx, id, currentUser.Email)

	result := &models.ResponseDelete{
		Id:      strconv.Itoa(id),
		Message: "Success Delete",
	}

	return result, nil
}

func (m ArticleCategoryUsecase) Update(c context.Context, ar *models.NewCommandArticleCategory, token string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser, err := m.userAdminUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return models.ErrUnAuthorize
	}

	getArticleCategory, err := m.articlecategoryRepo.GetByID(ctx, ar.Id)
	if err != nil {
		return err
	}
	var modifyBy string = currentUser.Email
	now := time.Now()
	getArticleCategory.ArticleCategoryName = ar.ArticleCategoryName
	getArticleCategory.ModifiedBy = &modifyBy
	getArticleCategory.ModifiedDate = &now
	err = m.articlecategoryRepo.Update(ctx, getArticleCategory)
	if err != nil {
		return err
	}
	return nil
}

func (m ArticleCategoryUsecase) Create(c context.Context, ar *models.NewCommandArticleCategory, token string) (*models.NewCommandArticleCategory, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser, err := m.userAdminUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return nil, models.ErrUnAuthorize
	}
	insert := models.ArticleCategory{
		Id:                  0,
		CreatedBy:           currentUser.Email,
		CreatedDate:         time.Now(),
		ModifiedBy:          nil,
		ModifiedDate:        nil,
		DeletedBy:           nil,
		DeletedDate:         nil,
		IsDeleted:           0,
		IsActive:            1,
		ArticleCategoryName: ar.ArticleCategoryName,
	}

	err = m.articlecategoryRepo.Insert(ctx, &insert)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

func (m ArticleCategoryUsecase) GetById(c context.Context, id int, token string) (*models.ArticleCategoryDto, error) {
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
func (m ArticleCategoryUsecase) List(ctx context.Context, page, limit, offset int, search string) (*models.ArticleCategoryWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.articlecategoryRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	users := make([]*models.ArticleCategoryDto, len(list))
	for i, item := range list {
		users[i] = &models.ArticleCategoryDto{
			Id:          item.Id,
			ArticleCategoryName: item.ArticleCategoryName,
		}
	}
	totalRecords, _ := m.articlecategoryRepo.Count(ctx)
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

	response := &models.ArticleCategoryWithPagination{
		Data: users,
		Meta: meta,
	}

	return response, nil
}
