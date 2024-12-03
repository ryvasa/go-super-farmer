package repository

import "github.com/ryvasa/go-super-farmer/internal/domain"

type UserRepository interface {
	Create(user *domain.User) error
	FindById(id uint) (*domain.User, error)
}
