package models

import "time"

type Review struct {
	Id           string     `json:"id" `
	CreatedBy    string     `json:"created_by"`
	CreatedDate  time.Time  `json:"created_date"`
	ModifiedBy   *string    `json:"modified_by"`
	ModifiedDate *time.Time `json:"modified_date"`
	DeletedBy    *string    `json:"deleted_by"`
	DeletedDate  *time.Time `json:"deleted_date"`
	IsDeleted    int        `json:"is_deleted" validate:"required"`
	IsActive     int        `json:"is_active" validate:"required"`
	Values       float64        `json:"values"`
	Desc         string     `json:"desc"`
	UserId      *string 		`json:"user_id"`
	TransactionId *string `json:"transaction_id"`
}