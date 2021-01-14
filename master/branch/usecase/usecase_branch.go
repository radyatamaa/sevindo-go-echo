package usecase

import (
	"context"
	"github.com/auth/user_admin"
	"github.com/master/branch"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/models"
)

type branchUsecase struct {
	userAdminUsecase user_admin.Usecase
	branchRepo    branch.Repository
	contextTimeout time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewbranchUsecase(	userAdminUsecase user_admin.Usecase,branchRepo branch.Repository, timeout time.Duration) branch.Usecase {
	return &branchUsecase{
		userAdminUsecase:userAdminUsecase,
		branchRepo:    branchRepo,
		contextTimeout: timeout,
	}
}

func (m branchUsecase) Delete(c context.Context, id string, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return nil,models.ErrUnAuthorize
	}

	err = m.branchRepo.Delete(ctx,id,currentUser.Email)

	result := &models.ResponseDelete{
		Id:      id,
		Message: "Success Delete",
	}

	return result,nil
}

func (m branchUsecase) Update(c context.Context, ar *models.NewCommandBranch, token string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return models.ErrUnAuthorize
	}

	getBranch ,err := m.branchRepo.GetByID(ctx,ar.Id)
	if err != nil{
		return err
	}
	var modifyBy string = currentUser.Email
	now := time.Now()
	getBranch.BranchName = ar.BranchName
	getBranch.BranchDesc = ar.BranchDesc
	getBranch.BranchPicture = ar.BranchPicture
	getBranch.Balance = ar.Balance
	getBranch.Address = ar.Address
	getBranch.ModifiedBy = &modifyBy
	getBranch.ModifiedDate = &now
	err = m.branchRepo.Update(ctx,getBranch)
	if err != nil{
		return err
	}
	return nil
}

func (m branchUsecase) List(ctx context.Context, page, limit, offset int, search string) (*models.BranchWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.branchRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	users := make([]*models.BranchDto, len(list))
	for i, item := range list {
		users[i] = &models.BranchDto{
			Id:          item.Id,
			BranchName: item.BranchName,
			BranchDesc: item.BranchDesc,
			BranchPicture: item.BranchPicture,
			Balance: item.Balance,
			Address: item.Address,
		}
	}
	totalRecords, _ := m.branchRepo.Count(ctx)
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

	response := &models.BranchWithPagination{
		Data: users,
		Meta: meta,
	}

	return response, nil
}

func (m branchUsecase) Create(c context.Context, ar *models.NewCommandBranch, token string) (*models.NewCommandBranch, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return nil,models.ErrUnAuthorize
	}

	insert := models.Branch{
		//Id:           0,
		Id:            uuid.New().String(),
		CreatedBy:    currentUser.Email,
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

	err = m.branchRepo.Insert(ctx, &insert)
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