package models

import "time"

type Resort struct {
	Id                   string     `json:"id" validate:"required"`
	CreatedBy            string     `json:"created_by" validate:"required"`
	CreatedDate          time.Time  `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int        `json:"is_deleted" validate:"required"`
	IsActive             int        `json:"is_active" validate:"required"`
	ResortTitle				string `json:"resort_title"`
	ResortDesc				string `json:"resort_desc"`
	ResortLongitude			float64 `json:"resort_longitude"`
	ResortLatitude 			float64 `json:"resort_latitude"`
	Status 				int `json:"status"`
	Rating 				float64 `json:"rating"`
	BranchId 			string `json:"branch_id"`
	DistrictsId 				int `json:"districts_id"`
}

type ResortJoin struct {
	Id                   string     `json:"id" validate:"required"`
	CreatedBy            string     `json:"created_by" validate:"required"`
	CreatedDate          time.Time  `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int        `json:"is_deleted" validate:"required"`
	IsActive             int        `json:"is_active" validate:"required"`
	ResortTitle				string `json:"resort_title"`
	ResortDesc				string `json:"resort_desc"`
	ResortLongitude			float64 `json:"resort_longitude"`
	ResortLatitude 			float64 `json:"resort_latitude"`
	Status 				int `json:"status"`
	Rating 				float64 `json:"rating"`
	BranchId 			string `json:"branch_id"`
	DistrictsId 				int `json:"districts_id"`
	DistrictsName 		string `json:"districts_name"`
	CityId				int `json:"city_id"`
	CityName 			string `json:"city_name"`
	ProvinceId 			string `json:"province_id"`
	ProvinceName 		string `json:"province_name"`
	BranchName 			string `json:"branch_name"`
	Price 				float64 `json:"price"`
	Currency 			string `json:"currency"`
}
type ResortJoinDto struct {
	Id                   string     `json:"id" validate:"required"`
	ResortTitle				string `json:"resort_title"`
	ResortDesc				string `json:"resort_desc"`
	ResortLongitude			float64 `json:"resort_longitude"`
	ResortLatitude 			float64 `json:"resort_latitude"`
	Status 				int `json:"status"`
	Rating 				float64 `json:"rating"`
	BranchId 			string `json:"branch_id"`
	DistrictsId 				int `json:"districts_id"`
	DistrictsName 		string `json:"districts_name"`
	CityId				int `json:"city_id"`
	CityName 			string `json:"city_name"`
	ProvinceId 			string `json:"province_id"`
	ProvinceName 		string `json:"province_name"`
	BranchName 			string `json:"branch_name"`
	Price 				float64 `json:"price"`
	ResortPhoto 		[]string `json:"resort_photo"`
}

type ResortJoinDtoWithPagination struct {
	Data []*ResortJoinDto  `json:"data"`
	Meta *MetaPagination `json:"meta"`
}