package usecase

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"



	"github.com/auth/identityserver"
	"github.com/auth/user"
	"github.com/models"
	"golang.org/x/net/context"
)

type userUsecase struct {
	userRepo         user.Repository
	identityServerUc identityserver.Usecase
	contextTimeout   time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewuserUsecase( a user.Repository, is identityserver.Usecase,  timeout time.Duration) user.Usecase {
	return &userUsecase{
		userRepo:         a,
		identityServerUc: is,
		contextTimeout:   timeout,
	}
}

func (m userUsecase) GetUserByReferralCode(ctx context.Context, code string) (*models.UserDto, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	getUser, _ := m.userRepo.GetByID(ctx, "", code)
	var result *models.UserDto
	if getUser != nil {
		result = &models.UserDto{
			Id:             getUser.Id,
			CreatedDate:    getUser.CreatedDate,
			UpdatedDate:    getUser.ModifiedDate,
			IsActive:       getUser.IsActive,
			UserEmail:      getUser.UserEmail,
			Password:       "",
			FirstName:       getUser.FirstName,
			LastName:       getUser.LastName,
			PhoneNumber:    getUser.PhoneNumber,
			ProfilePictUrl: getUser.ProfilePictUrl,
			Address:        getUser.Address,
			Dob:            getUser.Dob,
			Gender:         getUser.Gender,
			IdType:         getUser.IdType,
			IdNumber:       getUser.IdNumber,
			ReferralCode:   getUser.ReferralCode,
			Points:         getUser.Points,
		}
	}

	return result, nil
}
func (m userUsecase) CheckEmailORPhoneNumber(ctx context.Context, email string, phoneNumber string, otp string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()
	var user *models.User
	if email != "" {
		user, _ = m.userRepo.GetByUserEmail(ctx, email, "email", "")
		if user == nil {
			return nil, models.ErrNotFound
		}
	} else if phoneNumber != "" {
		user, _ = m.userRepo.GetByUserEmail(ctx, "", "phone_number", phoneNumber)
		if user == nil {
			return nil, models.ErrNotFound
		}
	}
	if otp != "" {
		user.VerificationCode = otp

		err := m.userRepo.Update(ctx, user)
		if err != nil {
			return nil, err
		}
	}

	result := models.ResponseDelete{
		Id:      user.Id,
		Message: user.VerificationCode,
	}

	return &result, nil

}

func (m userUsecase) ChangePassword(c context.Context, token string, password string, email string, phoneNumber string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	//currentUsers, err := m.ValidateTokenUser(ctx, token)
	getUser := models.User{}
	username := ""
	if email != "" {
		user, _ := m.userRepo.GetByUserEmail(ctx, email, "email", "")
		if user == nil {
			return nil, models.ErrNotFound
		}
		username = user.UserEmail
		getUser = *user
	} else if phoneNumber != "" {
		user, _ := m.userRepo.GetByUserEmail(ctx, "", "phone_number", phoneNumber)
		if user == nil {
			return nil, models.ErrNotFound
		}
		username = user.PhoneNumber
		getUser = *user
	}

	updateUser := models.RegisterAndUpdateUser{
		Id:            getUser.Id,
		Username:      username,
		Password:      password,
		Name:          getUser.FirstName,
		GivenName:     "",
		FamilyName:    "",
		Email:         getUser.UserEmail,
		EmailVerified: true,
		Website:       "",
		Address:       "",
		OTP:           "",
		UserType:      1,
		PhoneNumber:   getUser.PhoneNumber,
		LoginType:     *getUser.LoginType,
		UserRoles:     nil,
	}

	_, err := m.identityServerUc.UpdateUser(&updateUser)
	if err != nil {
		return nil, err
	}

	result := models.ResponseDelete{
		Id:      getUser.Id,
		Message: "Success Update Password",
	}
	return &result, nil
}

func (m userUsecase) GetUserByEmail(ctx context.Context, email string, loginType string) (*models.UserDto, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()
	var getUser *models.User
	if loginType == "phone_number" {
		getUser, _ = m.userRepo.GetByUserEmail(ctx, "", loginType, email)
	} else {
		getUser, _ = m.userRepo.GetByUserEmail(ctx, email, loginType, "")
	}

	var result *models.UserDto
	if getUser != nil {
		result = &models.UserDto{
			Id:             getUser.Id,
			CreatedDate:    getUser.CreatedDate,
			UpdatedDate:    getUser.ModifiedDate,
			IsActive:       getUser.IsActive,
			UserEmail:      getUser.UserEmail,
			Password:       "",
			FirstName:       getUser.FirstName,
			LastName:       getUser.LastName,
			PhoneNumber:    getUser.PhoneNumber,
			ProfilePictUrl: getUser.ProfilePictUrl,
			Address:        getUser.Address,
			Dob:            getUser.Dob,
			Gender:         getUser.Gender,
			IdType:         getUser.IdType,
			IdNumber:       getUser.IdNumber,
			ReferralCode:   getUser.ReferralCode,
			Points:         getUser.Points,
		}
	}

	return result, nil
}
func (m userUsecase) LoginByGoogle(c context.Context, loginType string, email string, profilePicture string, name string) (*models.GetToken, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()
	var requestToken *models.GetToken

	checkUser, err := m.userRepo.GetByUserEmail(ctx, email, loginType, "")
	if err != nil {
		return nil, err
	}
	if checkUser == nil {
		password, _ := generateRandomString(11)
		registerUser := models.RegisterAndUpdateUser{
			Id:            "",
			Username:      email,
			Password:      password,
			Name:          name,
			GivenName:     "",
			FamilyName:    "",
			Email:         email,
			EmailVerified: false,
			Website:       "",
			Address:       "",
			OTP:           "",
			UserType:      1,
			PhoneNumber:   "",
			LoginType:     loginType,
			UserRoles:     nil,
		}
		isUser, _ := m.identityServerUc.CreateUser(&registerUser)
		referralCode, _ := generateRandomString(9)
		userModel := models.User{}
		userModel.Id = isUser.Id
		userModel.CreatedBy = email
		userModel.UserEmail = email
		userModel.FirstName = name
		userModel.PhoneNumber = ""
		userModel.VerificationSendDate = time.Now()
		userModel.VerificationCode = isUser.OTP
		userModel.ProfilePictUrl = profilePicture
		userModel.Address = ""
		userModel.Dob = time.Time{}
		userModel.Gender = 0
		userModel.IdType = 0
		userModel.IdNumber = ""
		userModel.ReferralCode = referralCode
		userModel.Points = 0
		userModel.LoginType = &loginType
		err := m.userRepo.Insert(ctx, &userModel)
		if err != nil {
			return nil, err
		}
		getToken, err := m.identityServerUc.GetToken(registerUser.Email, registerUser.Password, "", "1", loginType)
		if err != nil {
			return nil, err
		}
		requestToken = getToken
	} else {
		getDetailUser, err := m.identityServerUc.GetDetailUserById(checkUser.Id, "", "true")
		if err != nil {
			return nil, err
		}
		getToken, err := m.identityServerUc.GetToken(getDetailUser.Email, getDetailUser.Password, "", "1", loginType)
		if err != nil {
			return nil, err
		}
		requestToken = getToken
	}
	return requestToken, nil
}
func (m userUsecase) Delete(c context.Context, userId string, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()
	//currentUserAdmin, err := m.adminUsecase.ValidateTokenAdmin(ctx, token)
	//if err != nil {
	//	return nil, err
	//}
	error := m.userRepo.Delete(ctx, userId, "admin")
	_ = m.identityServerUc.DeleteUser(userId)
	if error != nil {
		response := models.ResponseDelete{
			Id:      userId,
			Message: error.Error(),
		}
		return &response, nil
	}
	response := models.ResponseDelete{
		Id:      userId,
		Message: "Deleted Success",
	}

	return &response, nil
}
func (m userUsecase) RequestOTP(ctx context.Context, phoneNumber string) (*models.RequestOTP, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	requestOTP, err := m.identityServerUc.RequestOTP(phoneNumber)
	if err != nil {
		return nil, err
	}

	getUserByPhoneNumber, err := m.userRepo.GetByUserNumberOTP(ctx, phoneNumber, "")
	if getUserByPhoneNumber != nil {
		user := getUserByPhoneNumber
		user.VerificationCode = requestOTP.OTP

		err = m.userRepo.Update(ctx, user)
		if err != nil {
			return nil, err
		}
	}

	return requestOTP, nil
}
func (m userUsecase) List(ctx context.Context, page, limit, offset int, search string) (*models.UserWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.userRepo.List(ctx, limit, offset, search)
	if err != nil {
		return nil, err
	}

	users := make([]*models.UserInfoDto, len(list))
	for i, item := range list {
		users[i] = &models.UserInfoDto{
			Id:             item.Id,
			CreatedDate:    item.CreatedDate,
			UpdatedDate:    item.ModifiedDate,
			IsActive:       item.IsActive,
			UserEmail:      item.UserEmail,
			FirstName:       item.FirstName,
			LastName:       item.LastName,
			PhoneNumber:    item.PhoneNumber,
			ProfilePictUrl: item.ProfilePictUrl,
			Address:        item.Address,
			Dob:            item.Dob,
			Gender:         item.Gender,
			IdType:         item.IdType,
			IdNumber:       item.IdNumber,
			ReferralCode:   item.ReferralCode,
			Points:         item.Points,
		}
	}
	totalRecords, _ := m.userRepo.Count(ctx)
	totalPage := int(math.Ceil(float64(totalRecords) / float64(limit)))
	prev := page
	next := page
	if page != 1 {
		prev = page - 1
	}

	if page != totalPage {
		next = page + 1
	}
	meta := &models.MetaPagination{
		Page:          page,
		Total:         totalPage,
		TotalRecords:  totalRecords,
		Prev:          prev,
		Next:          next,
		RecordPerPage: len(list),
	}

	response := &models.UserWithPagination{
		Data: users,
		Meta: meta,
	}

	return response, nil
}

func (m userUsecase) Login(ctx context.Context, ar *models.Login) (*models.GetToken, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()
	var requestToken *models.GetToken

	if ar.Scope == "phone_number" {
		checkPhoneNumberExists, _ := m.userRepo.GetByUserNumberOTP(ctx, ar.Email, "")
		if checkPhoneNumberExists == nil {
			return nil, models.ErrNotYetRegister
		}
		existeduser, _ := m.userRepo.GetByUserNumberOTP(ctx, ar.Email, ar.Password)
		if existeduser == nil {
			return nil, models.ErrInvalidOTP
		}
		getToken, err := m.identityServerUc.GetToken(ar.Email, ar.Password, ar.Scope, "1", ar.Scope)
		if err != nil {
			return nil, err
		}

		requestToken = getToken
	} else {
		getToken, err := m.identityServerUc.GetToken(ar.Email, ar.Password, ar.Scope, "1", ar.LoginType)
		if err != nil {
			return nil, err
		}
		if ar.LoginType == "phone_number" {
			existeduser, _ := m.userRepo.GetByUserEmail(ctx, "", ar.LoginType, ar.Email)
			if existeduser == nil {
				return nil, models.ErrUsernamePassword
			}
		} else {
			existeduser, _ := m.userRepo.GetByUserEmail(ctx, ar.Email, ar.LoginType, "")
			if existeduser == nil {
				return nil, models.ErrUsernamePassword
			}
		}
		requestToken = getToken
	}
	return requestToken, nil
}

func (m userUsecase) ValidateTokenUser(ctx context.Context, token string) (*models.UserInfoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	getInfoToIs, err := m.identityServerUc.GetUserInfo(token)
	if err != nil {
		return nil, err
	}
	var existeduser *models.User
	if getInfoToIs.LoginType == "phone_number" {
		existeduser, _ = m.userRepo.GetByUserEmail(ctx, "", getInfoToIs.LoginType, getInfoToIs.PhoneNumber)
		if existeduser == nil {
			return nil, models.ErrUnAuthorize
		}
	} else {
		existeduser, _ = m.userRepo.GetByUserEmail(ctx, getInfoToIs.Email, getInfoToIs.LoginType, "")
		if existeduser == nil {
			return nil, models.ErrUnAuthorize
		}
	}

	userInfo := &models.UserInfoDto{
		Id:             existeduser.Id,
		UserEmail:      existeduser.UserEmail,
		FirstName:       existeduser.FirstName,
		LastName:       existeduser.LastName,
		PhoneNumber:    existeduser.PhoneNumber,
		ProfilePictUrl: existeduser.ProfilePictUrl,
		Points:         existeduser.Points,
		LoginType:      *existeduser.LoginType,
		ReferralCode:   existeduser.ReferralCode,
	}

	return userInfo, nil
}

func (m userUsecase) VerifiedEmail(ctx context.Context, token string, codeOTP string) (*models.UserInfoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	getInfoToIs, err := m.identityServerUc.GetUserInfo(token)
	if err != nil {
		return nil, err
	}
	var existeduser *models.User
	if getInfoToIs.LoginType == "phone_number" {
		existeduser, _ = m.userRepo.GetByUserEmail(ctx, "", getInfoToIs.LoginType, getInfoToIs.PhoneNumber)
		if existeduser == nil {
			return nil, models.ErrUnAuthorize
		}
	} else {
		existeduser, _ = m.userRepo.GetByUserEmail(ctx, getInfoToIs.Email, getInfoToIs.LoginType, "")
		if existeduser == nil {
			return nil, models.ErrUnAuthorize
		}
	}
	verifiedEmail := models.VerifiedEmail{
		Email:   existeduser.UserEmail,
		CodeOTP: codeOTP,
	}
	_, error := m.identityServerUc.VerifiedEmail(&verifiedEmail)
	if error != nil {
		return nil, error
	}
	userInfo := &models.UserInfoDto{
		Id:             existeduser.Id,
		UserEmail:      existeduser.UserEmail,
		FirstName:       existeduser.FirstName,
		LastName:       existeduser.LastName,
		PhoneNumber:    existeduser.PhoneNumber,
		ProfilePictUrl: existeduser.ProfilePictUrl,
		ReferralCode:   existeduser.ReferralCode,
	}

	return userInfo, nil
}
func (m userUsecase) GetUserInfo(ctx context.Context, token string, orderId string) (*models.UserInfoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	getInfoToIs, err := m.identityServerUc.GetUserInfo(token)
	if err != nil {
		return nil, err
	}
	var existeduser *models.User
	if getInfoToIs.LoginType == "phone_number" {
		existeduser, _ = m.userRepo.GetByUserEmail(ctx, "", getInfoToIs.LoginType, getInfoToIs.PhoneNumber)
		if existeduser == nil {
			return nil, models.ErrUnAuthorize
		}
	} else {
		existeduser, _ = m.userRepo.GetByUserEmail(ctx, getInfoToIs.Email, getInfoToIs.LoginType, "")
		if existeduser == nil {
			return nil, models.ErrUnAuthorize
		}
	}
	//if orderId != "" {
	//	booking, _ := m.bookingRepo.GetByID(ctx, orderId)
	//	if booking != nil && booking.UserId != nil && booking.IsRefund != nil {
	//		if *booking.UserId == existeduser.Id && booking.Status != 2  {
	//			m.paymentRepo.ConfirmPaymentRefund(ctx, orderId)
	//			if booking.Points != nil{
	//				booking.TotalPrice = booking.TotalPrice + *booking.Points
	//			}
	//			m.userRepo.UpdatePointByID(ctx, booking.TotalPrice, existeduser.Id, true)
	//
	//			existeduser.Points = existeduser.Points + int(booking.TotalPrice)
	//			if booking.TransId != nil {
	//				book, _ := m.bookingRepo.GetDetailTransportBookingID(ctx, booking.Id, booking.OrderId, nil)
	//				bookingDetail := book[0]
	//				creditUse := models.CreditsUse{
	//					Id:                    0,
	//					CreatedBy:             bookingDetail.BookedByEmail,
	//					CreatedDate:           time.Now(),
	//					ModifiedBy:            nil,
	//					ModifiedDate:          nil,
	//					DeletedBy:             nil,
	//					DeletedDate:           nil,
	//					IsDeleted:             0,
	//					IsActive:              1,
	//					UserId:                existeduser.Id,
	//					ReferralCode:          *bookingDetail.ReferralCode,
	//					TransactionId:         *bookingDetail.TransactionId,
	//					RuleReferralCodeId:    nil,
	//					Point:                 booking.TotalPrice,
	//					UserIdUseReferralCode: nil,
	//				}
	//				_, err = m.CreditsUseRepo.Insert(ctx, &creditUse)
	//				if err != nil {
	//					return nil, err
	//				}
	//			} else if booking.ExpId != nil {
	//				book, _ := m.bookingRepo.GetDetailBookingID(ctx, booking.Id, booking.Id)
	//				bookingDetail := book
	//				creditUse := models.CreditsUse{
	//					Id:                    0,
	//					CreatedBy:             bookingDetail.BookedByEmail,
	//					CreatedDate:           time.Now(),
	//					ModifiedBy:            nil,
	//					ModifiedDate:          nil,
	//					DeletedBy:             nil,
	//					DeletedDate:           nil,
	//					IsDeleted:             0,
	//					IsActive:              1,
	//					UserId:                existeduser.Id,
	//					ReferralCode:          *bookingDetail.ReferralCode,
	//					TransactionId:         *bookingDetail.TransactionId,
	//					RuleReferralCodeId:    nil,
	//					Point:                 booking.TotalPrice,
	//					UserIdUseReferralCode: nil,
	//				}
	//				_, err = m.CreditsUseRepo.Insert(ctx, &creditUse)
	//				if err != nil {
	//					return nil, err
	//				}
	//			}
	//		} else {
	//			return nil, errors.New("order id already refund")
	//		}
	//
	//	} else {
	//		// return nil, errors.New("invalid Order Id")
	//	}
	//
	//}
	userInfo := models.UserInfoDto{
		Id:             existeduser.Id,
		UserEmail:      existeduser.UserEmail,
		FirstName:       existeduser.FirstName,
		LastName:       existeduser.LastName,
		PhoneNumber:    existeduser.PhoneNumber,
		ProfilePictUrl: existeduser.ProfilePictUrl,
		ReferralCode:   existeduser.ReferralCode,
		Points:         existeduser.Points,
	}

	return &userInfo, nil
}

func (m userUsecase) Update(c context.Context, ar *models.NewCommandUser, isAdmin bool, token string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()
	var currentUser string
	username := ar.UserEmail
	if ar.LoginType == "phone_number" {
		username = ar.PhoneNumber
	}
	existeduser, _ := m.userRepo.GetByID(ctx, ar.Id, "")
	loginType := ""
	if existeduser.LoginType != nil {
		loginType = *existeduser.LoginType
	}
	updateUser := models.RegisterAndUpdateUser{}
	if isAdmin == true {
		//currentUserAdmin, err := m.adminUsecase.ValidateTokenAdmin(ctx, token)
		//if err != nil {
		//	return err
		//}
		//currentUser = currentUserAdmin.Name
		//updateUser = models.RegisterAndUpdateUser{
		//	Id:            ar.Id,
		//	Username:      username,
		//	Password:      ar.Password,
		//	Name:          ar.FullName,
		//	GivenName:     "",
		//	FamilyName:    "",
		//	Email:         ar.UserEmail,
		//	EmailVerified: true,
		//	Website:       "",
		//	Address:       "",
		//	OTP:           "",
		//	UserType:      1,
		//	PhoneNumber:   ar.PhoneNumber,
		//	LoginType:     loginType,
		//	UserRoles:     nil,
		//}
	} else {
		currentUsers, err := m.ValidateTokenUser(ctx, token)
		if err != nil {
			return err
		}
		currentUser = currentUsers.UserEmail

		getUserById, err := m.identityServerUc.GetDetailUserById(ar.Id, token, "true")
		updateUser = models.RegisterAndUpdateUser{
			Id:            ar.Id,
			Username:      username,
			Password:      getUserById.Password,
			Name:          ar.FirstName,
			GivenName:     "",
			FamilyName:    "",
			Email:         ar.UserEmail,
			EmailVerified: true,
			Website:       "",
			Address:       "",
			OTP:           "",
			UserType:      1,
			PhoneNumber:   ar.PhoneNumber,
			LoginType:     loginType,
			UserRoles:     nil,
		}
	}
	//var roles []string

	_, err := m.identityServerUc.UpdateUser(&updateUser)
	if err != nil {
		return err
	}
	var dob time.Time
	if ar.Dob != "" {

		layoutFormat := "2006-01-02 15:04:05"
		dobParse, errDateDob := time.Parse(layoutFormat, ar.Dob)
		if errDateDob != nil {
			return errDateDob
		}
		dob = dobParse
	}

	userModel := models.User{}
	userModel.Id = existeduser.Id
	userModel.ModifiedBy = &currentUser
	userModel.UserEmail = ar.UserEmail
	userModel.FirstName=       ar.FirstName
		userModel.LastName=       ar.LastName
	userModel.PhoneNumber = ar.PhoneNumber
	userModel.VerificationSendDate = existeduser.VerificationSendDate
	userModel.VerificationCode = existeduser.VerificationCode
	if ar.ProfilePictUrl != "" {
		userModel.ProfilePictUrl = ar.ProfilePictUrl
	} else {
		userModel.ProfilePictUrl = existeduser.ProfilePictUrl
	}

	userModel.Address = ar.Address
	userModel.Dob = dob
	userModel.Gender = ar.Gender
	userModel.IdType = ar.IdType
	userModel.IdNumber = ar.IdNumber
	userModel.ReferralCode = existeduser.ReferralCode
	userModel.Points = ar.Points
	userModel.LoginType = existeduser.LoginType
	return m.userRepo.Update(ctx, &userModel)
}
func replaceCapitalAlphabet(char string) string {
	char = strings.ReplaceAll(char, "a", "A")
	char = strings.ReplaceAll(char, "b", "B")
	char = strings.ReplaceAll(char, "c", "C")
	char = strings.ReplaceAll(char, "d", "D")
	char = strings.ReplaceAll(char, "e", "E")
	char = strings.ReplaceAll(char, "f", "F")
	char = strings.ReplaceAll(char, "g", "G")
	char = strings.ReplaceAll(char, "h", "H")
	char = strings.ReplaceAll(char, "i", "I")
	char = strings.ReplaceAll(char, "j", "J")
	char = strings.ReplaceAll(char, "k", "K")
	char = strings.ReplaceAll(char, "l", "L")
	char = strings.ReplaceAll(char, "m", "M")
	char = strings.ReplaceAll(char, "n", "N")
	char = strings.ReplaceAll(char, "o", "O")
	char = strings.ReplaceAll(char, "p", "P")
	char = strings.ReplaceAll(char, "q", "Q")
	char = strings.ReplaceAll(char, "r", "R")
	char = strings.ReplaceAll(char, "s", "S")
	char = strings.ReplaceAll(char, "t", "T")
	char = strings.ReplaceAll(char, "u", "U")
	char = strings.ReplaceAll(char, "v", "V")
	char = strings.ReplaceAll(char, "w", "W")
	char = strings.ReplaceAll(char, "x", "X")
	char = strings.ReplaceAll(char, "y", "Y")
	char = strings.ReplaceAll(char, "z", "Z")
	return char
}
func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes, err := generateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}
func (m userUsecase) Create(c context.Context, ar *models.NewCommandUser, isAdmin bool, token string) (*models.NewCommandUser, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	if ar.LoginType == "phone_number" {
		existeduser, _ := m.userRepo.GetByUserEmail(ctx, "", ar.LoginType, ar.PhoneNumber)
		if existeduser != nil {
			return nil, models.ErrConflict
		}
	} else {
		existeduser, _ := m.userRepo.GetByUserEmail(ctx, ar.UserEmail, ar.LoginType, "")
		if existeduser != nil {
			return nil, models.ErrConflict
		}
	}

	//var roles []string

	var createdBy string
	if isAdmin == true {
		if token == "" {
			return nil, models.ErrUnAuthorize
		} else {
			//currentUserAdmin, err := m.adminUsecase.ValidateTokenAdmin(ctx, token)
			//if err != nil {
			//	return nil, models.ErrUnAuthorize
			//}
			//createdBy = currentUserAdmin.Name
		}
	}
	username := ar.UserEmail
	if ar.LoginType == "phone_number" {
		username = ar.PhoneNumber
	}
	name := ar.FirstName
	name = strings.ReplaceAll(name, " ", "")
	name = strings.ReplaceAll(name, ".", "")
	name = strings.ReplaceAll(name, ",", "")
	name = string(name[0:4])
	name = replaceCapitalAlphabet(name)
	referralCode, er := generateRandomString(5)
	if er != nil {
		return nil, er
	}
	registerUser := models.RegisterAndUpdateUser{
		Id:            "",
		Username:      username,
		Password:      ar.Password,
		Name:          ar.FirstName,
		GivenName:     "",
		FamilyName:    "",
		Email:         ar.UserEmail,
		EmailVerified: false,
		Website:       "",
		Address:       "",
		OTP:           "",
		UserType:      1,
		PhoneNumber:   ar.PhoneNumber,
		LoginType:     ar.LoginType,
		UserRoles:     nil,
	}
	isUser, errorIs := m.identityServerUc.CreateUser(&registerUser)
	if isAdmin == false {

		createdBy = ar.UserEmail
	}

	ar.Id = isUser.Id
	var dob time.Time
	if ar.Dob != "" {

		layoutFormat := "2006-01-02 15:04:05"

		dobs, errDateDob := time.Parse(layoutFormat, ar.Dob)

		if errDateDob != nil {
			return nil, errDateDob
		}
		dob = dobs
	}

	if errorIs != nil {
		return nil, errorIs
	}

	referralCode = name + referralCode
	userModel := models.User{}
	userModel.Id = isUser.Id
	userModel.CreatedBy = createdBy
	userModel.UserEmail = ar.UserEmail
	userModel.FirstName = ar.FirstName
	userModel.LastName = ar.LastName
	userModel.PhoneNumber = ar.PhoneNumber
	userModel.VerificationSendDate = time.Now()
	userModel.VerificationCode = isUser.OTP
	userModel.ProfilePictUrl = ar.ProfilePictUrl
	userModel.Address = ar.Address
	userModel.Dob = dob
	userModel.Gender = ar.Gender
	userModel.IdType = ar.IdType
	userModel.IdNumber = ar.IdNumber
	userModel.ReferralCode = referralCode
	userModel.Points = ar.Points
	userModel.LoginType = &ar.LoginType
	err := m.userRepo.Insert(ctx, &userModel)
	if err != nil {
		return nil, err
	}
	requestToken, err := m.identityServerUc.GetToken(username, ar.Password, "", strconv.Itoa(registerUser.UserType), ar.LoginType)

	ar.Token = &requestToken.AccessToken
	return ar, nil
}

