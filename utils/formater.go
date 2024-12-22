package utils

import (
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	"github.com/ryvasa/go-super-farmer/service_api/model/dto"
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

// func PriceDtoFormat(data *domain.Price) *dto.PriceResponseDTO {
// 	return &dto.PriceResponseDTO{
// 		ID:        data.ID,
// 		Price:     data.Price,
// 		CreatedAt: data.CreatedAt,
// 		UpdatedAt: data.UpdatedAt,
// 	}
// }

func AuthDtoFormat(user *domain.User, token string) *dto.AuthResponseDTO {
	return &dto.AuthResponseDTO{
		User:  UserDtoFormat(user),
		Token: token,
	}
}
