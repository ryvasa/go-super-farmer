package usecase_interface

import (
	"context"

	"github.com/ryvasa/go-super-farmer/internal/model/dto"
)

type AuthUsecase interface {
	Login(ctx context.Context, req *dto.AuthDTO) (*dto.AuthResponseDTO, error)
}
