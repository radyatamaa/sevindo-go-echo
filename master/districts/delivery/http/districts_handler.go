package http

import (
	"context"
	"github.com/master/districts"

	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/models"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// countryHandler  represent the http handler for country
type districtsHandler struct {
	districtsUsecase districts.Usecase
}

// NewcityHandler will initialize the countrys/ resources endpoint
func NewdistrictsHandler(e *echo.Echo, us districts.Usecase) {
	handler := &districtsHandler{
		districtsUsecase: us,
	}
	e.POST("/master/districts", handler.Create)
	e.PUT("/master/districts/:id", handler.UpdateDistricts)
	e.DELETE("/master/districts/:id", handler.Delete)
	// e.GET("/countrys/:id/credit", handler.GetCreditByID)
	e.GET("/master/districts/:id", handler.GetDetailID)
	e.GET("/master/districts", handler.List)
}

func isRequestValid(m *models.NewCommandDistricts) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (a *districtsHandler) List(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	qpage := c.QueryParam("page")
	qperPage := c.QueryParam("size")
	search := c.QueryParam("search")

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

	result, err := a.districtsUsecase.List(ctx, page, limit, offset,search)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}


func (a *districtsHandler) Create(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}


	var districts models.NewCommandDistricts
	err := c.Bind(&districts)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res,err := a.districtsUsecase.Create(ctx, &districts, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}


func (a *districtsHandler) GetDetailID(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	id := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	ids ,_:= strconv.Atoi(id)
	result, err := a.districtsUsecase.GetById(ctx, ids, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}
func (a *districtsHandler) UpdateDistricts(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	id := c.Param("id")
	var districts models.NewCommandDistricts
	ids ,_:= strconv.Atoi(id)
	districts.Id = ids
	err := c.Bind(&districts)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = a.districtsUsecase.Update(ctx, &districts, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, districts)
}

func (a *districtsHandler) Delete (c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	id := c.Param("id")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	ids ,_:= strconv.Atoi(id)
	result, err := a.districtsUsecase.Delete(ctx, ids, token)
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
