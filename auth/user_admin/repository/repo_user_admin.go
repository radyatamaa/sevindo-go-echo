package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"

	"time"

	"github.com/sirupsen/logrus"

	"github.com/auth/user_admin"
	"github.com/models"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type userAdminRepository struct {
	Conn *sql.DB
}

// NewuserAdminRepository will create an object that represent the article.repository interface
func NewuserAdminRepository(Conn *sql.DB) user_admin.Repository {
	return &userAdminRepository{Conn}
}

func (m *userAdminRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT count(*) AS count FROM user_admins WHERE is_deleted = 0 and is_active = 1`

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

func (m *userAdminRepository) List(ctx context.Context, limit, offset int, search string) ([]*models.UserAdmin, error) {
	query := `SELECT * FROM user_admins WHERE is_deleted = 0 and is_active = 1 `
	if search != "" {
		//query = query + ` AND ( email LIKE '%` + search + `%' ` +
		//	`OR full_name LIKE '%` + search + `%' ` +
		//	`OR phone_number LIKE '%` + search + `%' ` +
		//	`OR address LIKE '%` + search + `%' ` +
		//	`OR dob LIKE '%` + search + `%' ` +
		//	`OR points LIKE '%` + search + `%' )`
	}
	query = query + ` LIMIT ? OFFSET ?`
	list, err := m.fetch(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (m *userAdminRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.UserAdmin, error) {
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

	result := make([]*models.UserAdmin, 0)
	for rows.Next() {
		t := new(models.UserAdmin)
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
			&t.Email,
			&t.FullName,
			&t.BranchId,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *userAdminRepository) GetByID(ctx context.Context, id string) (res *models.UserAdmin, err error) {
	query := `SELECT * FROM user_admins WHERE `

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

func (m *userAdminRepository) GetByUserEmail(ctx context.Context, userEmail string,isAdmin bool) (res *models.UserAdmin, err error) {
	query := `SELECT * FROM user_admins WHERE is_deleted = 0 AND is_active = 1 `

	if userEmail != "" {
		query = query + ` AND email = '` + userEmail + `' `
	}
	if isAdmin != true{
		query = query + ` AND branch_id is null `
	}
	list, err := m.fetch(ctx, query)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return nil, nil
	}
	return
}

func (m *userAdminRepository) Insert(ctx context.Context, a *models.UserAdmin) error {
	query := `INSERT user_admins SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , deleted_date=? , 
			is_deleted=? , is_active=? , email=? , full_name=? , branch_id=? `
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.Email, a.FullName, a.BranchId)
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

func (m *userAdminRepository) Delete(ctx context.Context, id string, deleted_by string) error {
	query := `UPDATE user_admins SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=? WHERE id =?`
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
func (m *userAdminRepository) Update(ctx context.Context, a *models.UserAdmin) error {
	query := `UPDATE user_admins set modified_by=?, modified_date=? , email=? , full_name=? ,branch_id=? WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	res, err := stmt.ExecContext(ctx, a.ModifiedBy, time.Now(), a.Email, a.FullName, a.BranchId, a.Id)
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

// DecodeCursor will decode cursor from user for mysql
func DecodeCursor(encodedTime string) (time.Time, error) {
	byt, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		return time.Time{}, err
	}

	timeString := string(byt)
	t, err := time.Parse(timeFormat, timeString)

	return t, err
}

// EncodeCursor will encode cursor from mysql to user
func EncodeCursor(t time.Time) string {
	timeString := t.Format(timeFormat)

	return base64.StdEncoding.EncodeToString([]byte(timeString))
}
