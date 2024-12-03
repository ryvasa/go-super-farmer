package usecase

import "github.com/ryvasa/go-super-farmer/internal/model/domain"

type UserUsecase interface {
	Register(user *domain.User) error
	GetUserByID(id uint) (*domain.User, error)
}
