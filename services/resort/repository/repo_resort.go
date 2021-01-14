package repository

import (
	"context"
	"database/sql"
	"github.com/models"
	"github.com/services/resort"
	"github.com/sirupsen/logrus"
	"strconv"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type resortRepository struct {
	Conn *sql.DB
}

// NewuserRepository will create an object that represent the article.repository interface
func NewresortRepository(Conn *sql.DB) resort.Repository {
	return &resortRepository{Conn}
}
func (m *resortRepository) fetchJoin(ctx context.Context, query string, args ...interface{}) ([]*models.ResortJoin, error) {
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

func (m resortRepository) GetAll(ctx context.Context, id []string, capacity int,limit ,offset int) ([]*models.ResortJoin, error) {
	query := `SELECT r.*,	
			d.districts_name,
			d.city_id,
			c.city_name,
			c.province_id,
			p.province_name,
			b.branch_name,
			(SELECT rrp.price from resort_rooms rs join resort_room_payments rrp on rrp.resort_room_id = rs.id where rs.resort_id = r.id and rrp.is_active = 1 and rrp.is_deleted = 0 and rs.is_active = 1 and rs.is_deleted = 0 and rs.resort_capacity >= `+strconv.Itoa(capacity)+` ORDER BY rrp.price asc LIMIT 1) as price,
			(SELECT cur.currency_name from resort_rooms rs join resort_room_payments rrp on rrp.resort_room_id = rs.id join currencies cur on rrp.currency = cur.id where rs.resort_id = r.id and rrp.is_active = 1 and rrp.is_deleted = 0 and rs.is_active = 1 and rs.is_deleted = 0 and rs.resort_capacity >= `+strconv.Itoa(capacity)+` ORDER BY rrp.price asc LIMIT 1) as currency
		FROM resorts r
		JOIN branches b on b.id = r.branch_id
		JOIN districts d on d.id = r.districts_id
		JOIN cities c on c.id = d.city_id
		JOIN provinces p on p.id = c.province_id
		JOIN resort_rooms rr on rr.id = (SELECT rs.id from resort_rooms rs join resort_room_payments rrp on rrp.resort_room_id = rs.id where rs.resort_id = r.id and rrp.is_active = 1 and rrp.is_deleted = 0 and rs.is_active = 1 and rs.is_deleted = 0 and rs.resort_capacity >= `+strconv.Itoa(capacity)+` ORDER BY rrp.price asc LIMIT 1)
		WHERE r.is_active = 1 and r.is_deleted = 0  and r.status = 2
 `

	if limit != 0{
		query = query + ` LIMIT `+strconv.Itoa(limit)+` OFFSET `+ strconv.Itoa(offset) + ` `
	}
	list, err := m.fetchJoin(ctx, query)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (m *resortRepository) GetAllCount(ctx context.Context, id []string, capacity int) (int, error) {
	query := `SELECT count(*) AS count
		FROM resorts r
		JOIN branches b on b.id = r.branch_id
		JOIN districts d on d.id = r.districts_id
		JOIN cities c on c.id = d.city_id
		JOIN provinces p on p.id = c.province_id
		WHERE r.is_active = 1 and r.is_deleted = 0 and r.status = 2
 `
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
