package repository_implementation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository_interface.UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Select("users.id", "users.name", "users.email", "users.phone", "users.created_at", "users.updated_at").Preload("Role", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).
		First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindAll(ctx context.Context) (*[]domain.User, error) {
	var users []domain.User
	err := r.db.WithContext(ctx).Select("users.id", "users.name", "users.email", "users.phone", "users.created_at", "users.updated_at").Preload("Role", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).Find(&users).Error

	if err != nil {
		return nil, err
	}
	return &users, nil
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
	err := r.db.WithContext(ctx).Preload("Role", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).
		Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
