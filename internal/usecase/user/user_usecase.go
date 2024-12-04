package usecase

import (
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
)

type UserUsecase interface {
	Register(user *domain.User) error
	GetUserByID(id int64) (*domain.User, error)
	GetAllUsers() ([]domain.User, error)
	UpdateUser(id int64, user *domain.User) error
	DeleteUser(id int64) error
	RestoreUser(id int64) error
}
