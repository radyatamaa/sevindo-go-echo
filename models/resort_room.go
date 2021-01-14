package models

import "time"

type ResortRoom struct {
	Id                   string     `json:"id" validate:"required"`
	CreatedBy            string     `json:"created_by" validate:"required"`
	CreatedDate          time.Time  `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int        `json:"is_deleted" validate:"required"`
	IsActive             int        `json:"is_active" validate:"required"`
	ResortRoomTitle				string `json:"resort_title"`
	ResortRoomDesc				string `json:"resort_desc"`
	ResortMaximumBookingAmount int `json:"resort_maximum_booking_amount"`
	ResortImage 	 string `json:"resort_image"`
	ResortId 			string `json:"resort_id"`
}