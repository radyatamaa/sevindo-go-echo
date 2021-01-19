package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/master/promo"
	"github.com/models"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type promoRepository struct {
	Conn *sql.DB
}



// NewuserRepository will create an object that represent the article.repository interface
func NewPromoRepository(Conn *sql.DB) promo.Repository {
	return &promoRepository{Conn}
}
func (m *promoRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Promo, error) {
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

	result := make([]*models.Promo, 0)
	for rows.Next() {
		t := new(models.Promo)
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
			&t.PromoCode,
			&t.PromoName,
			&t.PromoDesc,
			&t.PromoValue,
			&t.PromoType,
			&t.PromoImage,
			&t.StartDate,
			&t.EndDate,
			&t.HowToGet,
			&t.HowToUse,
			&t.TermCondition,
			&t.Disclaimer,
			&t.MaxDiscount,
			&t.MaxUsage,
			&t.ProductionCapacity,
			&t.CurrencyId,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (m *promoRepository) GetByID(ctx context.Context, id string) (res *models.Promo, err error) {
	query := `SELECT * FROM promos WHERE `

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

func (m *promoRepository) Update(ctx context.Context, a *models.Promo) error {
	query := `UPDATE promos set modified_by=?, modified_date=? ,
	promo_code=?, promo_name=?, promo_desc=?,
	promo_value=?, promo_type=?, promo_image=?, 
	start_date=?, end_date=?, how_to_get=?, how_to_use=?,
	term_condition=?, disclaimer=?, max_discount=?, max_usage=?,
	production_capacity=?, currency_id=?
	WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	res, err := stmt.ExecContext(ctx, a.ModifiedBy, time.Now(), a.PromoCode, a.PromoName,
		a.PromoDesc, a.PromoValue, a.PromoType, a.PromoImage, a.StartDate, a.EndDate,
		a.HowToGet, a.HowToUse, a.TermCondition, a.Disclaimer, a.MaxDiscount, a.MaxUsage,
		a.ProductionCapacity, a.CurrencyId,a.Id)
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

func (m *promoRepository) Delete(ctx context.Context, id string, deleted_by string) error {
	query := `UPDATE promos SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=? WHERE id =?`
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


func (m *promoRepository) Insert(ctx context.Context, a *models.Promo) error {
	query := `INSERT promos SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , deleted_date=? , is_deleted=? , is_active=? ,
	promo_code=?, promo_name=?, promo_desc=?,
	promo_value=?, promo_type=?, promo_image=?, 
	start_date=?, end_date=?, how_to_get=?, how_to_use=?,
	term_condition=?, disclaimer=?, max_discount=?, max_usage=?,
	production_capacity=?, currency_id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1,  a.PromoCode, a.PromoName,
		a.PromoDesc, a.PromoValue, a.PromoType, a.PromoImage, a.StartDate, a.EndDate,
		a.HowToGet, a.HowToUse, a.TermCondition, a.Disclaimer, a.MaxDiscount, a.MaxUsage,
		a.ProductionCapacity, a.CurrencyId)
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

func (m *promoRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT count(*) AS count FROM promos WHERE is_deleted = 0 and is_active = 1`

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

func (m *promoRepository) List(ctx context.Context, limit, offset int) ([]*models.Promo, error) {
	query := `SELECT * FROM promos WHERE is_deleted = 0 and is_active = 1 `

	query = query + ` LIMIT ? OFFSET ?`
	list, err := m.fetch(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return list, nil
}
