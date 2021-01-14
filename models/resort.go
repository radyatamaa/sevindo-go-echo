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
	CityId 				int `json:"city_id"`
}
