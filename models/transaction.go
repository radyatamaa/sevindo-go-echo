package models

import "time"

type Transaction struct {
	Id                  string     `json:"id" validate:"required"`
	CreatedBy           string     `json:"created_by" validate:"required"`
	CreatedDate         time.Time  `json:"created_date" validate:"required"`
	ModifiedBy          *string    `json:"modified_by"`
	ModifiedDate        *time.Time `json:"modified_date"`
	DeletedBy           *string    `json:"deleted_by"`
	DeletedDate         *time.Time `json:"deleted_date"`
	IsDeleted           int        `json:"is_deleted" validate:"required"`
	IsActive            int        `json:"is_active" validate:"required"`
	BookingType         int        `json:"booking_type"`
	BookingId        *string `json:"booking_id"`
	PromoId             *string    `json:"promo_id"`
	PaymentMethodId     *string     `json:"payment_method_id"`
	ResortRoomPayment *string `json:"resort_room_payment"`
	Status              int        `json:"status"`
	TotalPrice          float64    `json:"total_price"`
	Currency            string     `json:"currency"`
	OrderId             *string    `json:"order_id"`
	VaNumber            *string    `json:"va_number"`
	ExChangeRates 		*float64	`json:"ex_change_rates"`
	ExChangeCurrency 	*string		`json:"ex_change_currency"`
	Points 				*float64	`json:"points"`
	OriginalPrice 		*float64	`json:"original_price"`
	Remarks				*string 	`json:"remarks"`
	TicketPrice 		*float64 `json:"ticket_price"`
	ReferralCode 		*string `json:"referral_code"`
}

type TransactionIn struct {
	BookingType         int     `json:"booking_type,omitempty"`
	BookingId           string  `json:"booking_id"`
	OrderId             string  `json:"order_id"`
	PaypalOrderId       string  `json:"paypal_order_id"`
	CcTokenId           string  `json:"cc_token_id"`
	CcAuthId            string  `json:"cc_auth_id"`
	PromoId             string  `json:"promo_id"`
	PaymentMethodId     *string  `json:"payment_method_id"`
	ResortRoomPayment string `json:"resort_room_payment"`
	Status              int     `json:"status,omitempty"`
	TotalPrice          float64 `json:"total_price,omitempty"`
	Currency            string  `json:"currency"`
	Points              float64 `json:"points"`
	ExChangeRates 		float64	`json:"ex_change_rates"`
	ExChangeCurrency 	string		`json:"ex_change_currency"`
	OriginalPrice 		*float64	`json:"original_price"`
	ReferralCode 		*string `json:"referral_code"`
}