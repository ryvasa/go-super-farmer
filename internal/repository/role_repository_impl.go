package repository

import (
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"gorm.io/gorm"
)

type RoleRepositoryImpl struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &RoleRepositoryImpl{db}
}

func (r *RoleRepositoryImpl) Create(role *domain.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepositoryImpl) FindAll() (*[]domain.Role, error) {
	var roles []domain.Role
	if err := r.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return &roles, nil
}

func (r *RoleRepositoryImpl) FindByID(id int64) (*domain.Role, error) {
	var role domain.Role
	if err := r.db.First(&role, id).Error; err != nil {
		return &role, err
	}
	return &role, nil
}
