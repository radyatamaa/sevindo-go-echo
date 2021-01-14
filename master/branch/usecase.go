package branch

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Delete(ctx context.Context, id string, token string) (*models.ResponseDelete, error)
	Update(ctx context.Context, ar *models.NewCommandBranch,  token string) error
	List(ctx context.Context, page, limit, offset int, search string) (*models.BranchWithPagination, error)
	Create(ctx context.Context, ar *models.NewCommandBranch, token string) (*models.NewCommandBranch, error)
	GetById(ctx context.Context, id string, token string) (*models.BranchDto, error)
}