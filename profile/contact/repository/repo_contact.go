package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/models"
	"github.com/profile/contact"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type contactRepository struct {
	Conn *sql.DB
}

// NewuserRepository will create an object that represent the article.repository interface
func NewAcontactRepository(Conn *sql.DB) contact.Repository {
	return &contactRepository{Conn}
}
func (m *contactRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Contact, error) {
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

	result := make([]*models.Contact, 0)
	for rows.Next() {
		t := new(models.Contact)
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
			&t.FullName ,
			&t.TypeAs ,
			&t.PhoneNumber ,
			&t.UserId 	,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *contactRepository) GetByID(ctx context.Context, id string) (res *models.Contact, err error) {
	query := `SELECT * FROM contacts WHERE `

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

func (m *contactRepository) Update(ctx context.Context, a *models.Contact) error {
	query := `UPDATE contacts set modified_by=?, modified_date=? , full_name=?,type_as=?,phone_number=?,user_id=?  WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	res, err := stmt.ExecContext(ctx, a.ModifiedBy, time.Now(), a.FullName,a.TypeAs,a.PhoneNumber,a.UserId, a.Id)
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

func (m *contactRepository) Delete(ctx context.Context, id string, deleted_by string) error {
	query := `UPDATE contacts SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=? WHERE id =?`
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

func (m *contactRepository) Insert(ctx context.Context, a *models.Contact) error {
	query := `INSERT contacts SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , deleted_date=? , is_deleted=? , is_active=? ,
	 			full_name=?,type_as=?,phone_number=?,user_id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.FullName,a.TypeAs,a.PhoneNumber,a.UserId)
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

func (m *contactRepository) Count(ctx context.Context,userId string) (int, error) {
	query := `SELECT count(*) AS count FROM contacts WHERE is_deleted = 0 and is_active = 1`
	if userId != "" {
		query = query + ` AND user_id = '` + userId + `' `
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

func (m *contactRepository) List(ctx context.Context, limit, offset int,userId string) ([]*models.Contact, error) {
	query := `SELECT * FROM contacts WHERE is_deleted = 0 and is_active = 1 `
	if userId != "" {
		query = query + ` AND user_id = '` + userId + `' `
	}
	query = query + ` LIMIT ? OFFSET ?`
	list, err := m.fetch(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return list, nil
}


