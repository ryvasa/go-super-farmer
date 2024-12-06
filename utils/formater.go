package utils

import (
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
)

func UserDtoFormat(data *domain.User) *dto.UserResponseDTO {
	return &dto.UserResponseDTO{
		ID:        data.ID,
		Name:      data.Name,
		Email:     data.Email,
		Phone:     data.Phone,
		Password:  data.Password,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}
