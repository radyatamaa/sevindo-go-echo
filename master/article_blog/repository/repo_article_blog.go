package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/master/article_blog"
	"strconv"

	"time"

	"github.com/sirupsen/logrus"

	"github.com/models"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type articleblogRepository struct {
	Conn *sql.DB
}

// NewuserRepository will create an object that represent the article.repository interface
func NewArticleBlogRepository(Conn *sql.DB) article_blog.Repository {
	return &articleblogRepository{Conn}
}
func (m *articleblogRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.ArticleBlog, error) {
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

	result := make([]*models.ArticleBlog, 0)
	for rows.Next() {
		t := new(models.ArticleBlog)
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
			&t.ArticleBlogName,
			&t.Title,
			&t.Description,
			&t.CategoryId,
			&t.ArticlePicture,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (m *articleblogRepository) GetByID(ctx context.Context, id int) (res *models.ArticleBlog, err error) {
	query := `SELECT * FROM article_blogs WHERE `

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

func (m *articleblogRepository) Update(ctx context.Context, ar *models.ArticleBlog) error {
	query := `UPDATE article_blogs set modified_by=?, modified_date=? , article_blog_name=? , title=?, description=?, category_id=?, article_picture=?  WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	res, err := stmt.ExecContext(ctx, ar.ModifiedBy, time.Now(), ar.ArticleBlogName, ar.Id, ar.Title, ar.Description, ar.CategoryId)
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

func (m *articleblogRepository) Delete(ctx context.Context, id int, deleted_by string) error {
	query := `UPDATE article_blogs SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=? WHERE id =?`
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

func (m *articleblogRepository) Insert(ctx context.Context, a *models.ArticleBlog) error {
	query := `INSERT article_blogs SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , deleted_date=? , is_deleted=? , is_active=? ,
	article_blog_name=?, title=?, description=?, category_id=?, article_picture=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.ArticleBlogName, a.Title,a.Description, a.CategoryId ,a.ArticlePicture)
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

func (m *articleblogRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT count(*) AS count FROM article_blogs WHERE is_deleted = 0 and is_active = 1`

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

func (m *articleblogRepository) List(ctx context.Context, limit, offset int) ([]*models.ArticleBlog, error) {

	query := `SELECT * FROM article_blogs WHERE is_deleted = 0 and is_active = 1 `

	query = query + ` LIMIT ? OFFSET ?`
	list, err := m.fetch(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return list, nil
}
