package usecase

import (
	"context"
	"github.com/auth/user_admin"
	"github.com/google/uuid"
	"github.com/master/gallery_experience"

	"math"
	"time"

	"github.com/models"
)

type galerryexperienceUsecase struct {
	userAdminUsecase user_admin.Usecase
	galleryexperienceRepo gallery_experience.Repository
	contextTimeout time.Duration
}



// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewGalleyExperienceUsecase(	userAdminUsecase user_admin.Usecase,galleryexperienceRepo gallery_experience.Repository, timeout time.Duration) gallery_experience.Usecase {
	return &galerryexperienceUsecase{
		userAdminUsecase:userAdminUsecase,
		galleryexperienceRepo:galleryexperienceRepo,
		contextTimeout: timeout,
	}
}
func (m galerryexperienceUsecase) Delete(c context.Context, id string, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return nil,models.ErrUnAuthorize
	}

	err = m.galleryexperienceRepo.Delete(ctx,id,currentUser.Email)

	result := &models.ResponseDelete{
		Id:      id,
		Message: "Success Delete",
	}

	return result,nil
}

func (m galerryexperienceUsecase) Update(c context.Context, ar *models.NewCommandGalleryExperience, token string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return models.ErrUnAuthorize
	}

	getPromo ,err := m.galleryexperienceRepo.GetByID(ctx,ar.Id)
	if err != nil{
		return err
	}
	var modifyBy string = currentUser.Email
	now := time.Now()
	getPromo.ExperienceName = ar.ExperienceName
	getPromo.ExperienceDesc = ar.ExperienceDesc
	getPromo.ExperiencePicture = ar.ExperiencePicture
	getPromo.Longitude = ar.Longitude
	getPromo.Latitude = ar.Latitude
	getPromo.ModifiedBy = &modifyBy
	getPromo.ModifiedDate = &now
	err = m.galleryexperienceRepo.Update(ctx,getPromo)
	if err != nil{
		return err
	}
	return nil
}

func (m galerryexperienceUsecase) List(ctx context.Context, page, limit, offset int, search string) (*models.GalleryExperienceWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.galleryexperienceRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	users := make([]*models.GalleryExperienceDto, len(list))
	for i, item := range list {
		users[i] = &models.GalleryExperienceDto{
			Id:          item.Id,
			ExperienceName: item.ExperienceName,
			ExperienceDesc: item.ExperienceDesc,
			ExperiencePicture: item.ExperiencePicture,
			Longitude: item.Longitude,
			Latitude: item.Latitude,
		}
	}
	totalRecords, _ := m.galleryexperienceRepo.Count(ctx)
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

	response := &models.GalleryExperienceWithPagination{
		Data: users,
		Meta: meta,
	}

	return response, nil
}

func (m galerryexperienceUsecase) Create(c context.Context, ar *models.NewCommandGalleryExperience, token string) (*models.NewCommandGalleryExperience, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return nil,models.ErrUnAuthorize
	}

	insert := models.GalleryExperience{
		//Id:           0,
		Id:            uuid.New().String(),
		CreatedBy:    currentUser.Email,
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		ExperienceName: ar.ExperienceName,
		ExperienceDesc: ar.ExperienceDesc,
		ExperiencePicture: ar.ExperiencePicture,
		Longitude: ar.Longitude,
		Latitude: ar.Latitude,
	}

	err = m.galleryexperienceRepo.Insert(ctx, &insert)
	if err != nil {
		return nil, err
	}
	ar.Id = insert.Id
	return ar, nil
}

func (m galerryexperienceUsecase) GetById(c context.Context, id string, token string) (*models.GalleryExperienceDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	promo, err := m.galleryexperienceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, models.ErrNotFound
	}

	result := &models.GalleryExperienceDto{
		Id:          promo.Id,
		ExperienceName: promo.ExperienceName,
		ExperienceDesc: promo.ExperienceDesc,
		ExperiencePicture: promo.ExperiencePicture,
		Longitude: promo.Longitude,
		Latitude: promo.Latitude,
	}

	return result, nil
}
