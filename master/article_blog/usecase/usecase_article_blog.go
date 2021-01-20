package usecase

import (
	"context"

	"math"
	"strconv"
	"time"
	"github.com/auth/user_admin"
	"github.com/master/article_blog"
	"github.com/models"
)

type ArticleBlogUsecase struct {
	userAdminUsecase user_admin.Usecase
	articleblogRepo article_blog.Repository
	contextTimeout      time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewArticleBlogUsecase(userAdminUsecase user_admin.Usecase, articleblogRepo article_blog.Repository, timeout time.Duration) article_blog.Usecase {
	return &ArticleBlogUsecase{
		userAdminUsecase: userAdminUsecase,
		articleblogRepo: articleblogRepo,
		contextTimeout:      timeout,
	}
}

func (m ArticleBlogUsecase) Delete(c context.Context, id int, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser, err := m.userAdminUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return nil, models.ErrUnAuthorize
	}

	err = m.articleblogRepo.Delete(ctx, id, currentUser.Email)

	result := &models.ResponseDelete{
		Id:      strconv.Itoa(id),
		Message: "Success Delete",
	}

	return result, nil
}

func (m ArticleBlogUsecase) Update(c context.Context, ar *models.NewCommandArticleBlog, token string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser, err := m.userAdminUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return models.ErrUnAuthorize
	}

	getArticleBlog, err := m.articleblogRepo.GetByID(ctx, ar.Id)
	if err != nil {
		return err
	}
	var modifyBy string = currentUser.Email
	now := time.Now()
	getArticleBlog.ArticleBlogName = ar.ArticleBlogName
	getArticleBlog.Title = ar.Title
	getArticleBlog.Description = ar.Description
	getArticleBlog.CategoryId = ar.CategoryId
	getArticleBlog.ArticlePicture = ar.ArticlePicture
	getArticleBlog.ModifiedBy = &modifyBy
	getArticleBlog.ModifiedDate = &now
	err = m.articleblogRepo.Update(ctx, getArticleBlog)
	if err != nil {
		return err
	}
	return nil
}

func (m ArticleBlogUsecase) Create(c context.Context, ar *models.NewCommandArticleBlog, token string) (*models.NewCommandArticleBlog, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser, err := m.userAdminUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return nil, models.ErrUnAuthorize
	}
	insert := models.ArticleBlog{
		Id:                  0,
		CreatedBy:           currentUser.Email,
		CreatedDate:         time.Now(),
		ModifiedBy:          nil,
		ModifiedDate:        nil,
		DeletedBy:           nil,
		DeletedDate:         nil,
		IsDeleted:           0,
		IsActive:            1,
		ArticleBlogName: ar.ArticleBlogName,
		Title: ar.Title,
		Description: ar.Description,
		CategoryId: ar.CategoryId,
		ArticlePicture: ar.ArticlePicture,
	}

	err = m.articleblogRepo.Insert(ctx, &insert)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

func (m ArticleBlogUsecase) GetById(c context.Context, id int, token string) (*models.ArticleBlogDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	articleblog, err := m.articleblogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, models.ErrNotFound
	}

	result := &models.ArticleBlogDto{
		Id:                  articleblog.Id,
		ArticleBlogName: articleblog.ArticleBlogName,
		Title: articleblog.Title,
		Description: articleblog.Description,
		CategoryId: articleblog.CategoryId,
	}

	return result, nil
}
func (m ArticleBlogUsecase) List(ctx context.Context, page, limit, offset int, search string) (*models.ArticleBlogWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.articleblogRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	users := make([]*models.ArticleBlogDto, len(list))
	for i, item := range list {
		users[i] = &models.ArticleBlogDto{
			Id:          item.Id,
			ArticleBlogName: item.ArticleBlogName,
			Title: item.Title,
			Description: item.Description,
			CategoryId: item.CategoryId,
		}
	}
	totalRecords, _ := m.articleblogRepo.Count(ctx)
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

	response := &models.ArticleBlogWithPagination{
		Data: users,
		Meta: meta,
	}

	return response, nil
}
