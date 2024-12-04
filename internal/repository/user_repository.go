package repository

import (
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindById(id int64) (*domain.User, error)
	FindAll() (*[]domain.User, error)
	Delete(id int64) error
	Restore(id int64) error
	Update(id int64, user *domain.User) error
	FindDeletedById(id int64) (*domain.User, error)
}