//func (m userUsecase) Subscription(ctx context.Context, s *models.NewCommandSubscribe) (*models.ResponseDelete, error) {
//	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
//	defer cancel()
//	subs := models.Subscribe{
//		Id:              0,
//		CreatedBy:       s.SubscriberEmail,
//		CreatedDate:     time.Time{},
//		ModifiedBy:      nil,
//		ModifiedDate:    nil,
//		DeletedBy:       nil,
//		DeletedDate:     nil,
//		IsDeleted:       0,
//		IsActive:        0,
//		SubscriberEmail: s.SubscriberEmail,
//	}
//	err := m.userRepo.SubscriptionUser(ctx, &subs)
//	if err != nil {
//		return nil, err
//	}
//
//	sendEmail := models.SendingEmail{
//		Subject: s.SubscriberEmail + " has subscribed to your mailing list.",
//		Message: s.SubscriberEmail + " has subscribed to cGO's mailing list. You can start sending latest news about cGO.",
//		From:    "",
//		To:      "info@cgo.co.id",
//	}
//	_, err = m.identityServerUc.SendingEmail(&sendEmail)
//	if err != nil {
//		return nil, err
//	}
//
//	result := models.ResponseDelete{
//		Id:      "",
//		Message: "Subscription Success",
//	}
//	//var roles []string
//	return &result, nil
//}

