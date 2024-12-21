package repository_implementation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	"github.com/ryvasa/go-super-farmer/utils"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository_interface.UserRepository {
	return &UserRepositoryImpl{db}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Select("users.id", "users.name", "users.email", "users.phone", "users.created_at", "users.updated_at").
		First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindAll(ctx context.Context, params *dto.PaginationDTO) ([]*domain.User, error) {
	var users []*domain.User

	err := r.db.WithContext(ctx).
		Scopes(
			utils.ApplyFilters(&params.Filter),
			utils.GetPaginationScope(params),
		).
		Select("users.id", "users.name", "users.email", "users.phone", "users.created_at", "users.updated_at").
		Find(&users).Error

	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.User{}).Error
}

func (r *UserRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Unscoped().Model(&domain.User{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *UserRepositoryImpl) Update(ctx context.Context, id uuid.UUID, user *domain.User) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Updates(user).Error
}

func (r *UserRepositoryImpl) FindDeletedByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Unscoped().Where("id = ? AND deleted_at IS NOT NULL", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		Preload("Role").
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) Count(ctx context.Context, filter *dto.PaginationFilterDTO) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Scopes(
			utils.ApplyFilters(filter),
		).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
