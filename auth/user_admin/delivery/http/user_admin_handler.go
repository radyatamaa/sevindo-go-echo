package http

import (
	"context"
	"github.com/auth/identityserver"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"net/http"

	"github.com/auth/user_admin"
	"github.com/models"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// userHandler  represent the http handler for user
type userAdminHandler struct {
	userAdminUsecase user_admin.Usecase
	isUsecase   identityserver.Usecase
}

// NewuserHandler will initialize the user_admin/ resources endpoint
func NewuserAdminHandler(e *echo.Echo, us user_admin.Usecase, is identityserver.Usecase) {
	handler := &userAdminHandler{
		userAdminUsecase: us,
		isUsecase:   is,
	}
	e.POST("/user_admin", handler.CreateUser)
	e.PUT("/user_admin/:id", handler.UpdateUser)
	e.DELETE("/user_admin/:id", handler.Delete)
	//e.GET("/user_admin/:id", handler.GetDetailID)
	//e.GET("/user_admin", handler.List)
}

func isRequestValid(m *models.NewCommandUser) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (a *userAdminHandler) Delete(c echo.Context) error {
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

	result, err := a.userAdminUsecase.Delete(ctx, id, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}
// Store will store the user by given request body
func (a *userAdminHandler) CreateUser(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")
	//filupload, image, _ := c.Request().FormFile("profile_pict_url")
	//dir, err := os.Getwd()
	//if err != nil {
	//	return models.ErrInternalServerError
	//}
	//var imagePath string
	//if filupload != nil {
	//	fileLocation := filepath.Join(dir, "files", image.Filename)
	//	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	//	if err != nil {
	//		os.MkdirAll(filepath.Join(dir, "files"), os.ModePerm)
	//		return models.ErrInternalServerError
	//	}
	//	defer targetFile.Close()
	//
	//	if _, err := io.Copy(targetFile, filupload); err != nil {
	//		return models.ErrInternalServerError
	//	}
	//
	//	//w.Write([]byte("done"))
	//	imagePat, _ := a.isUsecase.UploadFileToBlob(fileLocation, "UserProfile")
	//	imagePath = imagePat
	//	targetFile.Close()
	//	errRemove := os.Remove(fileLocation)
	//	if errRemove != nil {
	//		return models.ErrInternalServerError
	//	}
	//}

	var branchId *string
	branch := c.FormValue("branch_id")
	if branch != ""{
		branchId = &branch
	}
	userCommand := models.NewCommandUserAdmin{
		Id:       c.FormValue("id"),
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
		FullName: c.FormValue("full_name"),
		BranchId: branchId,
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	user,error := a.userAdminUsecase.Create(ctx, &userCommand, token)

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (a *userAdminHandler) UpdateUser(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	//filupload, image, _ := c.Request().FormFile("profile_pict_url")
	//dir, err := os.Getwd()
	//if err != nil {
	//	return models.ErrInternalServerError
	//}
	//var imagePath string
	//if filupload != nil {
	//	fileLocation := filepath.Join(dir, "files", image.Filename)
	//	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	//	if err != nil {
	//		os.MkdirAll(filepath.Join(dir, "files"), os.ModePerm)
	//		return models.ErrInternalServerError
	//	}
	//	defer targetFile.Close()
	//
	//	if _, err := io.Copy(targetFile, filupload); err != nil {
	//		return models.ErrInternalServerError
	//	}
	//
	//	//w.Write([]byte("done"))
	//	imagePath, _ = a.isUsecase.UploadFileToBlob(fileLocation, "UserProfile")
	//	targetFile.Close()
	//	errRemove := os.Remove(fileLocation)
	//	if errRemove != nil {
	//		return models.ErrInternalServerError
	//	}
	//}
	var branchId *string
	branch := c.FormValue("branch_id")
	if branch != ""{
		branchId = &branch
	}
	userCommand := models.NewCommandUserAdmin{
		Id:       c.FormValue("id"),
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
		FullName: c.FormValue("full_name"),
		BranchId: branchId,
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	error := a.userAdminUsecase.Update(ctx, &userCommand, token)

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, userCommand)
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