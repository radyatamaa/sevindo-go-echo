package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/master/accessibility"
	"github.com/models"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type accessibilityRepository struct {
	Conn *sql.DB
}



// NewuserRepository will create an object that represent the article.repository interface
func NewAccessibilityRepository(Conn *sql.DB) accessibility.Repository {
	return &accessibilityRepository{Conn}
}
func (m *accessibilityRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Accessibility, error) {
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

	result := make([]*models.Accessibility, 0)
	for rows.Next() {
		t := new(models.Accessibility)
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
			&t.Name,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (m *accessibilityRepository) GetByID(ctx context.Context, id int) (res *models.Accessibility, err error) {
	query := `SELECT * FROM accessibilities WHERE `

	if id != 0 {
		query = query + ` id = '` + strconv.Itoa(id) + `' `
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

func (m *accessibilityRepository) Update(ctx context.Context, a *models.Accessibility) error {
	query := `UPDATE accessibilities set modified_by=?, modified_date=? , name=?  WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	res, err := stmt.ExecContext(ctx, a.ModifiedBy, time.Now(), a.Name, a.Id)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affect != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", affect)

		return err
	}

	return nil
}

func (m *accessibilityRepository) Delete(ctx context.Context, id int, deleted_by string) error {
	query := `UPDATE accessibilities SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=? WHERE id =?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, deleted_by, time.Now(), 1, 0, id)
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


func (m *accessibilityRepository) Insert(ctx context.Context, a *models.Accessibility) error {
	query := `INSERT accessibilities SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , deleted_date=? , is_deleted=? , is_active=? ,
	name=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.Name)
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

func (m *accessibilityRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT count(*) AS count FROM accessibilities WHERE is_deleted = 0 and is_active = 1`

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

func (m *accessibilityRepository) List(ctx context.Context, limit, offset int) ([]*models.Accessibility, error) {
	query := `SELECT * FROM accessibilities WHERE is_deleted = 0 and is_active = 1 `

	query = query + ` LIMIT ? OFFSET ?`
	list, err := m.fetch(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return list, nil
}