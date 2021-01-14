package models

import "time"

type Role struct {
	Id           int     `json:"id" validate:"required"`
	CreatedBy    string     `json:"created_by" validate:"required"`
	CreatedDate  time.Time  `json:"created_date" validate:"required"`
	ModifiedBy   *string    `json:"modified_by"`
	ModifiedDate *time.Time `json:"modified_date"`
	DeletedBy    *string    `json:"deleted_by"`
	DeletedDate  *time.Time `json:"deleted_date"`
	IsDeleted    int        `json:"is_deleted" validate:"required"`
	IsActive     int        `json:"is_active" validate:"required"`
	RoleName  	 string     `json:"role_name"`
	RoleType     int        `json:"role_type"`
	Description  string     `json:"description"`
}

type NewCommandRole struct {
	Id          int `json:"id" validate:"required"`
	RoleName    string `json:"role_name"`
	RoleType     int        `json:"role_type"`
	Description  string     `json:"description"`
}

type RoleDto struct {
	Id          int `json:"id" validate:"required"`
	RoleName    string `json:"role_name"`
	RoleType     int        `json:"role_type"`
	Description  string     `json:"description"`
}

type RoleWithPagination struct {
	Data []*RoleDto  `json:"data"`
	Meta *MetaPagination `json:"meta"`
}