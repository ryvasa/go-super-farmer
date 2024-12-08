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
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockUserRepository(ctrl)
	uc := usecase.NewUserUsecase(repo)
	ctx := context.TODO()
	t.Run("Test Register successfully", func(t *testing.T) {
		userID := uuid.New()

		req := &dto.UserCreateDTO{Name: "Test User", Email: "test@example.com", Password: "password"}
		mockUser := &domain.User{
			Name:  "Test User",
			Email: "test@example.com",
		}

		// Mock HashPassword untuk mengembalikan password yang di-hash
		utils.MockHashPassword = func(password string) (string, error) {
			return "hashed_password", nil
		}
		defer func() { utils.MockHashPassword = nil }() // Reset mock setelah test selesai

		repo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *domain.User) error {
			user.ID = userID
			return nil
		}).Times(1)
		repo.EXPECT().FindByID(ctx, userID).Return(mockUser, nil).Times(1)

		resp, err := uc.Register(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, req.Name, resp.Name)
		assert.Equal(t, req.Email, resp.Email)
	})

	t.Run("Test Register validation error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockUserRepository(ctrl)
		uc := usecase.NewUserUsecase(repo)
		ctx := context.TODO()

		req := &dto.UserCreateDTO{Email: "", Password: ""}

		resp, err := uc.Register(ctx, req)

		assert.Error(t, err)
		assert.Empty(t, resp)
	})

	t.Run("Test Register hashing error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockUserRepository(ctrl)
		uc := usecase.NewUserUsecase(repo)
		ctx := context.TODO()

		req := &dto.UserCreateDTO{Name: "Test User", Email: "test@example.com", Password: "password"}

		// Mock HashPassword untuk mengembalikan password yang di-hash
		utils.MockHashPassword = func(password string) (string, error) {
			return "", errors.New("hashing error")
		}
		defer func() { utils.MockHashPassword = nil }() // Reset mock setelah test selesai

		repo.EXPECT().Create(ctx, gomock.Any()).Times(0)

		res, err := uc.Register(ctx, req)

		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockUserRepository(ctrl)
	uc := usecase.NewUserUsecase(repo)
	ctx := context.Background()

	t.Run("Test GetUserByID successfully", func(t *testing.T) {
		userID := uuid.New()

		mockUser := &domain.User{ID: userID, Name: "Test User", Email: "test@example.com"}

		repo.EXPECT().FindByID(ctx, mockUser.ID).Return(mockUser, nil).Times(1)

		resp, err := uc.GetUserByID(ctx, mockUser.ID)

		assert.NoError(t, err)
		assert.Equal(t, mockUser.Name, resp.Name)
		assert.Equal(t, mockUser.Email, resp.Email)
	})

	t.Run("Test GetUserByID not found", func(t *testing.T) {
		userID := uuid.New()

		repo.EXPECT().FindByID(ctx, userID).Return(nil, errors.New("user not found")).Times(1)

		resp, err := uc.GetUserByID(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestGetAllUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockUserRepository(ctrl)
	uc := usecase.NewUserUsecase(repo)
	ctx := context.Background()

	t.Run("Test GetAllUsers successfully", func(t *testing.T) {
		userID1 := uuid.New()
		userID2 := uuid.New()

		mockUsers := &[]domain.User{
			{ID: userID1, Name: "Test User 1", Email: "test1@example.com"},
			{ID: userID2, Name: "Test User 2", Email: "test2@example.com"},
		}

		repo.EXPECT().FindAll(ctx).Return(mockUsers, nil).Times(1)

		resp, err := uc.GetAllUsers(ctx)

		assert.NoError(t, err)
		assert.Len(t, *resp, len(*mockUsers))
	})

	t.Run("Test GetAllUsers internal error", func(t *testing.T) {
		repo.EXPECT().FindAll(ctx).Return(nil, errors.New("internal error")).Times(1)

		resp, err := uc.GetAllUsers(ctx)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockUserRepository(ctrl)
	uc := usecase.NewUserUsecase(repo)
	ctx := context.Background()

	t.Run("Test UpdateUser validation error", func(t *testing.T) {
		userID := uuid.New()
		updateReq := &dto.UserUpdateDTO{Name: "k", Email: "test@example.com", Password: "password"}

		resp, err := uc.UpdateUser(ctx, userID, updateReq)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Test UpdateUser witout password, successfully", func(t *testing.T) {
		userID := uuid.New()

		mockUser := &domain.User{ID: userID, Name: "Test User", Email: "test@example.com"}
		mockUserUpdated := &domain.User{ID: userID, Name: "Test User Updated", Email: "test@example.com"}

		repo.EXPECT().FindByID(ctx, mockUser.ID).Return(mockUser, nil).Times(1)
		repo.EXPECT().Update(ctx, mockUser.ID, gomock.Any()).Return(nil).Times(1)
		repo.EXPECT().FindByID(ctx, mockUser.ID).Return(mockUserUpdated, nil).Times(1)

		req := &dto.UserUpdateDTO{Name: "Test User Updated", Email: "test@example.com"}
		resp, err := uc.UpdateUser(ctx, mockUser.ID, req)

		assert.NoError(t, err)
		assert.Equal(t, req.Name, resp.Name)
		assert.Equal(t, req.Email, resp.Email)
	})

	t.Run("Test UpdateUser with password, successfully", func(t *testing.T) {
		userID := uuid.New()

		mockUser := &domain.User{ID: userID, Name: "Test User", Email: "test@example.com"}
		mockUserUpdated := &domain.User{ID: userID, Name: "Test User Updated", Email: "test@example.com"}

		repo.EXPECT().FindByID(ctx, mockUser.ID).Return(mockUser, nil).Times(1)
		repo.EXPECT().Update(ctx, mockUser.ID, gomock.Any()).Return(nil).Times(1)
		repo.EXPECT().FindByID(ctx, mockUser.ID).Return(mockUserUpdated, nil).Times(1)

		// Mock HashPassword untuk mengembalikan password yang di-hash
		utils.MockHashPassword = func(password string) (string, error) {
			return "hashed_password", nil
		}
		defer func() { utils.MockHashPassword = nil }() // Reset mock setelah test selesai

		req := &dto.UserUpdateDTO{Name: "Test User Updated", Email: "test@example.com", Password: "password"}
		resp, err := uc.UpdateUser(ctx, mockUser.ID, req)

		assert.NoError(t, err)
		assert.Equal(t, req.Name, resp.Name)
		assert.Equal(t, req.Email, resp.Email)
	})

	t.Run("Test UpdateUser with password, hashing error", func(t *testing.T) {
		userID := uuid.New()

		mockUser := &domain.User{ID: userID, Name: "Test User", Email: "test@example.com"}

		repo.EXPECT().FindByID(ctx, mockUser.ID).Return(mockUser, nil).Times(1)

		// Mock HashPassword untuk mengembalikan password yang di-hash
		utils.MockHashPassword = func(password string) (string, error) {
			return "", errors.New("hashing error")
		}
		defer func() { utils.MockHashPassword = nil }() // Reset mock setelah test selesai

		req := &dto.UserUpdateDTO{Name: "Test User Updated", Email: "test@example.com", Password: "password"}
		resp, err := uc.UpdateUser(ctx, mockUser.ID, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Test UpdateUser not found", func(t *testing.T) {
		userID := uuid.New()

		repo.EXPECT().FindByID(ctx, userID).Return(nil, errors.New("user not found")).Times(1)

		req := &dto.UserUpdateDTO{Name: "Test User Updated", Email: "test@example.com", Password: "password"}
		resp, err := uc.UpdateUser(ctx, userID, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockUserRepository(ctrl)
	uc := usecase.NewUserUsecase(repo)
	ctx := context.Background()

	t.Run("Test DeleteUser successfully", func(t *testing.T) {
		userID := uuid.New()

		mockUser := &domain.User{ID: userID, Name: "Test User", Email: "test@example.com"}

		repo.EXPECT().FindByID(ctx, mockUser.ID).Return(mockUser, nil).Times(1)
		repo.EXPECT().Delete(ctx, mockUser.ID).Return(nil).Times(1)

		err := uc.DeleteUser(ctx, mockUser.ID)

		assert.NoError(t, err)
	})

	t.Run("Test DeleteUser not found", func(t *testing.T) {
		userID := uuid.New()

		repo.EXPECT().FindByID(ctx, userID).Return(nil, errors.New("user not found")).Times(1)

		err := uc.DeleteUser(ctx, userID)

		assert.Error(t, err)
	})
}

func TestRestoreUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockUserRepository(ctrl)
	uc := usecase.NewUserUsecase(repo)
	ctx := context.Background()

	t.Run("Test RestoreUser successfully", func(t *testing.T) {
		userID := uuid.New()

		mockUser := &domain.User{ID: userID, Name: "Test User", Email: "test@example.com"}

		repo.EXPECT().FindDeletedByID(ctx, mockUser.ID).Return(mockUser, nil).Times(1)
		repo.EXPECT().Restore(ctx, mockUser.ID).Return(nil).Times(1)
		repo.EXPECT().FindByID(ctx, mockUser.ID).Return(mockUser, nil).Times(1)

		resp, err := uc.RestoreUser(ctx, mockUser.ID)

		assert.NoError(t, err)
		assert.Equal(t, mockUser.Name, resp.Name)
		assert.Equal(t, mockUser.Email, resp.Email)
	})

	t.Run("Test RestoreUser not found", func(t *testing.T) {
		userID := uuid.New()

		repo.EXPECT().FindDeletedByID(ctx, userID).Return(nil, errors.New("user not found")).Times(1)

		resp, err := uc.RestoreUser(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
