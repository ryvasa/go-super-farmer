package usecase

import (
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
)

type RoleUsecase interface {
	CreateRole(role *dto.RoleCreateDTO) (*domain.Role, error)
	GetAllRoles() (*[]domain.Role, error)
}
