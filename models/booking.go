package models

import "time"

type Booking struct {
	Id                 string     `json:"id" validate:"required"`
	CreatedBy          string     `json:"created_by":"required"`
	CreatedDate        time.Time  `json:"created_date" validate:"required"`
	ModifiedBy         *string    `json:"modified_by"`
	ModifiedDate       *time.Time `json:"modified_date"`
	DeletedBy          *string    `json:"deleted_by"`
	DeletedDate        *time.Time `json:"deleted_date"`
	IsDeleted          int        `json:"is_deleted" validate:"required"`
	IsActive           int        `json:"is_active" validate:"required"`
	OrderId            string     `json:"order_id"`
	GuestDesc          string     `json:"guest_desc"`
	BookedBy           string     `json:"booked_by"`
	BookedByEmail      string     `json:"booked_by_email"`
	BookingDate        time.Time  `json:"booking_date"`
	ExpiredDatePayment *time.Time `json:"expired_date_payment"`
	UserId             *string    `json:"user_id"`
	Status             int        `json:"status"`
	TicketCode         string     `json:"ticket_code"`
	TicketQRCode       string     `json:"ticket_qr_code"`
	PaymentUrl         *string    `json:"payment_url"`
	ResortId 		 	*string `json:"resort_id"`
	ResortRoomId 		*string `json:"resort_room_id"`
	CheckInDate 		*string `json:"check_in_date"`
	CheckOutDate 		*string `json:"check_out_date"`
}

type NewBookingCommand struct {
	Id                string   `json:"id"`
	GuestDesc         []GuestDescObj   `json:"guest_desc"`
	BookedByEmail     string   `json:"booked_by_email"`
	BookingDate       string   `json:"booking_date"`
	UserId            string  `json:"user_id"`
	CheckInDate 		string `json:"check_in_date"`
	CheckOutDate 		string `json:"check_out_date"`
	ResortId 		 	string `json:"resort_id"`
	ResortRoomId 		string `json:"resort_room_id"`
}
type GuestDescObj struct {
	Title       string      `json:"title"`
	FullName    string      `json:"fullname"`
	PhoneNumber string `json:"phone_number"`
}
