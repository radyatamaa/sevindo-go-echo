package models

import (
	"time"
)

type ArticleCategory struct {
	Id                  int        `json:"id" validate:"required"`
	CreatedBy           string     `json:"created_by" validate:"required"`
	CreatedDate         time.Time  `json:"created_date" validate:"required"`
	ModifiedBy          *string    `json:"modified_by"`
	ModifiedDate        *time.Time `json:"modified_date"`
	DeletedBy           *string    `json:"deleted_by"`
	DeletedDate         *time.Time `json:"deleted_date"`
	IsDeleted           int        `json:"is_deleted" validate:"required"`
	IsActive            int        `json:"is_active" validate:"required"`
	ArticleCategoryName string     `json:"article_category_name"`
}

type NewCommandArticleCategory struct {
	Id                  int    `json:"id" validate:"required"`
	ArticleCategoryName string `json:"article_category_name"`
}

type ArticleCategoryDto struct {
	Id                  int    `json:"id" validate:"required"`
	ArticleCategoryName string `json:"article_category_name"`
}
