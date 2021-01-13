package user_admin

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	GetUserByEmail(ctx context.Context, email string) (*models.UserAdminDto, error)
	Delete(ctx context.Context, userId string, token string) (*models.ResponseDelete, error)
	Update(ctx context.Context, ar *models.NewCommandUserAdmin,  token string) error
	Create(ctx context.Context, ar *models.NewCommandUserAdmin,  token string) (*models.NewCommandUserAdmin, error)
	ValidateTokenUser(ctx context.Context, token string) (*models.UserAdminInfoDto, error)
	Login(ctx context.Context, ar *models.Login) (*models.GetToken, error)
	GetUserInfo(ctx context.Context, token string) (*models.UserAdminInfoDto, error)
	List(ctx context.Context, page, limit, offset int, search string) (*models.UserAdminWithPagination, error)
}