func (m userUsecase) GetCreditByID(ctx context.Context, id string) (*models.UserPoint, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	point, err := m.userRepo.GetCreditByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.UserPoint{Points: point}, nil
}

func (m userUsecase) GetUserDetailById(ctx context.Context, id string, token string) (*models.UserDto, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	getUserIdentity, err := m.identityServerUc.GetDetailUserById(id, token, "true")
	if err != nil {
		return nil, err
	}
	getUserById, err := m.userRepo.GetByID(ctx, id, "")

	result := models.UserDto{
		Id:             getUserById.Id,
		CreatedDate:    getUserById.CreatedDate,
		UpdatedDate:    getUserById.ModifiedDate,
		IsActive:       getUserById.IsActive,
		UserEmail:      getUserById.UserEmail,
		Password:       getUserIdentity.Password,
		FirstName:       getUserById.FirstName,
		LastName:       getUserById.LastName,
		PhoneNumber:    getUserById.PhoneNumber,
		ProfilePictUrl: getUserById.ProfilePictUrl,
		Address:        getUserById.Address,
		Dob:            getUserById.Dob,
		Gender:         getUserById.Gender,
		IdType:         getUserById.IdType,
		IdNumber:       getUserById.IdNumber,
		ReferralCode:   getUserById.ReferralCode,
		Points:         getUserById.Points,
	}

	return &result, nil

}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
