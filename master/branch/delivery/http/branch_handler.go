package http

import (
	"context"
	"github.com/labstack/echo"
	"github.com/master/branch"
	"github.com/models"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// branchHandler  represent the http handler for country
type branchHandler struct {
	branchUsecase branch.Usecase
}

// NewcountryHandler will initialize the countrys/ resources endpoint
func NewbranchHandler(e *echo.Echo, us branch.Usecase) {
	handler := &branchHandler{
		branchUsecase: us,
	}
	e.POST("/master/branch", handler.Create)
	// e.PUT("/branchys/:id", handler.UpdateBranch)
	// e.DELETE("/countrys/:id", handler.Delete)
	// e.GET("/countrys/:id/credit", handler.GetCreditByID)
	e.GET("/master/branch/:id", handler.GetDetailID)
	// e.GET("/countrys", handler.List)
	//e.POST("/subscribe", handler.Subscribe)
}

func isRequestValid(m *models.NewCommandBranch) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *branchHandler) Create(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")
	var branch models.NewCommandBranch
	err := c.Bind(&branch)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res,err := a.branchUsecase.Create(ctx, &branch, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (a *branchHandler) GetDetailID(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	id := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := a.branchUsecase.GetById(ctx, id, token)
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