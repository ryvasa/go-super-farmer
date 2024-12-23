package usecase_implementation

import (
	"context"

	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	"github.com/ryvasa/go-super-farmer/service_api/model/dto"
	repository_interface "github.com/ryvasa/go-super-farmer/service_api/repository/interface"
	usecase_interface "github.com/ryvasa/go-super-farmer/service_api/usecase/interface"
	"github.com/ryvasa/go-super-farmer/utils"
)

type RoleUsecaseImpl struct {
	repo repository_interface.RoleRepository
}

func NewRoleUsecase(repo repository_interface.RoleRepository) usecase_interface.RoleUsecase {
	return &RoleUsecaseImpl{repo: repo}
}

func (u *RoleUsecaseImpl) CreateRole(ctx context.Context, req *dto.RoleCreateDTO) (*domain.Role, error) {
	role := domain.Role{}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}
	roles, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	role.ID = int64(len(roles) + 1)
	role.Name = req.Name

	err = u.repo.Create(ctx, &role)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	createdRole, err := u.repo.FindByID(ctx, role.ID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return createdRole, nil
}

func (u *RoleUsecaseImpl) GetAllRoles(ctx context.Context) ([]*domain.Role, error) {
	roles, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return roles, nil
}
