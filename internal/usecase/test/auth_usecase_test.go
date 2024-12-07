package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/repository/mock"
	"github.com/ryvasa/go-super-farmer/internal/usecase"
	"github.com/ryvasa/go-super-farmer/utils"
	mockTokenUtil "github.com/ryvasa/go-super-farmer/utils/mock"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	utilToken := mockTokenUtil.NewMockTokenUtil(ctrl)
	repo := mock.NewMockUserRepository(ctrl)
	uc := usecase.NewAuthUsecase(repo, utilToken)
	ctx := context.Background()

	t.Run("Test Login, successfully", func(t *testing.T) {
		userID := uuid.New()
		hashedPassword, _ := utils.HashPassword("123456")
		mockUser := &domain.User{
			ID:       userID,
			Email:    "test@gmail.com",
			Password: hashedPassword,
			Role:     domain.Role{ID: 1, Name: "user"},
		}
		mockToken := "mocked.jwt.token"

		repo.EXPECT().FindByEmail(ctx, "test@gmail.com").Return(mockUser, nil).Times(1)
		utilToken.EXPECT().GenerateToken(userID, "user").Return(mockToken, nil).Times(1)

		req := &dto.AuthDTO{Email: "test@gmail.com", Password: "123456"}
		resp, err := uc.Login(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, mockToken, resp.Token)
	})

	t.Run("Test Login, validation error", func(t *testing.T) {
		req := &dto.AuthDTO{Email: "", Password: "123456"}
		resp, err := uc.Login(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Test Login, error get user", func(t *testing.T) {
		repo.EXPECT().FindByEmail(ctx, "test@gmail.com").Return(nil, errors.New("internal error")).Times(1)

		req := &dto.AuthDTO{Email: "test@gmail.com", Password: "123456"}
		resp, err := uc.Login(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Test Login, error generate token", func(t *testing.T) {
		userID := uuid.New()
		hashedPassword, _ := utils.HashPassword("123456")
		mockUser := &domain.User{
			ID:       userID,
			Email:    "test@gmail.com",
			Password: hashedPassword,
			Role:     domain.Role{Name: "user"},
		}

		repo.EXPECT().FindByEmail(ctx, "test@gmail.com").Return(mockUser, nil).Times(1)
		utilToken.EXPECT().GenerateToken(userID, "user").Return("", errors.New("internal error")).Times(1)

		req := &dto.AuthDTO{Email: "test@gmail.com", Password: "123456"}
		resp, err := uc.Login(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
