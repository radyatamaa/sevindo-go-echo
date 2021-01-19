package repository

import (
	"context"
	"database/sql"
	"github.com/models"
	"github.com/services/review"
	"github.com/sirupsen/logrus"
	"strconv"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type reviewRepository struct {
	Conn *sql.DB
}



// NewuserRepository will create an object that represent the article.repository interface
func NewreviewRepository(Conn *sql.DB) review.Repository {
	return &reviewRepository{Conn}
}

func (m *reviewRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.ReviewJoin, error) {
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

	result := make([]*models.ReviewJoin, 0)
	for rows.Next() {
		t := new(models.ReviewJoin)
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
			&t.Values ,
			&t.Desc    ,
			&t.UserId   ,
			&t.TransactionId ,
			&t.Name 	,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (m *reviewRepository) GetByResortIdJoinWithPayment(ctx context.Context, resortID string,limit int ,offset int) ([]*models.ReviewJoin, error) {
	query := `SELECT r.*,u.first_name as name FROM 
				reviews r
			JOIN users u ON u.id = r.user_id
			JOIN transactions t ON t.id = r.transaction_id
			JOIN bookings b ON b.id = t.booking_id
			WHERE r.is_deleted = 0 and r.is_active = 1 `

	if resortID != ""{
		query = query + ` and b.resort_id = '` + resortID + `' `
	}
	if limit != 0{
		query = query + ` LIMIT `+strconv.Itoa(limit)+` OFFSET `+ strconv.Itoa(offset) + ` `
	}
	list, err := m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (m *reviewRepository) Count(ctx context.Context, resortID string) (int, error) {
	query := `SELECT COUNT(*) AS count FROM 
				reviews r
			JOIN users u ON u.id = r.user_id
			JOIN transactions t ON t.id = r.transaction_id
			JOIN bookings b ON b.id = t.booking_id
			WHERE r.is_deleted = 0 and r.is_active = 1 `

	if resortID != ""{
		query = query + ` and b.resort_id = '` + resortID + `' `
	}
	rows, err := m.Conn.QueryContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	count, err := checkCount(rows)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	return count, nil
}

func checkCount(rows *sql.Rows) (count int, err error) {
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}
