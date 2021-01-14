package repository

import (
	"context"
	"database/sql"
	"github.com/models"
	"github.com/services/resort_photo"
	"github.com/sirupsen/logrus"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type resortPhotoRepository struct {
	Conn *sql.DB
}


// NewuserRepository will create an object that represent the article.repository interface
func NewresortPhotoRepository(Conn *sql.DB) resort_photo.Repository {
	return &resortPhotoRepository{Conn}
}
func (m *resortPhotoRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.ResortPhoto, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result := make([]*models.ResortPhoto, 0)
	for rows.Next() {
		t := new(models.ResortPhoto)
		err = rows.Scan(
			&t.Id,
			&t.CreatedBy,
			&t.CreatedDate,
			&t.ModifiedBy,
			&t.ModifiedDate,
			&t.DeletedBy,
			&t.DeletedDate,
			&t.IsDeleted,
			&t.IsActive,
			&t.ResortFolder,
			&t.ResortPhotos,
			&t.ResortId,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m resortPhotoRepository) GetByResortID(ctx context.Context, resortId string) ([]*models.ResortPhoto, error) {
	query := `SELECT * FROM resort_photos WHERE is_deleted = 0 and is_active = 1 `

	if resortId != ""{
		query = query + ` AND resort_id = '` + resortId + `' `
	}
	list, err := m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}

	return list, nil
}