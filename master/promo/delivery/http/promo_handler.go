package http

import (
	"context"
	"github.com/auth/identityserver"
	"github.com/master/promo"
	"io"
	"os"
	"path/filepath"

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
type promoHandler struct {
	isUsecase identityserver.Usecase
	promoUsecase promo.Usecase
}

// NewpromoHandler will initialize the countrys/ resources endpoint
func NewpromoHandler(e *echo.Echo, us promo.Usecase, isUsecase identityserver.Usecase) {
	handler := &promoHandler{
		isUsecase:isUsecase,
		promoUsecase: us,
	}
	e.POST("/master/promo", handler.Create)
	e.PUT("/master/promo/:id", handler.UpdatePromo)
	e.DELETE("/master/promo/:id", handler.Delete)
	// e.GET("/countrys/:id/credit", handler.GetCreditByID)
	e.GET("/master/promo/:id", handler.GetDetailID)
	e.GET("/master/promo", handler.List)
}

func isRequestValid(m *models.NewCommandPromo) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (a *promoHandler) List(c echo.Context) error {
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

	result, err := a.promoUsecase.List(ctx, page, limit, offset,search)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}


func (a *promoHandler) Create(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	filupload, image, _ := c.Request().FormFile("promo_image")
	dir, err := os.Getwd()
	if err != nil {
		return models.ErrInternalServerError
	}
	var imagePath string
	if filupload != nil {
		fileLocation := filepath.Join(dir, "files", image.Filename)
		targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			os.MkdirAll(filepath.Join(dir, "files"), os.ModePerm)
			return models.ErrInternalServerError
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, filupload); err != nil {
			return models.ErrInternalServerError
		}

		//w.Write([]byte("done"))
		imagePat, _ := a.isUsecase.UploadFileToBlob(fileLocation, "Promo")
		imagePath = imagePat
		targetFile.Close()
		errRemove := os.Remove(fileLocation)
		if errRemove != nil {
			return models.ErrInternalServerError
		}
	}
	var start_date *string
	promo := 	c.FormValue("start_date")
	if promo != ""{
		start_date = &promo
	}
	var end_date *string
	promo1 := c.FormValue("end_date")
	if promo1 != ""{
		end_date = &promo1
	}
	var how_to_get *string
	promo2 := c.FormValue("how_to_get")
	if promo2 != ""{
		how_to_get = &promo2
	}
	var how_to_use *string
	promo3 := c.FormValue("how_to_use")
	if promo3 != ""{
		how_to_use = &promo3
	}
	var term_condition *string
	promo4 := c.FormValue("term_condition")
	if promo4 != ""{
		term_condition = &promo4
	}
	var disclaimer *string
	promo5 := c.FormValue("disclaimer")
	if promo5 != ""{
		disclaimer = &promo5
	}

	promo_value, _ := strconv.ParseFloat(c.FormValue("promo_value"),64)
	max_discount, _ := strconv.ParseFloat(c.FormValue("max_discount"),32)
	promo_type, _ := strconv.Atoi(c.FormValue("promo_type"))
	max_usage, _ := strconv.Atoi(c.FormValue("max_usage"))
	production_capacity, _ := strconv.Atoi(c.FormValue("production_capacity"))
	currency_id, _ := strconv.Atoi(c.FormValue("currency_id"))
	max := float32(max_discount)
	userCommand := models.NewCommandPromo{
		Id:            		c.FormValue("id"),
		PromoCode:   	 	c.FormValue("promo_code"),
		PromoName:    		c.FormValue("promo_name"),
		PromoDesc:       	c.FormValue("promo_desc"),
		PromoValue:       	promo_value,
		PromoType:       	promo_type,
		PromoImage: 		imagePath,
		StartDate:       	start_date,
		EndDate:       		end_date,
		HowToGet:       	how_to_get,
		HowToUse:       	how_to_use,
		TermCondition:     term_condition,
		Disclaimer:       	disclaimer,
		MaxDiscount:        &max,
		MaxUsage:       	&max_usage,
		ProductionCapacity: &production_capacity,
		CurrencyId:       	&currency_id,
	}


	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res,err := a.promoUsecase.Create(ctx, &userCommand, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}


func (a *promoHandler) GetDetailID(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	id := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}


	result, err := a.promoUsecase.GetById(ctx, id, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}
func (a *promoHandler) UpdatePromo(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	id := c.Param("id")
	//ids ,_:= strconv.Atoi(id)
	filupload, image, _ := c.Request().FormFile("promo_image")
	dir, err := os.Getwd()
	if err != nil {
		return models.ErrInternalServerError
	}
	var imagePath string
	if filupload != nil {
		fileLocation := filepath.Join(dir, "files", image.Filename)
		targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			os.MkdirAll(filepath.Join(dir, "files"), os.ModePerm)
			return models.ErrInternalServerError
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, filupload); err != nil {
			return models.ErrInternalServerError
		}

		//w.Write([]byte("done"))
		imagePat, _ := a.isUsecase.UploadFileToBlob(fileLocation, "Promo")
		imagePath = imagePat
		targetFile.Close()
		errRemove := os.Remove(fileLocation)
		if errRemove != nil {
			return models.ErrInternalServerError
		}
	}

	var start_date *string
	promo := 	c.FormValue("start_date")
	if promo != ""{
		start_date = &promo
	}
	var end_date *string
	promo1 := c.FormValue("end_date")
	if promo1 != ""{
		end_date = &promo1
	}
	var how_to_get *string
	promo2 := c.FormValue("how_to_get")
	if promo2 != ""{
		how_to_get = &promo2
	}
	var how_to_use *string
	promo3 := c.FormValue("how_to_use")
	if promo3 != ""{
		how_to_use = &promo3
	}
	var term_condition *string
	promo4 := c.FormValue("term_condition")
	if promo4 != ""{
		term_condition = &promo4
	}
	var disclaimer *string
	promo5 := c.FormValue("disclaimer")
	if promo5 != ""{
		disclaimer = &promo5
	}

	promo_value, _ := strconv.ParseFloat(c.FormValue("promo_value"),64)
	max_discount, _ := strconv.ParseFloat(c.FormValue("max_discount"),32)
	promo_type, _ := strconv.Atoi(c.FormValue("promo_type"))
	max_usage, _ := strconv.Atoi(c.FormValue("max_usage"))
	production_capacity, _ := strconv.Atoi(c.FormValue("production_capacity"))
	currency_id, _ := strconv.Atoi(c.FormValue("currency_id"))
	max := float32(max_discount)
	userCommand := models.NewCommandPromo{
		Id:            id,
		PromoCode:   	 	c.FormValue("promo_code"),
		PromoName:    		c.FormValue("promo_name"),
		PromoDesc:       	c.FormValue("promo_desc"),
		PromoValue:       	promo_value,
		PromoType:       	promo_type,
		PromoImage: 		imagePath,
		StartDate:       	start_date,
		EndDate:       		end_date,
		HowToGet:       	how_to_get,
		HowToUse:       	how_to_use,
		TermCondition:     term_condition,
		Disclaimer:       	disclaimer,
		MaxDiscount:        &max,
		MaxUsage:       	&max_usage,
		ProductionCapacity: &production_capacity,
		CurrencyId:       	&currency_id,
	}


	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = a.promoUsecase.Update(ctx, &userCommand, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, userCommand)
}

func (a *promoHandler) Delete (c echo.Context) error {
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
	//ids ,_:= strconv.Atoi(id)
	result, err := a.promoUsecase.Delete(ctx, id, token)
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
