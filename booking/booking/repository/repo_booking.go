package repository

import (
	"context"
	"database/sql"
	"github.com/booking/booking"
	"github.com/models"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type bookingRepository struct {
	Conn *sql.DB
}



// NewMysqlArticleRepository will create an object that represent the article.repository interface
func NewbookingRepository(Conn *sql.DB) booking.Repository {
	return &bookingRepository{Conn}
}

func (b bookingRepository) Insert(ctx context.Context, a *models.Booking) (*models.Booking, error) {
	query := `INSERT bookings SET id=?,created_by=?,created_date=?,modified_by=?,modified_date=?,deleted_by=?,
				deleted_date=?,is_deleted=?,is_active=?,order_id=?,guest_desc=?,booked_by=?,booked_by_email=?,
				booking_date=?,expired_date_payment=?,user_id=?,status=?,ticket_code=?,ticket_qr_code=?,
				payment_url=?,resort_id=?,resort_room_id=?,check_in_date=?,check_out_date=?`

	stmt, err := b.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1,a.OrderId,a.GuestDesc,
		a.BookedBy,a.BookedByEmail,a.BookingDate,a.ExpiredDatePayment,a.UserId,a.Status,a.TicketCode,a.TicketQRCode,
		a.PaymentUrl,a.ResortId,a.ResortRoomId,a.CheckInDate,a.CheckOutDate)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (b bookingRepository) CheckBookingCode(ctx context.Context, bookingCode string) bool {
	var code string
	query := `SELECT order_id as code FROM booking_exps WHERE order_id = ?`

	_ = b.Conn.QueryRowContext(ctx, query, bookingCode).Scan(&code)

	if bookingCode == code {
		return true
	}

	return false
}