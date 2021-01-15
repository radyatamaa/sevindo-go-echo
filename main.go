package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/spf13/viper"

	_articleHttpDeliver "github.com/bxcodec/go-clean-arch/article/delivery/http"
	_articleRepo "github.com/bxcodec/go-clean-arch/article/repository"
	_articleUcase "github.com/bxcodec/go-clean-arch/article/usecase"
	_authorRepo "github.com/bxcodec/go-clean-arch/author/repository"
	_echoMiddleware "github.com/labstack/echo/middleware"

	_isHttpDeliver "github.com/auth/identityserver/delivery/http"
	_isUcase "github.com/auth/identityserver/usecase"

	_userHttpDeliver "github.com/auth/user/delivery/http"
	_userRepo "github.com/auth/user/repository"
	_userUcase "github.com/auth/user/usecase"

	_countryHttpDeliver "github.com/master/country/delivery/http"
	_countryRepo "github.com/master/country/repository"
	_countryUcase "github.com/master/country/usecase"

	_branchHttpDeliver "github.com/master/branch/delivery/http"
	_branchRepo "github.com/master/branch/repository"
	_branchUcase "github.com/master/branch/usecase"

	_currencyHttpDeliver "github.com/master/currency/delivery/http"
	_currencyRepo "github.com/master/currency/repository"
	_currencyUcase "github.com/master/currency/usecase"

	_languageHttpDeliver "github.com/master/language/delivery/http"
	_languageRepo "github.com/master/language/repository"
	_languageUcase "github.com/master/language/usecase"

	_userAdminHttpDeliver "github.com/auth/user_admin/delivery/http"
	_userAdminRepo "github.com/auth/user_admin/repository"
	_userAdminUcase "github.com/auth/user_admin/usecase"

	_provinceHttpDeliver "github.com/master/province/delivery/http"
	_provinceRepo "github.com/master/province/repository"
	_provinceUcase "github.com/master/province/usecase"

	_articlecategoryHttpDeliver "github.com/master/article_category/delivery/http"
	_articlecategoryRepo "github.com/master/article_category/repository"
	_articlecategoryUcase "github.com/master/article_category/usecase"

	_resortHttpDeliver "github.com/services/resort/delivery/http"
	_resortRepo "github.com/services/resort/repository"
	_resortUcase "github.com/services/resort/usecase"

	_cityHttpDeliver "github.com/master/city/delivery/http"
	_cityRepo "github.com/master/city/repository"
	_cityUcase "github.com/master/city/usecase"

	_resortPhotoRepo "github.com/services/resort_photo/repository"

	_roleHttpDeliver "github.com/master/role/delivery/http"
	_roleRepo "github.com/master/role/repository"
	_roleUcase "github.com/master/role/usecase"

	_bankHttpDeliver "github.com/master/bank/delivery/http"
	_bankRepo "github.com/master/bank/repository"
	_bankUcase "github.com/master/bank/usecase"

	_articleblogHttpDeliver "github.com/master/article_blog/delivery/http"
	_articleblogRepo "github.com/master/article_blog/repository"
	_articleblogUcase "github.com/master/article_blog/usecase"

	_districtsHttpDeliver "github.com/master/districts/delivery/http"
	_districtsRepo "github.com/master/districts/repository"
	_districtsUcase "github.com/master/districts/usecase"

	_amenitiesHttpDeliver "github.com/master/amenities/delivery/http"
	_amenitiesRepo "github.com/master/amenities/repository"
	_amenitiesUcase "github.com/master/amenities/usecase"
)

