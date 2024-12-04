package usecase

import (
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/repository"
)

type RoleUsecaseImpl struct {
	repo repository.RoleRepository
}

func NewRoleUsecase(repo repository.RoleRepository) RoleUsecase {
	return &RoleUsecaseImpl{repo: repo}
}

func (u *RoleUsecaseImpl) CreateRole(role *domain.Role) error {
	return u.repo.Create(role)
}

func (u *RoleUsecaseImpl) GetAllRoles() (*[]domain.Role, error) {
	return u.repo.FindAll()
}
