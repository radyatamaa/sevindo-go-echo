package http

import (
	"context"
	"github.com/auth/user_admin"
	"net/http"
	"strings"

	"github.com/auth/identityserver"
	"github.com/auth/user"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"

	"github.com/models"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// isHandler  represent the httphandler for is
type isHandler struct {
	isUsecase       identityserver.Usecase
	userUsecase     user.Usecase
	userAdminUsecase user_admin.Usecase
}

// NewisHandler will initialize the iss/ resources endpoint
func NewisHandler(e *echo.Echo,  u user.Usecase, is identityserver.Usecase,userAdminUsecase user_admin.Usecase) {
	handler := &isHandler{
		userAdminUsecase:userAdminUsecase,
		userUsecase:     u,
		isUsecase:       is,
	}
	e.GET("/account/info", handler.GetInfo)
	e.POST("/account/login", handler.Login)
	e.POST("/account/refresh-token", handler.RefreshToken)
	e.POST("/account/request-otp", handler.RequestOTP)
	e.POST("/account/request-otp-tmp", handler.RequestOTPTmp)
	e.GET("/account/verified-email", handler.VerifiedEmail)
	e.POST("/account/callback", handler.CallBack)
	e.GET("/account/forgot-password", handler.ForgotPassword)
	e.GET("/account/change-password", handler.ChangePassword)
	e.GET("/account/check-email-phone-number", handler.CheckEmailPhoneNumber)
}
func (a *isHandler) CheckEmailPhoneNumber(c echo.Context) error {

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	email := c.QueryParam("email")
	phoneNumber := c.QueryParam("phone_number")
	if phoneNumber != "" {
		phoneNumber = strings.ReplaceAll(phoneNumber, " ", "+")

	}
	responseOTP, err := a.userUsecase.CheckEmailORPhoneNumber(ctx, email, phoneNumber, "")
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, responseOTP)
}
func (a *isHandler) ChangePassword(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	email := c.QueryParam("email")
	phoneNumber := c.QueryParam("phone_number")
	if phoneNumber != "" {
		phoneNumber = strings.ReplaceAll(phoneNumber, " ", "+")

	}
	password := c.QueryParam("password")
	responseOTP, err := a.userUsecase.ChangePassword(ctx, token, password, email, phoneNumber)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, responseOTP)
}
func (a *isHandler) ForgotPassword(c echo.Context) error {

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	email := c.QueryParam("email")
	getUser, err := a.userUsecase.GetUserByEmail(ctx, email, "email")
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrNotFound)
	}
	getUserDetail, err := a.isUsecase.GetDetailUserById(getUser.Id, "", "true")
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrNotFound)
	}
	login, err := a.isUsecase.GetToken(email, getUserDetail.Password, "", "1", "email")
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrNotFound)
	}
	responseOTP, err := a.isUsecase.ForgotPassword(email, login.RefreshToken)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, responseOTP)
}
func (a *isHandler) RequestOTP(c echo.Context) error {
	var requestOTP models.RequestOTPNumber
	err := c.Bind(&requestOTP)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	requestOTP.PhoneNumber = c.Request().Form.Get("phone_number")
	responseOTP, err := a.userUsecase.RequestOTP(ctx, requestOTP.PhoneNumber)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, responseOTP)
}
func (a *isHandler) RequestOTPTmp(c echo.Context) error {
	var requestOTP models.RequestOTPTmpNumber
	err := c.Bind(&requestOTP)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	requestOTP.PhoneNumber = c.Request().Form.Get("phone_number")
	requestOTP.Email = c.Request().Form.Get("email")
	//if requestOTP.Email != "" {
	//	checkEmail, _ := a.userUsecase.GetUserByEmail(ctx, requestOTP.Email,"email")
	//	if checkEmail != nil {
	//		return c.JSON(getStatusCode(models.ErrConflict), ResponseError{Message: models.ErrConflict.Error()})
	//	}
	//}
	responseOTP, err := a.isUsecase.RequestOTPTmp(requestOTP.PhoneNumber, requestOTP.Email)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	a.userUsecase.CheckEmailORPhoneNumber(ctx, requestOTP.Email, requestOTP.PhoneNumber, responseOTP.OTP)
	return c.JSON(http.StatusOK, responseOTP)
}
func (a *isHandler) Login(c echo.Context) error {
	var isLogin models.Login
	err := c.Bind(&isLogin)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	isLogin.Email = c.Request().Form.Get("email")
	isLogin.Password = c.Request().Form.Get("password")
	isLogin.Type = c.Request().Form.Get("type")
	isLogin.Scope = c.Request().Form.Get("scope")
	isLogin.XMode = c.Request().Form.Get("x_mode")
	var responseToken *models.GetToken
	if isLogin.Type == "user" {
		isLogin.LoginType = c.Request().Form.Get("login_type")

		token, err := a.userUsecase.Login(ctx, &isLogin)

		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		responseToken = token
	}else if isLogin.Type == "admin_branch" {
		isLogin.LoginType = c.Request().Form.Get("login_type")

		token, err := a.userAdminUsecase.Login(ctx, &isLogin)

		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		responseToken = token
	} else {
		return c.JSON(http.StatusBadRequest, "Bad Request")
	}

	return c.JSON(http.StatusOK, responseToken)
}
func (a *isHandler) RefreshToken(c echo.Context) error {
	var isLogin models.RefreshTokenLogin
	err := c.Bind(&isLogin)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	isLogin.RefreshToken = c.Request().Form.Get("refresh_token")
	responseToken, err := a.isUsecase.RefreshToken(isLogin.RefreshToken)
	return c.JSON(http.StatusOK, responseToken)
}
func (a *isHandler) VerifiedEmail(c echo.Context) error {
	//c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	//c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	//token := c.Request().Header.Get("Authorization")
	email := c.QueryParam("email")
	phoneNumber := c.QueryParam("phone_number")
	if phoneNumber != "" {
		phoneNumber = strings.ReplaceAll(phoneNumber, " ", "+")

	}
	otpCode := c.QueryParam("otp")
	verified := models.VerifiedEmail{
		Email:       email,
		PhoneNumber: phoneNumber,
		CodeOTP:     otpCode,
	}
	//if typeUser == "user" {
	response, err := a.isUsecase.VerifiedEmail(&verified)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, response)
	//} else if typeUser == "merchant" {
	//	return c.JSON(http.StatusNotFound, "Not Implemented")
	//} else {
	//	return c.JSON(http.StatusBadRequest, "Bad Request")
	//}

	return c.JSON(http.StatusBadRequest, "Bad Request")
}
func (a *isHandler) CallBack(c echo.Context) error {
	loginType := c.FormValue("login_type")
	email := c.FormValue("email")
	profilePict := c.FormValue("profile_pict")
	name := c.FormValue("name")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	token, err := a.userUsecase.LoginByGoogle(ctx, loginType, email, profilePict, name)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, token)
}

func (a *isHandler) GetInfo(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	token := c.Request().Header.Get("Authorization")
	typeUser := c.QueryParam("type")
	orderId := c.QueryParam("order_id")

	if typeUser == "user" {
		response, err := a.userUsecase.GetUserInfo(ctx, token, orderId)

		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, response)
	} else if typeUser == "admin_branch" {
		response, err := a.userAdminUsecase.GetUserInfo(ctx, token)

		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, response)
	}else {
		return c.JSON(http.StatusBadRequest, "Bad Request")
	}

	return c.JSON(http.StatusBadRequest, "Bad Request")
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
	case models.ErrConflict:
		return http.StatusConflict
	case models.ErrUnAuthorize:
		return http.StatusUnauthorized
	case models.ErrBadParamInput:
		return http.StatusBadRequest
	case models.ErrUsernamePassword:
		return http.StatusUnauthorized
	case models.ErrInvalidOTP:
		return http.StatusUnauthorized
	case models.ErrNotYetRegister:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
