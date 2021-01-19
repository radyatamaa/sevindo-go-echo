package models

import (
	"time"
)

type ArticleBlog struct {
	Id              int        `json:"id" validate:"required"`
	CreatedBy       string     `json:"created_by" validate:"required"`
	CreatedDate     time.Time  `json:"created_date" validate:"required"`
	ModifiedBy      *string    `json:"modified_by"`
	ModifiedDate    *time.Time `json:"modified_date"`
	DeletedBy       *string    `json:"deleted_by"`
	DeletedDate     *time.Time `json:"deleted_date"`
	IsDeleted       int        `json:"is_deleted" validate:"required"`
	IsActive        int        `json:"is_active" validate:"required"`
	ArticleBlogName string     `json:"article_blog_name"`
	Title           string     `json:"title" validate:"required"`
	Description     string     `json:"description" validate:"required"`
	CategoryId              string       `json:"category_id"`
	ArticlePicture	string `json:"article_picture"`
}

type NewCommandArticleBlog struct {
	Id              int    `json:"id" validate:"required"`
	ArticleBlogName string `json:"article_blog_name"`
	Title           string `json:"title" validate:"required"`
	Description       string `json:"description" validate:"required"`
	CategoryId              string       `json:"category_id"`
	ArticlePicture	string `json:"article_picture"`
}

type ArticleBlogDto struct {
	Id              int    `json:"id" validate:"required"`
	ArticleBlogName string `json:"article_blog_name"`
	Title           string `json:"title" validate:"required"`
	Description       string `json:"description"`
	CategoryId              string       `json:"category_id"`
	ArticlePicture	string `json:"article_picture"`
}

type ArticleBlogWithPagination struct {
	Data []*ArticleBlogDto  `json:"data"`
	Meta *MetaPagination `json:"meta"`
}
