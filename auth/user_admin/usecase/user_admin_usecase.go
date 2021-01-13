package usecase

import (
	"time"



	"github.com/auth/identityserver"
	"github.com/auth/user_admin"
	"github.com/models"
	"golang.org/x/net/context"
)

type userAdminUsecase struct {
	userAdminRepo         user_admin.Repository
	identityServerUc identityserver.Usecase
	contextTimeout   time.Duration
	tokenSystem string
}



// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewuserAdminUsecase( tokenSystem string,a user_admin.Repository, is identityserver.Usecase,  timeout time.Duration) user_admin.Usecase {
	return &userAdminUsecase{
		tokenSystem:tokenSystem,
		userAdminRepo:         a,
		identityServerUc: is,
		contextTimeout:   timeout,
	}
}

func (m userAdminUsecase) GetUserByEmail(ctx context.Context, email string) (*models.UserAdminDto, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()
	var getUser *models.UserAdmin

	getUser, _ = m.userAdminRepo.GetByUserEmail(ctx, email,false)


	var result *models.UserAdminDto
	if getUser != nil {
		result = &models.UserAdminDto{
			Id:       getUser.Id,
			Email:    getUser.Email,
			FullName: getUser.FullName,
			BranchId: getUser.BranchId,
		}
	}

	return result, nil
}

func (m userAdminUsecase) Delete(c context.Context, userId string, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()
	currentUserAdmin, err := m.ValidateTokenUser(ctx, token)
	if err != nil {
		if token == m.tokenSystem {
			currentUserAdmin = &models.UserAdminInfoDto{
				Id:       "",
				Email:    "system",
				FullName: "",
				BranchId: nil,
				Roles:    "",
			}

		}else {
			return nil, models.ErrUnAuthorize
		}

	}
	if currentUserAdmin.Roles != "super_admin" && token != m.tokenSystem{
		return nil, models.ErrUnAuthorize
	}
	error := m.userAdminRepo.Delete(ctx, userId, currentUserAdmin.Email)
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
	return &response,nil
}

func (m userAdminUsecase) Update(c context.Context, ar *models.NewCommandUserAdmin, token string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	var currentUser string = "system"
	username := ar.Email

	existeduser, _ := m.userAdminRepo.GetByID(ctx, ar.Id)

	updateUser := models.RegisterAndUpdateUser{}
	if ar.BranchId == nil{
		if token != m.tokenSystem{
			return models.ErrUnAuthorize
		}
		updateUser = models.RegisterAndUpdateUser{
			Id:            ar.Id,
			Username:      username,
			Password:      ar.Password,
			Name:          ar.FullName,
			GivenName:     "",
			FamilyName:    "",
			Email:         ar.Email,
			EmailVerified: true,
			Website:       "",
			Address:       "",
			OTP:           "",
			UserType:      2,
			PhoneNumber:   "",
			LoginType:     "email",
			UserRoles:     nil,
		}
	} else {
		currentUsers, err := m.ValidateTokenUser(ctx, token)
		if err != nil {
			return models.ErrUnAuthorize
		}
		currentUser = currentUsers.Email

		getUserById, err := m.identityServerUc.GetDetailUserById(ar.Id, token, "true")
		updateUser = models.RegisterAndUpdateUser{
			Id:            ar.Id,
			Username:      username,
			Password:      getUserById.Password,
			Name:          ar.FullName,
			GivenName:     "",
			FamilyName:    "",
			Email:         ar.Email,
			EmailVerified: true,
			Website:       "",
			Address:       "",
			OTP:           "",
			UserType:      2,
			PhoneNumber:   "",
			LoginType:     "email",
			UserRoles:     nil,
		}
	}
	//var roles []string

	_, err := m.identityServerUc.UpdateUser(&updateUser)
	if err != nil {
		return err
	}
	now := time.Now()
	existeduser.ModifiedDate = &now
	existeduser.ModifiedBy = &currentUser
	existeduser.Email = ar.Email
	existeduser.BranchId = ar.BranchId
	existeduser.FullName = ar.FullName
	return m.userAdminRepo.Update(ctx, existeduser)
}

