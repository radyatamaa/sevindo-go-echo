package usecase

import (
	"context"
	"github.com/auth/user"
	"github.com/booking/transaction"
	guuid "github.com/google/uuid"
	"github.com/models"
	"time"
)

type transactionUsecase struct {
	userUsecase user.Usecase
	transactionRepo    transaction.Repository
	contextTimeout time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewtransactionUsecase(		userUsecase user.Usecase, transactionRepo    transaction.Repository ,timeout time.Duration) transaction.Usecase {
	return &transactionUsecase{
		userUsecase:userUsecase,
		transactionRepo:    transactionRepo,
		contextTimeout: timeout,
	}
}


func (p transactionUsecase) Insert(ctx context.Context, transaction *models.Transaction,token string) (*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	currentUser, err := p.userUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return nil, err
	}

	transactionM := models.Transaction{
		Id:                guuid.New().String(),
		CreatedBy:         currentUser.UserEmail,
		CreatedDate:       time.Now(),
		ModifiedBy:        nil,
		ModifiedDate:      nil,
		DeletedBy:         nil,
		DeletedDate:       nil,
		IsDeleted:         0,
		IsActive:          1,
		BookingType:       transaction.BookingType,
		BookingId:         transaction.BookingId,
		PromoId:           transaction.PromoId,
		PaymentMethodId:   transaction.PaymentMethodId,
		ResortRoomPayment: transaction.ResortRoomPayment,
		Status:            0,
		TotalPrice:        transaction.TotalPrice,
		Currency:          transaction.Currency,
		OrderId:           transaction.OrderId,
		VaNumber:          transaction.VaNumber,
		ExChangeRates:     transaction.ExChangeRates,
		ExChangeCurrency:  transaction.ExChangeCurrency,
		Points:            transaction.Points,
		OriginalPrice:     transaction.OriginalPrice,
		Remarks:           transaction.Remarks,
		TicketPrice:       transaction.TicketPrice,
		ReferralCode:      transaction.ReferralCode,
	}

	insert,err := p.transactionRepo.Insert(ctx,&transactionM)
	if err != nil {
		return nil, err
	}

	return insert,nil
}
