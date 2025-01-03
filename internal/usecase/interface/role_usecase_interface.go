package usecase_interface

import (
	"context"

	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
)

type RoleUsecase interface {
	CreateRole(ctx context.Context, role *dto.RoleCreateDTO) (*domain.Role, error)
	GetAllRoles(ctx context.Context) ([]*domain.Role, error)
}
