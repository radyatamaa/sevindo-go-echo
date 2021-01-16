package repository

import (
	"context"
	"database/sql"
	"github.com/models"
	"github.com/services/resort_room"
	"github.com/sirupsen/logrus"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type resortRoomRepository struct {
	Conn *sql.DB
}

// NewuserRepository will create an object that represent the article.repository interface
func NewresortRepository(Conn *sql.DB) resort_room.Repository {
	return &resortRoomRepository{Conn}
}
func (m *resortRoomRepository) fetchJoin(ctx context.Context, query string, args ...interface{}) ([]*models.ResortJoin, error) {
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

	result := make([]*models.ResortJoin, 0)
	for rows.Next() {
		t := new(models.ResortJoin)
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
			&t.ResortTitle		,
			&t.ResortDesc		,
			&t.ResortLongitude		,
			&t.ResortLatitude 	,
			&t.Status 			,
			&t.	Rating 		,
			&t.	BranchId 		,
			&t.	DistrictsId 		,
			&t.	DistrictsName 		,
			&t.	CityId			,
			&t.	CityName 		,
			&t.	ProvinceId 		,
			&t.	ProvinceName 	,
			&t.	BranchName 		,
			&t.	Price 			,
			&t.	Currency 			,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (r resortRoomRepository) GetByResortIdJoinWithPayment(ctx context.Context, resortID string) ([]*models.ResortRoom, error) {
	panic("implement me")
}