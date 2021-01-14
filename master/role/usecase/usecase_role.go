package usecase

import (
	"context"
	"github.com/auth/user_admin"
	"github.com/master/role"
	"github.com/models"
	"math"
	"strconv"
	"time"
)

type roleUsecase struct {
	userAdminUsecase user_admin.Usecase
	roleRepo    role.Repository
	contextTimeout time.Duration
}



// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewroleUsecase(	userAdminUsecase user_admin.Usecase,roleRepo role.Repository, timeout time.Duration) role.Usecase {
	return &roleUsecase{
		userAdminUsecase:userAdminUsecase,
		roleRepo:    roleRepo,
		contextTimeout: timeout,
	}
}
func (m roleUsecase) Delete(c context.Context, id int, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return nil,models.ErrUnAuthorize
	}

	err = m.roleRepo.Delete(ctx,id,currentUser.Email)

	result := &models.ResponseDelete{
		Id:      strconv.Itoa(id),
		Message: "Success Delete",
	}

	return result,nil
}

func (m roleUsecase) Update(c context.Context, ar *models.NewCommandRole, token string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return models.ErrUnAuthorize
	}

	getRole ,err := m.roleRepo.GetByID(ctx,ar.Id)
	if err != nil{
		return err
	}
	var modifyBy string = currentUser.Email
	now := time.Now()
	getRole.RoleName = ar.RoleName
	getRole.RoleType = ar.RoleType
	getRole.Description = ar.Description
	getRole.ModifiedBy = &modifyBy
	getRole.ModifiedDate = &now
	err = m.roleRepo.Update(ctx,getRole)
	if err != nil{
		return err
	}
	return nil
}

func (m roleUsecase) List(ctx context.Context, page, limit, offset int, search string) (*models.RoleWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.roleRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	users := make([]*models.RoleDto, len(list))
	for i, item := range list {
		users[i] = &models.RoleDto{
			Id:          item.Id,
			RoleName: item.RoleName,
			RoleType: item.RoleType,
			Description: item.Description,
		}
	}
	totalRecords, _ := m.roleRepo.Count(ctx)
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

	response := &models.RoleWithPagination{
		Data: users,
		Meta: meta,
	}

	return response, nil
}

func (m roleUsecase) Create(c context.Context, ar *models.NewCommandRole, token string) (*models.NewCommandRole, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	currentUser,err := m.userAdminUsecase.ValidateTokenUser(ctx,token)
	if err != nil{
		return nil,models.ErrUnAuthorize
	}

	insert := models.Role{
		Id:           0,
		CreatedBy:    currentUser.Email,
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		RoleName:  ar.RoleName,
		RoleType: ar.RoleType,
		Description: ar.Description,
	}

	err = m.roleRepo.Insert(ctx, &insert)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

func (m roleUsecase) GetById(c context.Context, id int, token string) (*models.RoleDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	role, err := m.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, models.ErrNotFound
	}

	result := &models.RoleDto{
		Id:          role.Id,
		RoleName: role.RoleName,
		RoleType: role.RoleType,
		Description: role.Description,
	}

	return result, nil
}
