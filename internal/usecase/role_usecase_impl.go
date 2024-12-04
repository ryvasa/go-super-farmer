package usecase

import (
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/repository"
	"github.com/ryvasa/go-super-farmer/utils"
)

type RoleUsecaseImpl struct {
	repo repository.RoleRepository
}

func NewRoleUsecase(repo repository.RoleRepository) RoleUsecase {
	return &RoleUsecaseImpl{repo: repo}
}

func (u *RoleUsecaseImpl) CreateRole(req *dto.RoleCreateDTO) (*domain.Role, error) {
	role := domain.Role{}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return &role, utils.NewValidationError(err)
	}
	roles, err := u.repo.FindAll()
	if err != nil {
		return &role, utils.NewInternalError(err.Error())
	}

	role.ID = int64(len(*roles) + 1)
	role.Name = req.Name

	err = u.repo.Create(&role)
	if err != nil {
		return &role, utils.NewInternalError(err.Error())
	}

	createdRole, err := u.repo.FindByID(role.ID)
	if err != nil {
		return &role, utils.NewInternalError(err.Error())
	}

	return createdRole, nil
}

func (u *RoleUsecaseImpl) GetAllRoles() (*[]domain.Role, error) {
	roles, err := u.repo.FindAll()
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return roles, nil
}
