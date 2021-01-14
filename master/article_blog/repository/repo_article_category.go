package repository

import (
	"context"
	"database/sql"

	"github.com/master/article_category"

	"time"

	"github.com/sirupsen/logrus"

	"github.com/models"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type articleCategoryRepository struct {
	Conn *sql.DB
}

// NewuserRepository will create an object that represent the article.repository interface
func NewArticleCategoryRepository(Conn *sql.DB) article_category.Repository {
	return &articleCategoryRepository{Conn}
}
func (m *articleCategoryRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.ArticleCategory, error) {
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

	result := make([]*models.ArticleCategory, 0)
	for rows.Next() {
		t := new(models.ArticleCategory)
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
			&t.ArticleCategoryName,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (m *articleCategoryRepository) GetByID(ctx context.Context, id string) (res *models.ArticleCategory, err error) {
	query := `SELECT * FROM article_categories WHERE `

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

func (m *articleCategoryRepository) Update(ctx context.Context, ar *models.ArticleCategory) error {
	panic("implement me")
}

func (m *articleCategoryRepository) Delete(ctx context.Context, id string, deleted_by string) error {
	panic("implement me")
}

func (m *articleCategoryRepository) Insert(ctx context.Context, a *models.ArticleCategory) error {
	query := `INSERT article_categories SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , deleted_date=? , is_deleted=? , is_active=? ,
	article_category_name=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.ArticleCategoryName)
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

func (m *articleCategoryRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT count(*) AS count FROM article_categories WHERE is_deleted = 0 and is_active = 1`

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

func (m *articleCategoryRepository) List(ctx context.Context, limit, offset int) ([]*models.ArticleCategory, error) {

	return nil, nil
}