func (m userAdminUsecase) Create(c context.Context, ar *models.NewCommandUserAdmin,  token string) (*models.NewCommandUserAdmin, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()
	var currentUser string = "system"
	var createuser models.RegisterAndUpdateUser
	if ar.BranchId == nil{
		if token != m.tokenSystem{
			return nil, models.ErrUnAuthorize
		}
		createuser = models.RegisterAndUpdateUser{
			Id:            ar.Id,
			Username:      ar.Email,
			Password:      ar.Password,
			Name:          ar.FullName,
			GivenName:     "",
			FamilyName:    "",
			Email:         ar.Email,
			EmailVerified: true,
			Website:       "",
			Address:       "",
			OTP:           "",
			UserType:      2,
			PhoneNumber:   "",
			LoginType:     "email",
			UserRoles:     nil,
		}
	}else {
		createuser = models.RegisterAndUpdateUser{
			Id:            ar.Id,
			Username:      ar.Email,
			Password:      ar.Password,
			Name:          ar.FullName,
			GivenName:     "",
			FamilyName:    "",
			Email:         ar.Email,
			EmailVerified: true,
			Website:       "",
			Address:       "",
			OTP:           "",
			UserType:      2,
			PhoneNumber:   "",
			LoginType:     "email",
			UserRoles:     nil,
		}
		user,err := m.ValidateTokenUser(ctx,token)
		if err != nil{
			return nil,models.ErrUnAuthorize
		}
		currentUser = user.Email
	}

	res, err := m.identityServerUc.CreateUser(&createuser)
	if err != nil{
		return nil,err
	}

	useradmin := models.UserAdmin{
		Id:           res.Id,
		CreatedBy:    currentUser,
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		Email:        ar.Email,
		FullName:     ar.FullName,
		BranchId:     ar.BranchId,
	}
	err = m.userAdminRepo.Insert(ctx,&useradmin)
	if err != nil{
		return nil,err
	}
	ar.Id = useradmin.Id
	return ar,nil
}

func (m userAdminUsecase) ValidateTokenUser(ctx context.Context, token string) (*models.UserAdminInfoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	getInfoToIs, err := m.identityServerUc.GetUserInfo(token)
	if err != nil {
		return nil, err
	}

	existeduser, _ := m.userAdminRepo.GetByUserEmail(ctx, getInfoToIs.Username,true)
	if existeduser == nil {
		existeduser, _ := m.userAdminRepo.GetByUserEmail(ctx, getInfoToIs.Username,false)
		if existeduser == nil {
			return nil, models.ErrUnAuthorize
		}
	}

	var roles string = "super_admin"
	if existeduser.BranchId != nil {
		roles = "admin_branch"
	}
	userInfo := &models.UserAdminInfoDto{
		Id:       existeduser.Id,
		Email:    existeduser.Email,
		FullName: existeduser.FullName,
		BranchId: existeduser.BranchId,
		Roles:roles,
	}
	return userInfo, nil
}

func (m userAdminUsecase) Login(ctx context.Context, ar *models.Login) (*models.GetToken, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()
	ar.LoginType = "email"
	getToken, err := m.identityServerUc.GetToken(ar.Email, ar.Password, ar.Scope, "2", ar.LoginType)
	if err != nil {
		return nil, err
	}
	existeduser, _ := m.userAdminRepo.GetByUserEmail(ctx, ar.Email, true)
	if existeduser == nil {
		existeduser, _ = m.userAdminRepo.GetByUserEmail(ctx, ar.Email, false)
		if existeduser == nil {
			return nil, models.ErrUsernamePassword
		}
	}

	return getToken,nil
}

func (m userAdminUsecase) GetUserInfo(ctx context.Context, token string) (*models.UserAdminInfoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()
	return m.ValidateTokenUser(ctx,token)
}

func (u userAdminUsecase) List(ctx context.Context, page, limit, offset int, search string) (*models.UserAdminWithPagination, error) {
	panic("implement me")
}