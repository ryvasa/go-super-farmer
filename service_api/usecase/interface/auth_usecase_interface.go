package usecase_interface

import (
	"context"

	"github.com/ryvasa/go-super-farmer/service_api/model/dto"
)

type AuthUsecase interface {
	Login(ctx context.Context, req *dto.AuthDTO) (*dto.AuthResponseDTO, error)
	SendOTP(ctx context.Context, req *dto.AuthSendDTO) error
	VerifyOTP(ctx context.Context, req *dto.AuthVerifyDTO) error
}
