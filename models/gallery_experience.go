package models

import "time"

type GalleryExperience struct {
	Id                   string     `json:"id" validate:"required"`
	CreatedBy            string     `json:"created_by" validate:"required"`
	CreatedDate          time.Time  `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int        `json:"is_deleted" validate:"required"`
	IsActive             int        `json:"is_active" validate:"required"`
	ExperienceName 		string `json:"experience_name"`
	ExperienceDesc 			string `json:"experience_desc"`
	ExperiencePicture 		*string `json:"experience_picture"`
	Longitude			float64 `json:"longitude"`
	Latitude 			float64 `json:"latitude"`
}

type NewCommandGalleryExperience struct {
	Id          string `json:"id" validate:"required"`
	ExperienceName    string     `json:"experience_name"`
	ExperienceDesc    string     `json:"experience_desc"`
	ExperiencePicture   *string     `json:"experience_picture"`
	Longitude  float64    `json:"longitude"`
	Latitude  float64    `json:"latitude"`
}

type GalleryExperienceDto struct {
	Id          string `json:"id" validate:"required"`
	ExperienceName    string     `json:"experience_name"`
	ExperienceDesc    string     `json:"experience_desc"`
	ExperiencePicture   *string     `json:"experience_picture"`
	Longitude  float64    `json:"longitude"`
	Latitude  float64    `json:"latitude"`
}
type GalleryExperienceWithPagination struct {
	Data []*GalleryExperienceDto  `json:"data"`
	Meta *MetaPagination `json:"meta"`
}
