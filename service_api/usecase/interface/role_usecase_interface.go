package usecase_interface

import (
	"context"

	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	"github.com/ryvasa/go-super-farmer/service_api/model/dto"
)

type RoleUsecase interface {
	CreateRole(ctx context.Context, role *dto.RoleCreateDTO) (*domain.Role, error)
	GetAllRoles(ctx context.Context) ([]*domain.Role, error)
}
