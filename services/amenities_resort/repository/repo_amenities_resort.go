package repository

import (
	"context"
	"database/sql"
	"github.com/models"
	"github.com/services/amenities_resort"
	"github.com/sirupsen/logrus"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type amenitiesResortRepository struct {
	Conn *sql.DB
}



// NewuserRepository will create an object that represent the article.repository interface
func NewamenitiesResortRepository(Conn *sql.DB) amenities_resort.Repository {
	return &amenitiesResortRepository{Conn}
}
func (m *amenitiesResortRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.AmenitiesResortJoin, error) {
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

	result := make([]*models.AmenitiesResortJoin, 0)
	for rows.Next() {
		t := new(models.AmenitiesResortJoin)
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
			&t.ResortId	,
			&t.AmenitiesId 		,
			&t.Name 		,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m amenitiesResortRepository) GetByResortId(ctx context.Context, resortId string) ([]*models.AmenitiesResortJoin, error) {
	query := `SELECT rrp.*,c.name FROM 
				amenities_resorts rrp 
			JOIN amenities c ON c.id = rrp.amenities_id 
			WHERE rrp.is_deleted = 0 and rrp.is_active = 1 `

	if resortId != ""{
		query = query + ` and resort_id = '` + resortId + `' `
	}
	list, err := m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}

	return list, nil
}
