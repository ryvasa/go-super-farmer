package usecase_implementation

import (
	"context"
	"fmt"
	"time"

	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/pkg/auth/token"
	"github.com/ryvasa/go-super-farmer/pkg/database/cache"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	"github.com/ryvasa/go-super-farmer/pkg/messages"
	"github.com/ryvasa/go-super-farmer/utils"
)

type AuthUsecaseImpl struct {
	userRepo repository_interface.UserRepository
	token    token.Token
	hash     utils.Hasher
	rabbitMQ messages.RabbitMQ
	cache    cache.Cache
	OTP      utils.OTP
}

func NewAuthUsecase(userRepo repository_interface.UserRepository, token token.Token, hash utils.Hasher, rabbitMQ messages.RabbitMQ, cache cache.Cache, OTP utils.OTP) usecase_interface.AuthUsecase {
	return &AuthUsecaseImpl{userRepo, token, hash, rabbitMQ, cache, OTP}
}

func (u *AuthUsecaseImpl) Login(ctx context.Context, req *dto.AuthDTO) (*dto.AuthResponseDTO, error) {
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}

	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, utils.NewBadRequestError("invalid password or email")
	}

	logrus.Log.Infof("user: %+v", user)

	res := u.hash.ValidatePassword(req.Password, user.Password)
	if res == false {
		return nil, utils.NewBadRequestError("invalid password or email")
	}

	token, err := u.token.GenerateToken(user.ID, user.Role.Name)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return utils.AuthDtoFormat(user, token), nil
}

func (u *AuthUsecaseImpl) SendOTP(ctx context.Context, req *dto.AuthSendDTO) error {
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return utils.NewValidationError(err)
	}

	_, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return utils.NewBadRequestError("user not found")
	}

	// Generate OTP
	otp, err := u.OTP.GenerateOTP(6)
	if err != nil {
		return utils.NewInternalError("Failed to generate OTP")
	}

	// Simpan OTP di Redis dengan expiry 5 menit
	err = u.cache.Set(ctx, fmt.Sprintf("otp:%s", req.Email), []byte(otp), 5*time.Minute)
	if err != nil {
		return utils.NewInternalError("Failed to store OTP")
	}

	// Prepare email message
	msg := struct {
		To  string `json:"to"`
		OTP string `json:"otp"`
	}{
		To:  req.Email,
		OTP: otp,
	}

	// Publish ke RabbitMQ
	err = u.rabbitMQ.PublishJSON(ctx, "mail-exchange", "verify-email", msg)
	if err != nil {
		return utils.NewInternalError(err.Error())
	}

	return nil
}

func (u *AuthUsecaseImpl) VerifyOTP(ctx context.Context, req *dto.AuthVerifyDTO) error {
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return utils.NewValidationError(err)
	}
	// Ambil OTP dari Redis
	storedOTP, err := u.cache.Get(ctx, fmt.Sprintf("otp:%s", req.Email))
	if err != nil {
		return utils.NewInternalError("Failed to get OTP")
	}
	if storedOTP == nil {
		return utils.NewBadRequestError("OTP expired or not found")
	}

	// Verifikasi OTP
	if string(storedOTP) != req.OTP {
		return utils.NewBadRequestError("Invalid OTP")
	}

	// Hapus OTP setelah diverifikasi
	err = u.cache.Delete(ctx, fmt.Sprintf("otp:%s", req.Email))
	if err != nil {
		return utils.NewInternalError("Failed to delete OTP")
	}

	return nil
}
