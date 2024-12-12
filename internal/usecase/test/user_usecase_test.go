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
	mockUtils "github.com/ryvasa/go-super-farmer/utils/mock"
	"github.com/stretchr/testify/assert"
)

type UserRepoMock struct {
	User *mock.MockUserRepository
	Hash *mockUtils.MockHasher
}

type UserIDs struct {
	UserID uuid.UUID
}

type UserMocks struct {
	Users       *[]domain.User
	User        *domain.User
	UpdatedUser *domain.User
}

type MockUserDTOs struct {
	Create             *dto.UserCreateDTO
	Update             *dto.UserUpdateDTO
	UpdateWithPassword *dto.UserUpdateDTO
}

func UserUsecaseUtils(t *testing.T) (*UserIDs, *UserMocks, *MockUserDTOs, *UserRepoMock, usecase.UserUsecase, context.Context) {
	userID := uuid.New()

	ids := &UserIDs{
		UserID: userID,
	}

	mocks := &UserMocks{
		Users: &[]domain.User{
			{
				ID:    userID,
				Name:  "test",
				Email: "test@example.com",
			},
		},
		User: &domain.User{
			ID:    userID,
			Name:  "test",
			Email: "test@example.com",
		},
		UpdatedUser: &domain.User{
			ID:   userID,
			Name: "updated",
		},
	}

	dto := &MockUserDTOs{
		Create: &dto.UserCreateDTO{
			Name:     "test",
			Email:    "test@example.com",
			Password: "password",
			Phone:    "1111111",
		},
		Update: &dto.UserUpdateDTO{
			Name: "updated",
		},

		UpdateWithPassword: &dto.UserUpdateDTO{
			Name:     "updated",
			Password: "password",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	hash := mockUtils.NewMockHasher(ctrl)
	uc := usecase.NewUserUsecase(userRepo, hash)
	ctx := context.TODO()

	repo := &UserRepoMock{User: userRepo, Hash: hash}

	return ids, mocks, dto, repo, uc, ctx
}

func TestRegister(t *testing.T) {
	ids, mocks, dtos, repo, uc, ctx := UserUsecaseUtils(t)

	t.Run("should register successfully", func(t *testing.T) {

		repo.Hash.EXPECT().HashPassword(dtos.Create.Password).Return("hashed_password", nil).Times(1)

		repo.User.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *domain.User) error {
			user.ID = ids.UserID
			return nil
		}).Times(1)
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)

		resp, err := uc.Register(ctx, dtos.Create)

		assert.NoError(t, err)
		assert.Equal(t, dtos.Create.Name, resp.Name)
		assert.Equal(t, dtos.Create.Email, resp.Email)
	})

	t.Run("should return error when validation error", func(t *testing.T) {
		resp, err := uc.Register(ctx, &dto.UserCreateDTO{Name: "", Email: "", Password: ""})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, err.Error(), "Validation failed")
	})

	t.Run("Test Register hashing error", func(t *testing.T) {
		repo.Hash.EXPECT().HashPassword(dtos.Create.Password).Return("", utils.NewInternalError("hashing error")).Times(1)

		resp, err := uc.Register(ctx, dtos.Create)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "hashing error")
	})

	t.Run("should return error when create user", func(t *testing.T) {
		repo.Hash.EXPECT().HashPassword(dtos.Create.Password).Return("hashed_password", nil).Times(1)

		repo.User.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("create user error")).Times(1)

		resp, err := uc.Register(ctx, dtos.Create)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "create user error")
	})

	t.Run("should return error when find created user by id", func(t *testing.T) {
		repo.Hash.EXPECT().HashPassword(dtos.Create.Password).Return("hashed_password", nil).Times(1)

		repo.User.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *domain.User) error {
			user.ID = ids.UserID
			return nil
		}).Times(1)
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(nil, errors.New("find user by id error")).Times(1)

		resp, err := uc.Register(ctx, dtos.Create)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "find user by id error")
	})
}

func TestGetUserByID(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := UserUsecaseUtils(t)

	t.Run("Test GetUserByID sUCcessfully", func(t *testing.T) {

		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)

		resp, err := uc.GetUserByID(ctx, ids.UserID)

		assert.NoError(t, err)
		assert.Equal(t, mocks.User.ID, resp.ID)
		assert.Equal(t, mocks.User.Name, resp.Name)
		assert.Equal(t, mocks.User.Email, resp.Email)
	})

	t.Run("Test GetUserByID not found", func(t *testing.T) {
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(nil, utils.NewNotFoundError("user not found")).Times(1)

		resp, err := uc.GetUserByID(ctx, ids.UserID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "user not found")
	})
}

func TestGetAllUsers(t *testing.T) {
	_, mocks, _, repo, uc, ctx := UserUsecaseUtils(t)

	t.Run("should return all users", func(t *testing.T) {
		repo.User.EXPECT().FindAll(ctx).Return(mocks.Users, nil).Times(1)

		resp, err := uc.GetAllUsers(ctx)

		assert.NoError(t, err)
		assert.Len(t, *resp, len(*mocks.Users))
		assert.Equal(t, (*mocks.Users)[0].ID, (*resp)[0].ID)
		assert.Equal(t, (*mocks.Users)[0].Name, (*resp)[0].Name)
		assert.Equal(t, (*mocks.Users)[0].Email, (*resp)[0].Email)
	})

	t.Run("should return error internal error", func(t *testing.T) {
		repo.User.EXPECT().FindAll(ctx).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetAllUsers(ctx)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})
}

