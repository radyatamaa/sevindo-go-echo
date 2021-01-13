package usecase

import (
	"context"
	"github.com/master/branch"
	"time"

	"github.com/models"
	"github.com/google/uuid"
)

type branchUsecase struct {
	branchRepo    branch.Repository
	contextTimeout time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewbranchUsecase(branchRepo branch.Repository, timeout time.Duration) branch.Usecase {
	return &branchUsecase{
		branchRepo:    branchRepo,
		contextTimeout: timeout,
	}
}

func (m branchUsecase) Create(c context.Context, ar *models.NewCommandBranch, token string) (*models.NewCommandBranch, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	insert := models.Branch{
		//Id:           ar.Id,
		Id:            uuid.New().String(),
		CreatedBy:    "admin",
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		BranchName:   ar.BranchName,
		BranchDesc:	  ar.BranchDesc,
		BranchPicture: ar.BranchPicture,
		Balance:	  ar.Balance,
		Address: 	  ar.Address,
	}

	err := m.branchRepo.Insert(ctx, &insert)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

func (m branchUsecase) GetById(c context.Context, id string, token string) (*models.BranchDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	branch, err := m.branchRepo.GetByID(ctx, id)
	if err != nil {
		return nil, models.ErrNotFound
	}

	result := &models.BranchDto{
		Id:          branch.Id,
		BranchName: branch.BranchName,
		BranchDesc: branch.BranchDesc,
		BranchPicture: branch.BranchPicture,
		Balance: branch.Balance,
		Address: branch.Address,
	}

	return result, nil
}