package usecase_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	mock_repo "github.com/ryvasa/go-super-farmer/internal/repository/mock"
	usecase_implementation "github.com/ryvasa/go-super-farmer/internal/usecase/implementation"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	mockToken "github.com/ryvasa/go-super-farmer/pkg/auth/token/mock"
	mock_pkg "github.com/ryvasa/go-super-farmer/pkg/mock"
	"github.com/ryvasa/go-super-farmer/utils"
	mock_utils "github.com/ryvasa/go-super-farmer/utils/mock"
	"github.com/stretchr/testify/assert"
)

type AuthRepoMock struct {
	User     *mock_repo.MockUserRepository
	Token    *mockToken.MockToken
	Hash     *mock_utils.MockHasher
	RabbitMQ *mock_pkg.MockRabbitMQ
	Cache    *mock_pkg.MockCache
	OTP      *mock_utils.MockOTP
}

type AuthIDs struct {
	UserID uuid.UUID
}

type AuthMocks struct {
	User  *domain.User
	Auth  *dto.AuthResponseDTO
	Token string
	OTP   string
}

type AuthDTOMock struct {
	Login     *dto.AuthDTO
	SendOTP   *dto.AuthSendDTO
	VerifyOTP *dto.AuthVerifyDTO
}

func AuthUsecaseUtils(t *testing.T) (*AuthIDs, *AuthMocks, *AuthDTOMock, *AuthRepoMock, usecase_interface.AuthUsecase, context.Context) {
	userID := uuid.New()

	ids := &AuthIDs{
		UserID: userID,
	}

	mocks := &AuthMocks{
		User: &domain.User{
			ID:       userID,
			Name:     "test",
			Email:    "test@example.com",
			Password: "password",
		},
		Auth: &dto.AuthResponseDTO{
			User: &dto.UserResponseDTO{
				ID:    userID,
				Name:  "test",
				Email: "test@example.com",
			},
			Token: "mocked.jwt.token",
		},
		Token: "generated token",
		OTP:   "123456",
	}

	dto := &AuthDTOMock{
		Login: &dto.AuthDTO{
			Email:    "test@example.com",
			Password: "password",
		},
		SendOTP: &dto.AuthSendDTO{
			Email: "test@example.com",
		},
		VerifyOTP: &dto.AuthVerifyDTO{
			Email: "test@example.com",
			OTP:   "123456",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	utilToken := mockToken.NewMockToken(ctrl)
	userRepo := mock_repo.NewMockUserRepository(ctrl)
	hash := mock_utils.NewMockHasher(ctrl)
	rabbitMQ := mock_pkg.NewMockRabbitMQ(ctrl)
	cache := mock_pkg.NewMockCache(ctrl)
	otp := mock_utils.NewMockOTP(ctrl)
	uc := usecase_implementation.NewAuthUsecase(userRepo, utilToken, hash, rabbitMQ, cache, otp)
	ctx := context.TODO()

	repo := &AuthRepoMock{User: userRepo, Token: utilToken, Hash: hash, RabbitMQ: rabbitMQ, Cache: cache, OTP: otp}

	return ids, mocks, dto, repo, uc, ctx
}

func TestAuthUsecase_Login(t *testing.T) {
	_, mocks, dtos, repo, uc, ctx := AuthUsecaseUtils(t)

	t.Run("should login successfully", func(t *testing.T) {
		repo.User.EXPECT().FindByEmail(ctx, dtos.Login.Email).Return(mocks.User, nil).Times(1)
		repo.Hash.EXPECT().ValidatePassword(dtos.Login.Password, mocks.User.Password).Return(true).Times(1)

		repo.Token.EXPECT().GenerateToken(mocks.User.ID, mocks.User.Role.Name).Return(mocks.Token, nil).Times(1)

		resp, err := uc.Login(ctx, dtos.Login)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, mocks.Token, resp.Token)
		assert.Equal(t, mocks.User.ID, resp.User.ID)
		assert.Equal(t, mocks.User.Name, resp.User.Name)
		assert.Equal(t, mocks.User.Email, resp.User.Email)
	})

	t.Run("should return error validation error", func(t *testing.T) {
		resp, err := uc.Login(ctx, &dto.AuthDTO{Email: "", Password: "123456"})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("should return error get user by email", func(t *testing.T) {

		repo.User.EXPECT().FindByEmail(ctx, dtos.Login.Email).Return(nil, utils.NewBadRequestError("invalid password or email")).Times(1)

		resp, err := uc.Login(ctx, dtos.Login)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "invalid password or email")
	})

	t.Run("should return error generate token", func(t *testing.T) {
		repo.User.EXPECT().FindByEmail(ctx, mocks.User.Email).Return(mocks.User, nil).Times(1)

		repo.Hash.EXPECT().ValidatePassword(dtos.Login.Password, mocks.User.Password).Return(true).Times(1)

		repo.Token.EXPECT().GenerateToken(mocks.User.ID, mocks.User.Role.Name).Return("", utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.Login(ctx, dtos.Login)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("should return error when password is not valid", func(t *testing.T) {
		dtos.Login.Password = "123456"

		repo.Hash.EXPECT().ValidatePassword(dtos.Login.Password, mocks.User.Password).Return(false).Times(1)

		repo.User.EXPECT().FindByEmail(ctx, mocks.User.Email).Return(mocks.User, nil).Times(1)
		resp, err := uc.Login(ctx, dtos.Login)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "invalid password or email")
	})
}

func TestAuthUsecase_SendOTP(t *testing.T) {
	_, mocks, dtos, repo, uc, ctx := AuthUsecaseUtils(t)

	t.Run("should successfully send OTP", func(t *testing.T) {
		repo.User.EXPECT().FindByEmail(ctx, dtos.SendOTP.Email).Return(mocks.User, nil)
		repo.OTP.EXPECT().GenerateOTP(gomock.Any()).Return("", nil)
		repo.Cache.EXPECT().Set(ctx, gomock.Any(), gomock.Any(), 5*time.Minute).Return(nil)
		repo.RabbitMQ.EXPECT().PublishJSON(ctx, "mail-exchange", "verify-email", gomock.Any()).Return(nil)

		err := uc.SendOTP(ctx, dtos.SendOTP)
		assert.NoError(t, err)
	})

	t.Run("should return error when email validation fails", func(t *testing.T) {
		invalidDTO := &dto.AuthSendDTO{Email: ""}

		err := uc.SendOTP(ctx, invalidDTO)
		assert.Error(t, err)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		repo.User.EXPECT().FindByEmail(ctx, dtos.SendOTP.Email).
			Return(nil, utils.NewBadRequestError("user not found"))

		err := uc.SendOTP(ctx, dtos.SendOTP)
		assert.Error(t, err)
		assert.EqualError(t, err, "user not found")
	})

	t.Run("should return error when cache fails", func(t *testing.T) {
		repo.User.EXPECT().FindByEmail(ctx, dtos.SendOTP.Email).Return(mocks.User, nil)
		repo.OTP.EXPECT().GenerateOTP(gomock.Any()).Return("", nil)
		repo.Cache.EXPECT().Set(ctx, gomock.Any(), gomock.Any(), 5*time.Minute).
			Return(utils.NewInternalError("cache error"))

		err := uc.SendOTP(ctx, dtos.SendOTP)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to store OTP")
	})

	t.Run("should return error when rabbitmq fails", func(t *testing.T) {
		repo.User.EXPECT().FindByEmail(ctx, dtos.SendOTP.Email).Return(mocks.User, nil)
		repo.OTP.EXPECT().GenerateOTP(gomock.Any()).Return("", nil)
		repo.Cache.EXPECT().Set(ctx, gomock.Any(), gomock.Any(), 5*time.Minute).Return(nil)
		repo.RabbitMQ.EXPECT().PublishJSON(ctx, "mail-exchange", "verify-email", gomock.Any()).
			Return(utils.NewInternalError("rabbitmq error"))

		err := uc.SendOTP(ctx, dtos.SendOTP)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rabbitmq error")
	})

	t.Run("should return error when generate OTP fails", func(t *testing.T) {
		repo.User.EXPECT().FindByEmail(ctx, dtos.SendOTP.Email).Return(mocks.User, nil)
		repo.Cache.EXPECT().Set(ctx, gomock.Any(), gomock.Any(), 5*time.Minute).Return(nil)
		repo.RabbitMQ.EXPECT().PublishJSON(ctx, "mail-exchange", "verify-email", gomock.Any()).Return(nil)
		repo.OTP.EXPECT().GenerateOTP(gomock.Any()).Return("", utils.NewInternalError("Failed to generate OTP"))

		err := uc.SendOTP(ctx, dtos.SendOTP)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to generate OTP")
	})
}

func TestAuthUsecase_VerifyOTP(t *testing.T) {
	_, mocks, dtos, repo, uc, ctx := AuthUsecaseUtils(t)

	t.Run("should successfully verify OTP", func(t *testing.T) {
		key := fmt.Sprintf("otp:%s", mocks.User.Email)
		repo.Cache.EXPECT().Get(ctx, key).Return([]byte(mocks.OTP), nil)
		repo.User.EXPECT().FindByEmail(ctx, dtos.VerifyOTP.Email).Return(mocks.User, nil)
		repo.User.EXPECT().Update(ctx, mocks.User.ID, mocks.User).Return(nil)
		repo.Cache.EXPECT().Delete(ctx, key).Return(nil)

		err := uc.VerifyOTP(ctx, dtos.VerifyOTP)
		assert.NoError(t, err)
	})
	t.Run("should return error when validation fails", func(t *testing.T) {
		invalidDTO := &dto.AuthVerifyDTO{
			Email: "",
			OTP:   "",
		}

		err := uc.VerifyOTP(ctx, invalidDTO)
		assert.Error(t, err)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("should return error when OTP not found", func(t *testing.T) {
		key := fmt.Sprintf("otp:%s", dtos.VerifyOTP.Email)
		repo.Cache.EXPECT().Get(ctx, key).Return(nil, nil)

		err := uc.VerifyOTP(ctx, dtos.VerifyOTP)
		assert.Error(t, err)
		assert.EqualError(t, err, "OTP expired or not found")
	})

	t.Run("should return error when OTP is invalid", func(t *testing.T) {
		key := fmt.Sprintf("otp:%s", dtos.VerifyOTP.Email)
		repo.Cache.EXPECT().Get(ctx, key).Return([]byte("wrong-otp"), nil)

		err := uc.VerifyOTP(ctx, dtos.VerifyOTP)
		assert.Error(t, err)
		assert.EqualError(t, err, "Invalid OTP")
	})

	t.Run("should return error when cache delete fails", func(t *testing.T) {
		key := fmt.Sprintf("otp:%s", dtos.VerifyOTP.Email)
		repo.Cache.EXPECT().Get(ctx, key).Return([]byte(mocks.OTP), nil)
		repo.User.EXPECT().FindByEmail(ctx, dtos.VerifyOTP.Email).Return(mocks.User, nil)
		repo.User.EXPECT().Update(ctx, mocks.User.ID, mocks.User).Return(nil)
		repo.Cache.EXPECT().Delete(ctx, key).Return(utils.NewInternalError("cache error"))

		err := uc.VerifyOTP(ctx, dtos.VerifyOTP)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to delete OTP")
	})

	t.Run("should return error when cache get fails", func(t *testing.T) {
		key := fmt.Sprintf("otp:%s", dtos.VerifyOTP.Email)
		repo.Cache.EXPECT().Get(ctx, key).Return(nil, utils.NewInternalError("cache error"))

		err := uc.VerifyOTP(ctx, dtos.VerifyOTP)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to get OTP")
	})
}
