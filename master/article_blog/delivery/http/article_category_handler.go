package http

import (
	"context"
	"net/http"

	"github.com/master/article_category"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/models"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// ArticleCategoryHandler  represent the http handler for language
type ArticleCategoryHandler struct {
	ArticleCategoryUsecase article_category.Usecase
}

// NewArticleCategoryHandler will initialize the languages/ resources endpoint
func NewArticleCategoryHandler(e *echo.Echo, us article_category.Usecase) {
	handler := &ArticleCategoryHandler{
		ArticleCategoryUsecase: us,
	}
	e.POST("/master/article_category", handler.Create)
	// e.PUT("/languages/:id", handler.Updatelanguage)
	// e.DELETE("/languages/:id", handler.Delete)
	// e.GET("/languages/:id/credit", handler.GetCreditByID)
	e.GET("/master/article_category/:id", handler.GetDetailID)
	// e.GET("/languages", handler.List)
	//e.POST("/subscribe", handler.Subscribe)
}

func isRequestValid(m *models.NewCommandArticleCategory) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *ArticleCategoryHandler) Create(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")
	var articlecategory models.NewCommandArticleCategory
	err := c.Bind(&articlecategory)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := a.ArticleCategoryUsecase.Create(ctx, &articlecategory, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (a *ArticleCategoryHandler) GetDetailID(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	id := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := a.ArticleCategoryUsecase.GetById(ctx, id, token)
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
