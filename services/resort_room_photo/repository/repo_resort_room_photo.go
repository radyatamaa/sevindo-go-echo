package repository

import (
	"context"
	"database/sql"
	"github.com/models"
	"github.com/services/resort_room_photo"
	"github.com/sirupsen/logrus"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type resortRoomPhotoRepository struct {
	Conn *sql.DB
}



// NewuserRepository will create an object that represent the article.repository interface
func NewresortRoomPhotoRepository(Conn *sql.DB) resort_room_photo.Repository {
	return &resortRoomPhotoRepository{Conn}
}

func (m *resortRoomPhotoRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.ResortRoomPhoto, error) {
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

	result := make([]*models.ResortRoomPhoto, 0)
	for rows.Next() {
		t := new(models.ResortRoomPhoto)
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
			&t.ResortFolder ,
			&t.ResortPhotos ,
			&t.ResortRoomId 	,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m resortRoomPhotoRepository) GetByResortRoomID(ctx context.Context, resortRoomId string) ([]*models.ResortRoomPhoto, error) {
	query := `SELECT * FROM 
				resort_room_photos rrp 
			WHERE rrp.is_deleted = 0 and rrp.is_active = 1 `

	if resortRoomId != ""{
		query = query + ` and resort_room_id = '` + resortRoomId + `' `
	}
	list, err := m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}

	return list, nil
}
