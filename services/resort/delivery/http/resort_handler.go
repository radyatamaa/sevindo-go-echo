package http

import (
	"context"
	"github.com/services/resort"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/models"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// provinceHandler  represent the http handler for province
type resortHandler struct {
	resortUsecase resort.Usecase
}

// NewprovinceHandler will initialize the provinces/ resources endpoint
func NewresortHandler(e *echo.Echo, us resort.Usecase) {
	handler := &resortHandler{
		resortUsecase: us,
	}
	e.GET("/service/resort-search", handler.GetAll)
}

func isRequestValid(m *models.NewCommandProvince) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *resortHandler) GetAll(c echo.Context) error {

	qpage := c.QueryParam("page")
	qperPage := c.QueryParam("size")
	startDate := c.QueryParam("start_date")
	endDate := c.QueryParam("end_date")
	capacity := c.QueryParam("capacity")

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
	capc ,_:= strconv.Atoi(capacity)
	result, err := a.resortUsecase.GetAll(ctx, page, limit, offset,capc,startDate,endDate)
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
