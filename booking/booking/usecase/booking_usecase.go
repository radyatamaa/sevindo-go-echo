package usecase

import (
	"context"
	"encoding/json"
	"github.com/auth/identityserver"
	"github.com/auth/user"
	"github.com/booking/booking"
	guuid "github.com/google/uuid"
	"github.com/models"
	"github.com/skip2/go-qrcode"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

type bookingUsecase struct {
	bookingRepo    booking.Repository
	contextTimeout time.Duration
	userUsecase user.Usecase
	isUsecase identityserver.Usecase
}



// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewbookingUsecase(			isUsecase identityserver.Usecase,userUsecase user.Usecase,bookingRepo    booking.Repository, timeout time.Duration) booking.Usecase {
	return &bookingUsecase{
		isUsecase:isUsecase,
		userUsecase:userUsecase,
		bookingRepo:bookingRepo,
		contextTimeout: timeout,
	}
}

func (b bookingUsecase) Insert(c context.Context, booking *models.NewBookingCommand,token string) ([]*models.NewBookingCommand, error, error) {
	ctx, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()

	layoutFormat := "2006-01-02 15:04:05"
	bookingDate, errDate := time.Parse(layoutFormat, booking.BookingDate)
	if errDate != nil {
		return nil, errDate, nil
	}
	orderId, err := generateRandomString(12)
	if err != nil {
		return nil, models.ErrInternalServerError, nil
	}

	// re-generate if duplicate order id
	if b.bookingRepo.CheckBookingCode(ctx, orderId) {
		orderId, err = generateRandomString(12)
		if err != nil {
			return nil, models.ErrInternalServerError, nil
		}
	}

	ticketCode, err := generateRandomString(12)
	if err != nil {
		return nil, models.ErrInternalServerError, nil
	}
	var createdBy string
	if token != "" {
		currentUser, err := b.userUsecase.ValidateTokenUser(ctx, token)
		if err != nil {
			return nil, err, nil
		}
		createdBy = currentUser.UserEmail
	} else {
		createdBy = booking.BookedByEmail
	}

	fileNameQrCode, err := generateQRCode(orderId)
	if err != nil {
		return nil, models.ErrInternalServerError, nil
	}
	imagePath, _ := b.isUsecase.UploadFileToBlob(*fileNameQrCode, "TicketBookingQRCode")

	errRemove := os.Remove(*fileNameQrCode)
	if errRemove != nil {
		return nil, models.ErrInternalServerError, nil
	}

	expIntineraryStartEndPoint, _ := json.Marshal(booking.GuestDesc)
	bookingM := models.Booking{
		Id:                 guuid.New().String(),
		CreatedBy:          createdBy,
		CreatedDate:        time.Now(),
		ModifiedBy:         nil,
		ModifiedDate:       nil,
		DeletedBy:          nil,
		DeletedDate:        nil,
		IsDeleted:          0,
		IsActive:           1,
		OrderId:            orderId,
		GuestDesc:          string(expIntineraryStartEndPoint),
		BookedBy:           createdBy,
		BookedByEmail:      createdBy,
		BookingDate:        bookingDate,
		ExpiredDatePayment: nil,
		UserId:             &booking.UserId,
		Status:             0,
		TicketCode:         ticketCode,
		TicketQRCode:       imagePath,
		PaymentUrl:         nil,
		ResortId:           &booking.ResortId,
		ResortRoomId:       &booking.ResortRoomId,
		CheckInDate:        &booking.CheckInDate,
		CheckOutDate:       &booking.CheckOutDate,
	}
	_,err = b.bookingRepo.Insert(ctx,&bookingM)
	if err != nil {
		return nil, err, nil
	}
	booking.Id = bookingM.Id
	var result []*models.NewBookingCommand
	result = append(result,booking)
	return result, nil, nil
}

func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes, err := generateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}
func generateQRCode(content string) (*string, error) {
	var png []byte
	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}
	name, err := generateRandomString(5)
	if err != nil {
		return nil, err
	}

	fileName := name + ".png"
	err = ioutil.WriteFile(fileName, png, 0700)
	copy, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	copy.Close()
	return &fileName, nil

	//err := qrcode.WriteFile("https://example.org", qrcode.Medium, 256, "qr.png")

}