func main() {
	//accountStorage := "cgostorage"
	//accessKeyStorage := "OwvEOlzf6e7QwVoV0H75GuSZHpqHxwhYnYL9QbpVPgBRJn+26K26aRJxtZn7Ip5AhaiIkw9kH11xrZSscavXfQ=="
	//dbHost := viper.GetString(`database.host`)
	//dbPort := viper.GetString(`database.port`)
	//dbUser := viper.GetString(`database.user`)
	//dbPass := viper.GetString(`database.pass`)
	//dbName := viper.GetString(`database.name`)
	//dev db

	//dev env
	//dev db
	dbHost := "bkni-ri.mysql.database.azure.com"
	dbPort := "3306"
	dbUser := "adminbkni@bkni-ri"
	dbPass := "Standar123."
	dbName := "sevindo_dev"
	//dev IS
	baseUrlis := "https://bkni-identity-server-dev.azurewebsites.net"
	//dev URL Forgot Password
	urlForgotPassword := "http://cgo-web-api-dev.azurewebsites.net/account/change-password"
	//google auth
	redirectUrlGoogle := "http://api.dev.cgo.co.id/account/callback"
	clientIDGoogle := "246939853284-f1san8r9bvsc4soef7in80bdti187op5.apps.googleusercontent.com"
	clientSecretGoogle := "TF-b8lBN77fiZLzJ3tuG3FFj"
	tokenSystem := "eyJhbGciOiJSUzI1NiIsImtpZCI6IjJmYzU2MjRkYjQ"
	////dev IS
	//baseUrlis := "http://identity-server-asparnas.azurewebsites.net"
	////dev URL Forgot Password
	//urlForgotPassword := "http://cgo-web-api-dev.azurewebsites.net/account/change-password"
	//basicAuth := "cm9jbGllbnQ6c2VjcmV0"
	//redirectUrlGoogle := "http://cgo-web-api.azurewebsites.net/account/callback"
	//clientIDGoogle := "422978617473-acff67dn9cgbomorrbvhqh2i1b6icm84.apps.googleusercontent.com"
	//clientSecretGoogle := "z_XfHM4DtamjRmJdpu8q0gQf"

	basicAuth := "cm9jbGllbnQ6c2VjcmV0"
	accountStorage := "bkniristorage"
	accessKeyStorage := "/qInIc1r2fMeHHjonpstK8H8HOO5GFIDM4TV/n5+Wk9be3t+UPD4OS0qiVABDIRK5y7XBdlQiHrGyu6M1DDjjQ=="

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)
	if err != nil && viper.GetBool("debug") {
		fmt.Println(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	//middL := middleware.InitMiddleware()
	//e.Use(middL.CORS)
	e.Use(_echoMiddleware.CORSWithConfig(_echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))
	authorRepo := _authorRepo.NewMysqlAuthorRepository(dbConn)
	ar := _articleRepo.NewMysqlArticleRepository(dbConn)
	userRepo := _userRepo.NewuserRepository(dbConn)
	countryRepo := _countryRepo.NewCountryRepository(dbConn)
	branchRepo := _branchRepo.NewBranchRepository(dbConn)
	currencyRepo := _currencyRepo.NewCurrencyRepository(dbConn)
	adminRepo := _userAdminRepo.NewuserAdminRepository(dbConn)
	languageRepo := _languageRepo.NewLanguageRepository(dbConn)
	provinceRepo := _provinceRepo.NewProvinceRepository(dbConn)
	articlecategoryRepo := _articlecategoryRepo.NewArticleCategoryRepository(dbConn)
	articleblogRepo := _articleblogRepo.NewArticleBlogRepository(dbConn)
	cityRepo := _cityRepo.NewCityRepository(dbConn)
	resortRepo := _resortRepo.NewresortRepository(dbConn)
	resortPhotoRepo := _resortPhotoRepo.NewresortPhotoRepository(dbConn)
	roleRepo := _roleRepo.NewRoleRepository(dbConn)
	bankRepo := _bankRepo.NewBankRepository(dbConn)
	districtsRepo := _districtsRepo.NewDistrictsRepository(dbConn)
	amenitiesRepo := _amenitiesRepo.NewAmenitiesRepository(dbConn)
	timeoutContext := 30 * time.Second
	au := _articleUcase.NewArticleUsecase(ar, authorRepo, timeoutContext)

	isUsecase := _isUcase.NewidentityserverUsecase(urlForgotPassword, redirectUrlGoogle, clientIDGoogle, clientSecretGoogle, baseUrlis, basicAuth, accountStorage, accessKeyStorage)
	adminUsecase := _userAdminUcase.NewuserAdminUsecase(tokenSystem, adminRepo, isUsecase, timeoutContext)
	branchUsecase := _branchUcase.NewbranchUsecase(adminUsecase, branchRepo, timeoutContext)
	currencyUsecase := _currencyUcase.NewcurrencyUsecase(adminUsecase, currencyRepo, timeoutContext)
	userUsecase := _userUcase.NewuserUsecase(userRepo, isUsecase, timeoutContext)
	countryUsecase := _countryUcase.NewcountryUsecase(adminUsecase, countryRepo, timeoutContext)
	languageUsecase := _languageUcase.NewlanguageUsecase(adminUsecase, languageRepo, timeoutContext)
	roleUsecase := _roleUcase.NewroleUsecase(adminUsecase, roleRepo, timeoutContext)
	articlecategoryUsecase := _articlecategoryUcase.NewArticleCategoryUsecase(adminUsecase, articlecategoryRepo, timeoutContext)
	articleblogUsecase := _articleblogUcase.NewArticleBlogUsecase(adminUsecase, articleblogRepo, timeoutContext)
	resortUsecase := _resortUcase.NewresortUsecase(resortPhotoRepo, resortRepo, timeoutContext)
	provinceUsecase := _provinceUcase.NewprovinceUsecase(adminUsecase, provinceRepo, timeoutContext)
	cityUsecase := _cityUcase.NewcityUsecase(adminUsecase, cityRepo, timeoutContext)
	bankUsecase := _bankUcase.NewbankUsecase(adminUsecase, bankRepo, timeoutContext)
	districtsUsecase := _districtsUcase.NewdistrictsUsecase(adminUsecase, districtsRepo, timeoutContext)
	amenitiesUsecase := _amenitiesUcase.NewAmenitiesUsecase(adminUsecase, amenitiesRepo, timeoutContext)

	_resortHttpDeliver.NewresortHandler(e, resortUsecase)
	_branchHttpDeliver.NewbranchHandler(e, branchUsecase)
	_currencyHttpDeliver.NewcurrencyHandler(e, currencyUsecase)
	_userAdminHttpDeliver.NewuserAdminHandler(e, adminUsecase, isUsecase)
	_countryHttpDeliver.NewcountryHandler(e, countryUsecase)
	_branchHttpDeliver.NewbranchHandler(e, branchUsecase)
	_userHttpDeliver.NewuserHandler(e, userUsecase, isUsecase)
	_isHttpDeliver.NewisHandler(e, userUsecase, isUsecase, adminUsecase)
	_languageHttpDeliver.NewlanguageHandler(e, languageUsecase)
	_provinceHttpDeliver.NewprovinceHandler(e, provinceUsecase)
	_articlecategoryHttpDeliver.NewArticleCategoryHandler(e, articlecategoryUsecase)
	_articleblogHttpDeliver.NewArticleBlogHandler(e, articleblogUsecase)
	_cityHttpDeliver.NewcityHandler(e, cityUsecase)
	_articleHttpDeliver.NewArticleHandler(e, au)
	_roleHttpDeliver.NewroleHandler(e, roleUsecase)
	_bankHttpDeliver.NewbankHandler(e, bankUsecase)
	_districtsHttpDeliver.NewdistrictsHandler(e, districtsUsecase)
	_amenitiesHttpDeliver.NewAmenitiesHandler(e, amenitiesUsecase)
	log.Fatal(e.Start(":9090"))
}
