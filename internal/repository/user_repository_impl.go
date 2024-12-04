package repository

import (
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(user *domain.User) error {
	return r.db.Create(user).Error
}
func (r *UserRepositoryImpl) FindById(id int64) (*domain.User, error) {
	var user domain.User
	err := r.db.
		Select("users.id", "users.name", "users.email", "users.phone", "users.created_at", "users.updated_at").Preload("Role", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).
		First(&user, id).Error
	return &user, err
}

func (r *UserRepositoryImpl) FindAll() (*[]domain.User, error) {
	var users []domain.User
	err := r.db.Find(&users).Error
	return &users, err
}

func (r *UserRepositoryImpl) Delete(id int64) error {
	return r.db.Delete(&domain.User{}, id).Error
}

func (r *UserRepositoryImpl) Restore(id int64) error {
	return r.db.Model(&domain.User{}).Where("deleted_at IS NOT NULL").Update("deleted_at", nil).Error
}

func (r *UserRepositoryImpl) Update(id int64, user *domain.User) error {
	return r.db.Model(&domain.User{}).Where("id = ?", id).Updates(user).Error
}

func (r *UserRepositoryImpl) FindDeletedById(id int64) (*domain.User, error) {
	var user domain.User
	err := r.db.Model(&user).Where("id = ?", id).Where("deleted_at IS NOT NULL").First(&user).Error
	return &user, err
}
