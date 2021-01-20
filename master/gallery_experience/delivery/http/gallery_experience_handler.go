package http

import (
	"context"
	"github.com/auth/identityserver"
	"github.com/master/gallery_experience"
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
type galleryexperienceHandler struct {
	isUsecase identityserver.Usecase
	galleyexperienceUsecase gallery_experience.Usecase
}

// NewpromoHandler will initialize the countrys/ resources endpoint
func NewGalleryExperienceHandler(e *echo.Echo, us gallery_experience.Usecase,	isUsecase identityserver.Usecase) {
	handler := &galleryexperienceHandler{
		isUsecase:isUsecase,
		galleyexperienceUsecase: us,
	}
	e.POST("/master/gallery_experience", handler.Create)
	e.PUT("/master/gallery_experience/:id", handler.UpdateGalleryExperience)
	e.DELETE("/master/gallery_experience/:id", handler.Delete)
	// e.GET("/countrys/:id/credit", handler.GetCreditByID)
	e.GET("/master/gallery_experience/:id", handler.GetDetailID)
	e.GET("/master/gallery_experience", handler.List)
}

func isRequestValid(m *models.NewCommandGalleryExperience) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (a *galleryexperienceHandler) List(c echo.Context) error {
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

	result, err := a.galleyexperienceUsecase.List(ctx, page, limit, offset,search)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}


func (a *galleryexperienceHandler) Create(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	filupload, image, _ := c.Request().FormFile("experience_picture")
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
		imagePat, _ := a.isUsecase.UploadFileToBlob(fileLocation, "GalleryExperience")
		imagePath = imagePat
		targetFile.Close()
		errRemove := os.Remove(fileLocation)
		if errRemove != nil {
			return models.ErrInternalServerError
		}
	}

	longitude, _ := strconv.ParseFloat(c.FormValue("longitude"),64)
	latitude, _ := strconv.ParseFloat(c.FormValue("latitude"),64)

	userCommand := models.NewCommandGalleryExperience{
		Id:                c.FormValue("id"),
		ExperienceName:    c.FormValue("experience_name"),
		ExperienceDesc:    c.FormValue("experience_desc"),
		ExperiencePicture: &imagePath,
		Longitude:         longitude,
		Latitude:          latitude,
	}


	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res,err := a.galleyexperienceUsecase.Create(ctx, &userCommand, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}


func (a *galleryexperienceHandler) GetDetailID(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	id := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}


	result, err := a.galleyexperienceUsecase.GetById(ctx, id, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}
func (a *galleryexperienceHandler) UpdateGalleryExperience(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	//id := c.Param("id")
	//==ids ,_:= strconv.Atoi(id)
	filupload, image, _ := c.Request().FormFile("experience_picture")
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
		imagePat, _ := a.isUsecase.UploadFileToBlob(fileLocation, "GalleryExperience")
		imagePath = imagePat
		targetFile.Close()
		errRemove := os.Remove(fileLocation)
		if errRemove != nil {
			return models.ErrInternalServerError
		}
	}

	longitude, _ := strconv.ParseFloat(c.FormValue("longitude"),64)
	latitude, _ := strconv.ParseFloat(c.FormValue("latitude"),64)

	userCommand := models.NewCommandGalleryExperience{
		Id:                c.FormValue("id"),
		ExperienceName:    c.FormValue("experience_name"),
		ExperienceDesc:    c.FormValue("experience_desc"),
		ExperiencePicture: &imagePath,
		Longitude:         longitude,
		Latitude:          latitude,
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = a.galleyexperienceUsecase.Update(ctx, &userCommand, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, userCommand)
}

func (a *galleryexperienceHandler) Delete (c echo.Context) error {
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
	result, err := a.galleyexperienceUsecase.Delete(ctx, id, token)
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
