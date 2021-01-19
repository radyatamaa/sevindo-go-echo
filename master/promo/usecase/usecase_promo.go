package usecase

import (
	"context"
	"github.com/auth/user_admin"
	"github.com/google/uuid"
	"github.com/master/promo"
	"math"
	"time"

	"github.com/models"
)

type promoUsecase struct {
	userAdminUsecase user_admin.Usecase
	promoRepo promo.Repository
	contextTimeout time.Duration
}



// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewpromoUsecase(	userAdminUsecase user_admin.Usecase,promoRepo promo.Repository, timeout time.Duration) promo.Usecase {
	return &promoUsecase{
		userAdminUsecase:userAdminUsecase,
		promoRepo:    promoRepo,
		contextTimeout: timeout,
	}
}
func (m promoUsecase) Delete(c context.Context, id string, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return nil,models.ErrUnAuthorize
	}

	err = m.promoRepo.Delete(ctx,id,currentUser.Email)

	result := &models.ResponseDelete{
		Id:      id,
		Message: "Success Delete",
	}

	return result,nil
}

func (m promoUsecase) Update(c context.Context, ar *models.NewCommandPromo, token string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return models.ErrUnAuthorize
	}

	getPromo ,err := m.promoRepo.GetByID(ctx,ar.Id)
	if err != nil{
		return err
	}
	var modifyBy string = currentUser.Email
	now := time.Now()
	getPromo.PromoCode = ar.PromoCode
	getPromo.PromoName = ar.PromoName
	getPromo.PromoDesc = ar.PromoDesc
	getPromo.PromoValue = ar.PromoValue
	getPromo.PromoType = ar.PromoType
	getPromo.PromoImage = ar.PromoImage
	getPromo.StartDate = ar.StartDate
	getPromo.EndDate = ar.EndDate
	getPromo.HowToGet = ar.HowToGet
	getPromo.HowToUse = ar.HowToUse
	getPromo.TermCondition = ar.TermCondition
	getPromo.Disclaimer = ar.Disclaimer
	getPromo.MaxDiscount = ar.MaxDiscount
	getPromo.MaxUsage = ar.MaxUsage
	getPromo.ProductionCapacity = ar.ProductionCapacity
	getPromo.CurrencyId = ar.CurrencyId
	getPromo.ModifiedBy = &modifyBy
	getPromo.ModifiedDate = &now
	err = m.promoRepo.Update(ctx,getPromo)
	if err != nil{
		return err
	}
	return nil
}

func (m promoUsecase) List(ctx context.Context, page, limit, offset int, search string) (*models.PromoWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.promoRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	users := make([]*models.PromoDto, len(list))
	for i, item := range list {
		users[i] = &models.PromoDto{
			Id:          item.Id,
			PromoCode: item.PromoCode,
			PromoName: item.PromoName,
			PromoDesc: item.PromoDesc,
			PromoValue: item.PromoValue,
			PromoType: item.PromoType,
			PromoImage: item.PromoImage,
			StartDate: item.StartDate,
			EndDate: item.EndDate,
			HowToGet: item.HowToGet,
			HowToUse: item.HowToUse,
			TermCondition: item.TermCondition,
			Disclaimer: item.Disclaimer,
			MaxDiscount: item.MaxDiscount,
			MaxUsage: item.MaxUsage,
			ProductionCapacity: item.ProductionCapacity,
			CurrencyId: item.CurrencyId,
		}
	}
	totalRecords, _ := m.promoRepo.Count(ctx)
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

	response := &models.PromoWithPagination{
		Data: users,
		Meta: meta,
	}

	return response, nil
}

func (m promoUsecase) Create(c context.Context, ar *models.NewCommandPromo, token string) (*models.NewCommandPromo, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return nil,models.ErrUnAuthorize
	}

	insert := models.Promo{
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
		PromoCode: ar.PromoCode,
		PromoName: ar.PromoName,
		PromoDesc: ar.PromoDesc,
		PromoValue: ar.PromoValue,
		PromoType: ar.PromoType,
		PromoImage: ar.PromoImage,
		StartDate: ar.StartDate,
		EndDate: ar.EndDate,
		HowToGet: ar.HowToGet,
		HowToUse: ar.HowToUse,
		TermCondition: ar.TermCondition,
		Disclaimer: ar.Disclaimer,
		MaxDiscount: ar.MaxDiscount,
		MaxUsage: ar.MaxUsage,
		ProductionCapacity: ar.ProductionCapacity,
		CurrencyId: ar.CurrencyId,
	}

	err = m.promoRepo.Insert(ctx, &insert)
	if err != nil {
		return nil, err
	}
	ar.Id = insert.Id
	return ar, nil
}

func (m promoUsecase) GetById(c context.Context, id string, token string) (*models.PromoDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	promo, err := m.promoRepo.GetByID(ctx, id)
	if err != nil {
		return nil, models.ErrNotFound
	}

	result := &models.PromoDto{
		Id:          promo.Id,
		PromoCode: promo.PromoCode,
		PromoName: promo.PromoName,
		PromoDesc: promo.PromoDesc,
		PromoValue: promo.PromoValue,
		PromoType: promo.PromoType,
		PromoImage: promo.PromoImage,
		StartDate: promo.StartDate,
		EndDate: promo.EndDate,
		HowToGet: promo.HowToGet,
		HowToUse: promo.HowToUse,
		TermCondition: promo.TermCondition,
		Disclaimer: promo.Disclaimer,
		MaxDiscount: promo.MaxDiscount,
		MaxUsage: promo.MaxUsage,
		ProductionCapacity: promo.ProductionCapacity,
		CurrencyId: promo.CurrencyId,
	}

	return result, nil
}