func TestUpdateUser(t *testing.T) {
	ids, mocks, dtos, repo, uc, ctx := UserUsecaseUtils(t)

	t.Run("should update user without password", func(t *testing.T) {
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)

		repo.User.EXPECT().Update(ctx, ids.UserID, gomock.Any()).Return(nil).Times(1)

		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.UpdatedUser, nil).Times(1)

		resp, err := uc.UpdateUser(ctx, ids.UserID, dtos.Update)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, mocks.UpdatedUser.ID, resp.ID)
		assert.Equal(t, dtos.Update.Name, resp.Name)
	})

	t.Run("should update user with password", func(t *testing.T) {
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)

		repo.User.EXPECT().Update(ctx, ids.UserID, gomock.Any()).Return(nil).Times(1)

		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.UpdatedUser, nil).Times(1)

		repo.Hash.EXPECT().HashPassword(dtos.UpdateWithPassword.Password).Return("hashed_password", nil).Times(1)

		resp, err := uc.UpdateUser(ctx, ids.UserID, dtos.UpdateWithPassword)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, mocks.UpdatedUser.ID, resp.ID)
		assert.Equal(t, dtos.UpdateWithPassword.Name, resp.Name)
	})

	t.Run("should return error validation error", func(t *testing.T) {
		resp, err := uc.UpdateUser(ctx, ids.UserID, &dto.UserUpdateDTO{Name: "1", Email: "", Password: ""})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, err.Error(), "Validation failed")
	})

	t.Run("should return error hashing error", func(t *testing.T) {

		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)

		repo.Hash.EXPECT().HashPassword(dtos.UpdateWithPassword.Password).Return("", errors.New("hashing error")).Times(1)

		resp, err := uc.UpdateUser(ctx, ids.UserID, dtos.UpdateWithPassword)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "hashing error")
	})

	t.Run("should return error not found", func(t *testing.T) {
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(nil, utils.NewNotFoundError("user not found")).Times(1)

		resp, err := uc.UpdateUser(ctx, ids.UserID, dtos.UpdateWithPassword)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "user not found")
	})

	t.Run("should return error not found when find updated user", func(t *testing.T) {
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)

		repo.Hash.EXPECT().HashPassword(dtos.UpdateWithPassword.Password).Return("hashed_password", nil).Times(1)

		repo.User.EXPECT().Update(ctx, ids.UserID, gomock.Any()).Return(nil).Times(1)
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(nil, utils.NewInternalError("user not found")).Times(1)

		resp, err := uc.UpdateUser(ctx, ids.UserID, dtos.UpdateWithPassword)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "user not found")
	})
}

func TestDeleteUser(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := UserUsecaseUtils(t)

	t.Run("should delete user sUCcessfully", func(t *testing.T) {
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)

		repo.User.EXPECT().Delete(ctx, ids.UserID).Return(nil).Times(1)

		err := uc.DeleteUser(ctx, ids.UserID)

		assert.NoError(t, err)
	})

	t.Run("should return error not found", func(t *testing.T) {
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(nil, utils.NewNotFoundError("user not found")).Times(1)

		err := uc.DeleteUser(ctx, ids.UserID)

		assert.Error(t, err)
		assert.EqualError(t, err, "user not found")
	})
}

func TestRestoreUser(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := UserUsecaseUtils(t)

	t.Run("should restore user sUCcessfully", func(t *testing.T) {
		repo.User.EXPECT().FindDeletedByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)

		repo.User.EXPECT().Restore(ctx, ids.UserID).Return(nil).Times(1)

		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)

		resp, err := uc.RestoreUser(ctx, ids.UserID)

		assert.NoError(t, err)
		assert.Equal(t, mocks.User.ID, resp.ID)
		assert.Equal(t, mocks.User.Name, resp.Name)
		assert.Equal(t, mocks.User.Email, resp.Email)
	})

	t.Run("should return error not found when find deleted user", func(t *testing.T) {
		repo.User.EXPECT().FindDeletedByID(ctx, ids.UserID).Return(nil, utils.NewNotFoundError("user not found")).Times(1)

		resp, err := uc.RestoreUser(ctx, ids.UserID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "user not found")
	})

	t.Run("should return error internal error when restore user", func(t *testing.T) {
		repo.User.EXPECT().FindDeletedByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)

		repo.User.EXPECT().Restore(ctx, ids.UserID).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.RestoreUser(ctx, ids.UserID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error internal error when find restored user", func(t *testing.T) {
		repo.User.EXPECT().FindDeletedByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)

		repo.User.EXPECT().Restore(ctx, ids.UserID).Return(nil).Times(1)

		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.RestoreUser(ctx, ids.UserID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})
}
