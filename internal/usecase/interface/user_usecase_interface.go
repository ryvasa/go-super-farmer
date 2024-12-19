package usecase_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
)

type UserUsecase interface {
	Register(ctx context.Context, req *dto.UserCreateDTO) (*dto.UserResponseDTO, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*dto.UserResponseDTO, error)
	GetAllUsers(ctx context.Context, pagination *dto.PaginationDTO) (*dto.PaginationResponseDTO, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req *dto.UserUpdateDTO) (*dto.UserResponseDTO, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	RestoreUser(ctx context.Context, id uuid.UUID) (*dto.UserResponseDTO, error)
}
