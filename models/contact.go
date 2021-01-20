package models

import "time"

type Contact struct {
	Id           string     `json:"id" validate:"required"`
	CreatedBy    string     `json:"created_by" validate:"required"`
	CreatedDate  time.Time  `json:"created_date" validate:"required"`
	ModifiedBy   *string    `json:"modified_by"`
	ModifiedDate *time.Time `json:"modified_date"`
	DeletedBy    *string    `json:"deleted_by"`
	DeletedDate  *time.Time `json:"deleted_date"`
	IsDeleted    int        `json:"is_deleted" validate:"required"`
	IsActive     int        `json:"is_active" validate:"required"`
	FullName string `json:"full_name"`
	TypeAs  int `json:"type_as"`
	PhoneNumber string `json:"phone_number"`
	UserId 		string `json:"user_id"`
}


type NewCommandContact struct {
	Id          string `json:"id" validate:"required"`
	FullName string `json:"full_name"`
	TypeAs  int `json:"type_as"`
	PhoneNumber string `json:"phone_number"`
}

type ContactDto struct {
	Id          string `json:"id" validate:"required"`
	FullName string `json:"full_name"`
	TypeAs  int `json:"type_as"`
	PhoneNumber string `json:"phone_number"`
	UserId 		string `json:"user_id"`
}

type ContactWithPagination struct {
	Data []*ContactDto  `json:"data"`
	Meta *MetaPagination `json:"meta"`
}
