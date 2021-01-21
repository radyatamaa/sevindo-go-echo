package http

import (
	"context"
	"github.com/booking/transaction"
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// accessibilityHandler  represent the http handler for country
type TransactionHandler struct {
	TransactionUsecase transaction.Usecase
}

// NewaccessibilityHandler will initialize the countrys/ resources endpoint
func NewaTransactionHandler(e *echo.Echo, us transaction.Usecase) {
	handler := &TransactionHandler{
		TransactionUsecase: us,
	}
	e.POST("/transaction/payments", handler.Create)
}

func isRequestValid(m *models.NewCommandAccessibility) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *TransactionHandler) Create(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}



	t := new(models.TransactionIn)
	if err := c.Bind(t); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrBadParamInput)
	}
	//
	//if ok, err := isRequestValid(t); !ok {
	//	return c.JSON(http.StatusBadRequest, err.Error())
	//}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	var promoId *string
	if t.PromoId != "" {
		promoId = &t.PromoId
	} else {
		promoId = nil
	}

	var resortPaymentId *string
	if t.ResortRoomPayment != "" {
		resortPaymentId = &t.ResortRoomPayment
	} else {
		resortPaymentId = nil
	}

	var orderId *string
	if t.OrderId != "" {
		orderId = &t.OrderId
	} else {
		orderId = nil
	}

	var bookingId *string
	if t.BookingId != "" {
		bookingId = &t.BookingId
	} else {
		bookingId = nil
	}

	tr := &models.Transaction{
		BookingType:         t.BookingType,
		BookingId:        bookingId,
		OrderId:             orderId,
		PromoId:             promoId,
		PaymentMethodId:     t.PaymentMethodId,
		ResortRoomPayment: resortPaymentId,
		Status:              t.Status,
		TotalPrice:          t.TotalPrice,
		Currency:            t.Currency,
		ExChangeRates:       &t.ExChangeRates,
		ExChangeCurrency:    &t.ExChangeCurrency,
		Points:              &t.Points,
		OriginalPrice:       t.OriginalPrice,
		ReferralCode:t.ReferralCode,
	}

	res,err:= a.TransactionUsecase.Insert(ctx, tr, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}


func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	logrus.Error(err)
	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrUnAuthorize:
		return http.StatusUnauthorized
	case models.ErrConflict:
		return http.StatusBadRequest
	case models.ErrBadParamInput:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}


