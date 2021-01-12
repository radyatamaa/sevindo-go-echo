package repository

import (
	"context"
	"database/sql"

	"time"

	"github.com/sirupsen/logrus"

	"github.com/master/country"
	"github.com/models"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type countryRepository struct {
	Conn *sql.DB
}



// NewuserRepository will create an object that represent the article.repository interface
func NewCountryRepository(Conn *sql.DB) country.Repository {
	return &countryRepository{Conn}
}
func (m *countryRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Country, error) {
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

	result := make([]*models.Country, 0)
	for rows.Next() {
		t := new(models.Country)
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
			&t.CountryName,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (m *countryRepository) GetByID(ctx context.Context, id string) (res *models.Country, err error) {
	query := `SELECT * FROM countries WHERE `

	if id != "" {
		query = query + ` id = '` + id + `' `
	}

	list, err := m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return nil, models.ErrNotFound
	}

	return
}

func (m *countryRepository) Update(ctx context.Context, ar *models.Country) error {
	panic("implement me")
}

func (m *countryRepository) Delete(ctx context.Context, id string, deleted_by string) error {
	panic("implement me")
}


func (m *countryRepository) Insert(ctx context.Context, a *models.Country) error {
	query := `INSERT countries SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , deleted_date=? , is_deleted=? , is_active=? ,
	country_name=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.CountryName)
	if err != nil {
		return err
	}

	//lastID, err := res.RowsAffected()
	if err != nil {
		return err
	}

	//a.Id = lastID
	return nil
}

func (m *countryRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT count(*) AS count FROM countries WHERE is_deleted = 0 and is_active = 1`

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

func (m *countryRepository) List(ctx context.Context, limit, offset int) ([]*models.Country, error) {
	

	return nil, nil
}
