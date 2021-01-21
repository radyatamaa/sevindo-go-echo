package repository

import (
	"context"
	"database/sql"
	"github.com/booking/transaction"
	"github.com/models"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type transactionRepository struct {
	Conn *sql.DB
}



// NewuserRepository will create an object that represent the article.repository interface
func NewTransactionRepository(Conn *sql.DB) transaction.Repository {
	return &transactionRepository{Conn}
}

func (b transactionRepository) Insert(ctx context.Context, a *models.Transaction) (*models.Transaction, error) {
	query := `INSERT transactions SET id=?,created_by=?,created_date=?,modified_by=?,modified_date=?,deleted_by=?,
				deleted_date=?,is_deleted=?,is_active=?,booking_type=?,booking_id=?,promo_id=?,payment_method_id=?,
				resort_room_payment=?,status=?,total_price=?,currency=?,order_id=?,va_number=?,
				ex_change_rates=?,ex_change_currency=?,points=?,original_price=?,remarks=?,ticket_price=?,
				referral_code=?`

	stmt, err := b.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1,a.BookingType,a.BookingId,
		a.PromoId,a.PaymentMethodId,a.ResortRoomPayment,a.Status,a.TotalPrice,a.Currency,a.OrderId,a.VaNumber,
		a.ExChangeRates,a.ExChangeCurrency,a.Points,a.OriginalPrice,a.Remarks,a.TicketPrice,a.ReferralCode)
	if err != nil {
		return nil, err
	}

	return a, nil
}