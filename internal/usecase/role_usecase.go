package usecase

import "github.com/ryvasa/go-super-farmer/internal/model/domain"

type RoleUsecase interface {
	CreateRole(role *domain.Role) error
	GetAllRoles() (*[]domain.Role, error)
}
