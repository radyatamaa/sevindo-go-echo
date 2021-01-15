package usecase

import (
	"context"
	"encoding/json"
	"github.com/services/accessibility_resort"
	"github.com/services/amenities_resort"
	"github.com/services/resort"
	"github.com/services/resort_photo"
	"github.com/services/resort_room"
	"github.com/services/resort_room_payment"
	"github.com/services/resort_room_photo"
	"math"
	"time"

	"github.com/models"
)

type resortUsecase struct {
	accessibilityResortRepo accessibility_resort.Repository
	amenitiesResortRepo amenities_resort.Repository
	resortRoomPhoto resort_room_photo.Repository
	resortPhoto resort_photo.Repository
	resortRepo    resort.Repository
	resortRoomRepo resort_room.Repository
	resortRoomPaymentRepo resort_room_payment.Repository
	contextTimeout time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewresortUsecase(		accessibilityResortRepo accessibility_resort.Repository,
amenitiesResortRepo amenities_resort.Repository,resortRoomPhoto resort_room_photo.Repository,resortRoomPaymentRepo resort_room_payment.Repository,resortRoomRepo resort_room.Repository,resortPhoto resort_photo.Repository,resortRepo resort.Repository, timeout time.Duration) resort.Usecase {
	return &resortUsecase{
		accessibilityResortRepo:accessibilityResortRepo,
		amenitiesResortRepo:amenitiesResortRepo,
		resortRoomPhoto:resortRoomPhoto,
		resortPhoto:resortPhoto,
		resortRepo:    resortRepo,
		resortRoomPaymentRepo:resortRoomPaymentRepo,
		resortRoomRepo:resortRoomRepo,
		contextTimeout: timeout,
	}
}
func (m resortUsecase) GetDetail(ctx context.Context, id string) (*models.ResortJoinDetailDto, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.resortRepo.GetById(ctx,id)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return  nil,models.ErrNotFound
	}

	result := &models.ResortJoinDetailDto{
		Id:              list[0].Id,
		ResortTitle:     list[0].ResortTitle,
		ResortDesc:      list[0].ResortDesc,
		ResortLongitude: list[0].ResortLongitude,
		ResortLatitude:  list[0].ResortLatitude,
		Status:          list[0].Status,
		Rating:          list[0].Rating,
		BranchId:        list[0].BranchId,
		DistrictsId:     list[0].DistrictsId,
		DistrictsName:   list[0].DistrictsName,
		CityId:          list[0].CityId,
		CityName:        list[0].CityName,
		ProvinceId:      list[0].ProvinceId,
		ProvinceName:    list[0].ProvinceName,
		BranchName:      list[0].BranchName,
		ResortRoom:      nil,
		ResortPhoto:     nil,
	}
	for _,element := range list{
		rm := models.ResortRoomObj{
			Id:                         element.ResortRoomId,
			ResortRoomTitle:            element.ResortRoomTitle,
			ResortRoomDesc:             element.ResortRoomDesc,
			ResortMaximumBookingAmount: element.ResortMaximumBookingAmount,
			ResortCapacity:             element.ResortCapacity,
			ResortRoomPayment:          nil,
			ResortRoomPhoto:            nil,
		}
		resortRoomPayment ,err := m.resortRoomPaymentRepo.GetByResortRoomID(ctx ,rm.Id)
		if err != nil {
			return nil, err
		}
		//for _,p := range resortRoomPayment{
		if len(resortRoomPayment) > 0 {
			rm.ResortRoomPayment = &models.ResortRoomPaymentObj{
				Id:       resortRoomPayment[0].Id,
				Currency: resortRoomPayment[0].CurrencyName,
				Price:    resortRoomPayment[0].Price,
			}
		}
		//}

		resortRoomPhoto ,err := m.resortRoomPhoto.GetByResortRoomID(ctx,rm.Id)
		for _,photo := range resortRoomPhoto{
			images := make([]string, 0)
			if photo.ResortPhotos != "" && photo.ResortPhotos != "[]" {
				if errUnmarshal := json.Unmarshal([]byte(photo.ResortPhotos), &images); errUnmarshal != nil {
					return nil, errUnmarshal
				}
			}

			pt := models.ResortRoomPhotoObj{
				ResortFolder: photo.ResortFolder,
				ResortPhotos: images,
			}

			rm.ResortRoomPhoto = append(rm.ResortRoomPhoto ,pt)
		}

		result.ResortRoom = append(result.ResortRoom,rm)

	}
	resortPhoto ,err := m.resortPhoto.GetByResortID(ctx,list[0].Id)
	if err != nil {
		return nil, err
	}
	for _, itemP := range resortPhoto {
		images := make([]string, 0)
		if itemP.ResortPhotos != "" && itemP.ResortPhotos != "[]" {
			if errUnmarshal := json.Unmarshal([]byte(itemP.ResortPhotos), &images); errUnmarshal != nil {
				return nil, errUnmarshal
			}
		}

		photo := models.ResortRoomPhotoObj{
			ResortFolder: itemP.ResortFolder,
			ResortPhotos: images,
		}

		result.ResortPhoto = append(result.ResortPhoto,photo)
	}

	return result,nil
}
func (m resortUsecase) GetAll(ctx context.Context, page, limit, offset int,capacity int, startDate string, endDate string) (*models.ResortJoinDtoWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.resortRepo.GetAll(ctx, make([]string,0),capacity,limit, offset)
	if err != nil {
		return nil, err
	}

	resort := make([]*models.ResortJoinDto, len(list))
	for i, item := range list {
		resort[i] = &models.ResortJoinDto{
			Id:              item.Id,
			ResortTitle:     item.ResortTitle,
			ResortDesc:      item.ResortDesc,
			ResortLongitude: item.ResortLongitude,
			ResortLatitude:  item.ResortLatitude,
			Status:          item.Status,
			Rating:          item.Rating,
			BranchId:        item.BranchId,
			DistrictsId:     item.DistrictsId,
			DistrictsName:   item.DistrictsName,
			CityId:          item.CityId,
			CityName:        item.CityName,
			ProvinceId:      item.ProvinceId,
			ProvinceName:    item.ProvinceName,
			BranchName:      item.BranchName,
			Price:           item.Price,
			ResortPhoto:     make([]string,0),
		}
		resortPhoto ,err := m.resortPhoto.GetByResortID(ctx,item.Id)
		if err != nil {
			return nil, err
		}
		for _, itemP := range resortPhoto {
			images := make([]string, 0)
			if itemP.ResortPhotos != "" && itemP.ResortPhotos != "[]" {
				if errUnmarshal := json.Unmarshal([]byte(itemP.ResortPhotos), &images); errUnmarshal != nil {
					return nil, errUnmarshal
				}
			}
			resort[i].ResortPhoto = append(resort[i].ResortPhoto,images...)
		}
	}
	totalRecords, _ := m.resortRepo.GetAllCount(ctx,make([]string,0),capacity)
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

	response := &models.ResortJoinDtoWithPagination{
		Data: resort,
		Meta: meta,
	}

	return response, nil
}