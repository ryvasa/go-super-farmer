package usecase_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/repository/mock"
	usecase_implementation "github.com/ryvasa/go-super-farmer/internal/usecase/implementation"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	mockToken "github.com/ryvasa/go-super-farmer/pkg/auth/token/mock"
	"github.com/ryvasa/go-super-farmer/utils"
	mockUtils "github.com/ryvasa/go-super-farmer/utils/mock"
	"github.com/stretchr/testify/assert"
)

type AuthRepoMock struct {
	User  *mock.MockUserRepository
	Token *mockToken.MockToken
	Hash  *mockUtils.MockHasher
}

type AuthIDs struct {
	UserID uuid.UUID
}

type AuthMocks struct {
	User  *domain.User
	Auth  *dto.AuthResponseDTO
	Token string
}

type AuthDTOMock struct {
	Login *dto.AuthDTO
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
	}

	dto := &AuthDTOMock{
		Login: &dto.AuthDTO{
			Email:    "test@example.com",
			Password: "password",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	utilToken := mockToken.NewMockToken(ctrl)
	userRepo := mock.NewMockUserRepository(ctrl)
	hash := mockUtils.NewMockHasher(ctrl)
	uc := usecase_implementation.NewAuthUsecase(userRepo, utilToken, hash)
	ctx := context.TODO()

	repo := &AuthRepoMock{User: userRepo, Token: utilToken, Hash: hash}

	return ids, mocks, dto, repo, uc, ctx
}

func TestLogin(t *testing.T) {
	_, mocks, dtos, repo, uc, ctx := AuthUsecaseUtils(t)

	t.Run("should login successfully", func(t *testing.T) {
		repo.Hash.EXPECT().ValidatePassword(dtos.Login.Password, mocks.User.Password).Return(true).Times(1)
		repo.User.EXPECT().FindByEmail(ctx, dtos.Login.Email).Return(mocks.User, nil).Times(1)

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
