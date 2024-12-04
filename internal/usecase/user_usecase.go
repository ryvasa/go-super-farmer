package usecase

import (
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
)

type UserUsecase interface {
	Register(req *dto.UserCreateDTO) (*dto.UserResponseDTO, error)
	GetUserByID(id int64) (*dto.UserResponseDTO, error)
	GetAllUsers() (*[]dto.UserResponseDTO, error)
	UpdateUser(id int64, req *dto.UserUpdateDTO) (*dto.UserResponseDTO, error)
	DeleteUser(id int64) error
	RestoreUser(id int64) (*dto.UserResponseDTO, error)
}
