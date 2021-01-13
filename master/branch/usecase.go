package branch

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Create(ctx context.Context, ar *models.NewCommandBranch, token string) (*models.NewCommandBranch, error)
	GetById(ctx context.Context, id string, token string) (*models.BranchDto, error)
}