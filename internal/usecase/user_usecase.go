package usecase

import (
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
)

type UserUsecase interface {
	Register(req *dto.UserCreateDTO) (*domain.User, error)
	GetUserByID(id int64) (*domain.User, error)
	GetAllUsers() (*[]domain.User, error)
	UpdateUser(id int64, req *dto.UserUpdateDTO) (*domain.User, error)
	DeleteUser(id int64) error
	RestoreUser(id int64) (*domain.User, error)
}
