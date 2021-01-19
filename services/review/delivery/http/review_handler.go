package http

import (
	"context"
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/services/review"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// provinceHandler  represent the http handler for province
type reviewHandler struct {
	reviewUsecase review.Usecase
}

// NewprovinceHandler will initialize the provinces/ resources endpoint
func NewreviewHandler(e *echo.Echo, us review.Usecase) {
	handler := &reviewHandler{
		reviewUsecase: us,
	}
	e.GET("/service/review", handler.GetReviewByResortId)
}

func isRequestValid(m *models.NewCommandProvince) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *reviewHandler) GetReviewByResortId(c echo.Context) error {

	resortId := c.QueryParam("resort_id")
	qpage := c.QueryParam("page")
	qperPage := c.QueryParam("size")

	var limit = 20
	var page = 1
	var offset = 0

	page, _ = strconv.Atoi(qpage)
	limit, _ = strconv.Atoi(qperPage)
	offset = (page - 1) * limit

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := a.reviewUsecase.GetAll(ctx,page,limit,offset,resortId)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
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


