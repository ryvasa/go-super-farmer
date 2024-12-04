package repository

import "github.com/ryvasa/go-super-farmer/internal/model/domain"

type RoleRepository interface {
	Create(role *domain.Role) error
	FindAll() (*[]domain.Role, error)
}
