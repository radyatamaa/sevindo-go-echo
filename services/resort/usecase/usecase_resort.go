package usecase

import (
	"context"
	"encoding/json"
	"github.com/services/resort"
	"github.com/services/resort_photo"
	"math"
	"time"

	"github.com/models"
)

type resortUsecase struct {
	resortPhoto resort_photo.Repository
	resortRepo    resort.Repository
	contextTimeout time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewresortUsecase(resortPhoto resort_photo.Repository,resortRepo resort.Repository, timeout time.Duration) resort.Usecase {
	return &resortUsecase{
		resortPhoto:resortPhoto,
		resortRepo:    resortRepo,
		contextTimeout: timeout,
	}
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