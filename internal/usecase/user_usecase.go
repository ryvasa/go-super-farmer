package usecase

import (
	"context"

	"github.com/ryvasa/go-super-farmer/internal/model/dto"
)

type UserUsecase interface {
	Register(ctx context.Context, req *dto.UserCreateDTO) (*dto.UserResponseDTO, error)
	GetUserByID(ctx context.Context, id uint64) (*dto.UserResponseDTO, error)
	GetAllUsers(ctx context.Context) (*[]dto.UserResponseDTO, error)
	UpdateUser(ctx context.Context, id uint64, req *dto.UserUpdateDTO) (*dto.UserResponseDTO, error)
	DeleteUser(ctx context.Context, id uint64) error
	RestoreUser(ctx context.Context, id uint64) (*dto.UserResponseDTO, error)
}
