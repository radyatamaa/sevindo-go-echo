package user

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	GetUserByReferralCode(ctx context.Context, code string) (*models.UserDto, error)
	CheckEmailORPhoneNumber(ctx context.Context, email string, phoneNumber string, otp string) (*models.ResponseDelete, error)
	ChangePassword(ctx context.Context, token string, password string, email string, phoneNumber string) (*models.ResponseDelete, error)
	GetUserByEmail(ctx context.Context, email string, loginType string) (*models.UserDto, error)
	LoginByGoogle(ctx context.Context, loginType string, email string, profilePicture string, name string) (*models.GetToken, error)
	Delete(ctx context.Context, userId string, token string) (*models.ResponseDelete, error)
	Update(ctx context.Context, ar *models.NewCommandUser, isAdmin bool, token string) error
	Create(ctx context.Context, ar *models.NewCommandUser, isAdmin bool, token string) (*models.NewCommandUser, error)
	ValidateTokenUser(ctx context.Context, token string) (*models.UserInfoDto, error)
	RequestOTP(ctx context.Context, phoneNumber string) (*models.RequestOTP, error)
	VerifiedEmail(ctx context.Context, token string, codeOTP string) (*models.UserInfoDto, error)
	Login(ctx context.Context, ar *models.Login) (*models.GetToken, error)
	GetUserInfo(ctx context.Context, token string, orderId string) (*models.UserInfoDto, error)
	GetCreditByID(ctx context.Context, id string) (*models.UserPoint, error)
	List(ctx context.Context, page, limit, offset int, search string) (*models.UserWithPagination, error)
	GetUserDetailById(ctx context.Context, id string, token string) (*models.UserDto, error)
	//Subscription(ctx context.Context, s *models.NewCommandSubscribe) (*models.ResponseDelete, error)
}
