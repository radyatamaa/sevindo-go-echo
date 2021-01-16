package repository

import (
	"context"
	"database/sql"
	"github.com/models"
	"github.com/services/resort_room_payment"
	"github.com/sirupsen/logrus"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type resortRoomPaymentRepository struct {
	Conn *sql.DB
}



// NewuserRepository will create an object that represent the article.repository interface
func NewresortRoomPaymentRepository(Conn *sql.DB) resort_room_payment.Repository {
	return &resortRoomPaymentRepository{Conn}
}
func (m *resortRoomPaymentRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.ResortRoomPaymentJoin, error) {
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

	result := make([]*models.ResortRoomPaymentJoin, 0)
	for rows.Next() {
		t := new(models.ResortRoomPaymentJoin)
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
			&t.Currency 	,
			&t.Price 		,
			&t.ResortRoomId	,
			&t.CurrencyName 		,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m resortRoomPaymentRepository) GetByResortRoomID(ctx context.Context, resortRoomId string) ([]*models.ResortRoomPaymentJoin, error) {
	query := `SELECT rrp.*,c.currency_name FROM 
				resort_room_payments rrp 
			JOIN currencies c ON c.id = rrp.currency 
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
