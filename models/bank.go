package models

import (
	"time"
)

type Bank struct {
	Id           int        `json:"id" validate:"required"`
	CreatedBy    string     `json:"created_by" validate:"required"`
	CreatedDate  time.Time  `json:"created_date" validate:"required"`
	ModifiedBy   *string    `json:"modified_by"`
	ModifiedDate *time.Time `json:"modified_date"`
	DeletedBy    *string    `json:"deleted_by"`
	DeletedDate  *time.Time `json:"deleted_date"`
	IsDeleted    int        `json:"is_deleted" validate:"required"`
	IsActive     int        `json:"is_active" validate:"required"`
	BankName string     `json:"bank_name"`
}

type NewCommandBank struct {
	Id           int    `json:"id" validate:"required"`
	BankName string `json:"bank_name"`
}

type BankDto struct {
	Id           int    `json:"id" validate:"required"`
	BankName string `json:"bank_name"`
}

type BankWithPagination struct {
	Data []*BankDto  `json:"data"`
	Meta *MetaPagination `json:"meta"`
}

