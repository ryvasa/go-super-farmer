package repository_implementation

import (
	"context"

	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	repository_interface "github.com/ryvasa/go-super-farmer/service_api/repository/interface"
	"gorm.io/gorm"
)

type RoleRepositoryImpl struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) repository_interface.RoleRepository {
	return &RoleRepositoryImpl{db}
}

func (r *RoleRepositoryImpl) Create(ctx context.Context, role *domain.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *RoleRepositoryImpl) FindAll(ctx context.Context) ([]*domain.Role, error) {
	var roles []*domain.Role
	if err := r.db.WithContext(ctx).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleRepositoryImpl) FindByID(ctx context.Context, id int64) (*domain.Role, error) {
	var role domain.Role
	if err := r.db.WithContext(ctx).Unscoped().First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
