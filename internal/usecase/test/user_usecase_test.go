package usecase_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mock_pkg "github.com/ryvasa/go-super-farmer/pkg/mock"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	mock_repo "github.com/ryvasa/go-super-farmer/internal/repository/mock"
	usecase_implementation "github.com/ryvasa/go-super-farmer/internal/usecase/implementation"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/utils"
	mockUtils "github.com/ryvasa/go-super-farmer/utils/mock"
	"github.com/stretchr/testify/assert"
)

type UserRepoMock struct {
	User  *mock_repo.MockUserRepository
	Hash  *mockUtils.MockHasher
	Cache *mock_pkg.MockCache
}

type UserIDs struct {
	UserID uuid.UUID
}

type UserMocks struct {
	Users                []*domain.User
	User                 *domain.User
	UpdatedUser          *domain.User
	UsersWithoutPassword []*dto.UserResponseDTO
}

type MockUserDTOs struct {
	Create             *dto.UserCreateDTO
	Update             *dto.UserUpdateDTO
	UpdateWithPassword *dto.UserUpdateDTO
}

func UserUsecaseUtils(t *testing.T) (*UserIDs, *UserMocks, *MockUserDTOs, *UserRepoMock, usecase_interface.UserUsecase, context.Context) {
	userID := uuid.New()

	ids := &UserIDs{
		UserID: userID,
	}

	mocks := &UserMocks{
		Users: []*domain.User{
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
		UsersWithoutPassword: []*dto.UserResponseDTO{
			{
				ID:    userID,
				Name:  "test",
				Email: "test@example.com",
			},
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

	userRepo := mock_repo.NewMockUserRepository(ctrl)
	hash := mockUtils.NewMockHasher(ctrl)
	cache := mock_pkg.NewMockCache(ctrl)
	uc := usecase_implementation.NewUserUsecase(userRepo, hash, cache)
	ctx := context.TODO()

	repo := &UserRepoMock{User: userRepo, Hash: hash, Cache: cache}

	return ids, mocks, dto, repo, uc, ctx
}

func TestUserUsecase_Register(t *testing.T) {
	ids, mocks, dtos, repo, uc, ctx := UserUsecaseUtils(t)

	t.Run("should register successfully", func(t *testing.T) {
		repo.User.EXPECT().FindByEmail(ctx, dtos.Create.Email).Return(&domain.User{}, nil).Times(1)
		// Mock hashing password
		repo.Hash.EXPECT().HashPassword(dtos.Create.Password).Return("hashed_password", nil).Times(1)
		// Mock create user
		repo.User.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *domain.User) error {
			user.ID = ids.UserID
			return nil
		}).Times(1)
		// Mock find user by ID
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)
		// Mock delete cache
		repo.Cache.EXPECT().DeleteByPattern(ctx, "user").Return(nil).Times(1)

		// Call the Register method
		resp, err := uc.Register(ctx, dtos.Create)

		// Assert no error
		assert.NoError(t, err)
		// Assert user data matches the input
		assert.Equal(t, dtos.Create.Name, resp.Name)
		assert.Equal(t, dtos.Create.Email, resp.Email)
	})

	t.Run("should return error when validation error", func(t *testing.T) {
		// Call the Register method with invalid data
		resp, err := uc.Register(ctx, &dto.UserCreateDTO{Name: "", Email: "", Password: ""})

		// Assert validation error
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, err.Error(), "Validation failed")
	})

	t.Run("should return error when user email exists", func(t *testing.T) {
		repo.User.EXPECT().FindByEmail(ctx, dtos.Create.Email).Return(mocks.User, nil).Times(1)

		resp, err := uc.Register(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "email already exists")
	})

	t.Run("should return error when hashing password", func(t *testing.T) {
		repo.User.EXPECT().FindByEmail(ctx, dtos.Create.Email).Return(&domain.User{}, nil).Times(1)

		// Mock hashing password failure
		repo.Hash.EXPECT().HashPassword(dtos.Create.Password).Return("", utils.NewInternalError("hashing error")).Times(1)

		// Call the Register method
		resp, err := uc.Register(ctx, dtos.Create)

		// Assert error and check the error message
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "hashing error")
	})

	t.Run("should return error when creating user", func(t *testing.T) {
		repo.User.EXPECT().FindByEmail(ctx, dtos.Create.Email).Return(&domain.User{}, nil).Times(1)

		// Mock hashing password
		repo.Hash.EXPECT().HashPassword(dtos.Create.Password).Return("hashed_password", nil).Times(1)
		// Mock create user failure
		repo.User.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("create user error")).Times(1)

		// Call the Register method
		resp, err := uc.Register(ctx, dtos.Create)

		// Assert error and check the error message
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "create user error")
	})

	t.Run("should return error when finding created user by ID", func(t *testing.T) {
		repo.User.EXPECT().FindByEmail(ctx, dtos.Create.Email).Return(&domain.User{}, nil).Times(1)

		// Mock hashing password
		repo.Hash.EXPECT().HashPassword(dtos.Create.Password).Return("hashed_password", nil).Times(1)
		// Mock create user
		repo.User.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *domain.User) error {
			user.ID = ids.UserID
			return nil
		}).Times(1)
		// Mock find user by ID failure
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(nil, errors.New("find user by id error")).Times(1)

		// Call the Register method
		resp, err := uc.Register(ctx, dtos.Create)

		// Assert error and check the error message
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "find user by id error")
	})

	t.Run("should return error when cache delete fails", func(t *testing.T) {
		repo.User.EXPECT().FindByEmail(ctx, dtos.Create.Email).Return(&domain.User{}, nil).Times(1)

		// Mock hashing password
		repo.Hash.EXPECT().HashPassword(dtos.Create.Password).Return("hashed_password", nil).Times(1)
		// Mock create user
		repo.User.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *domain.User) error {
			user.ID = ids.UserID
			return nil
		}).Times(1)
		// Mock find user by ID
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)
		// Mock cache delete failure
		repo.Cache.EXPECT().DeleteByPattern(ctx, "user").Return(utils.NewInternalError("cache delete error")).Times(1)

		// Call the Register method
		resp, err := uc.Register(ctx, dtos.Create)

		// Assert error and check the error message
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "cache delete error")
	})
}

func TestUserUsecase_GetUserByID(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := UserUsecaseUtils(t)

	t.Run("Test GetUserByID successfully", func(t *testing.T) {

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

func TestUserUsecase_GetAllUsers(t *testing.T) {
	_, mocks, _, repo, uc, ctx := UserUsecaseUtils(t)

	// Mock query params
	queryParams := &dto.PaginationDTO{
		Limit: 10,
		Page:  1,
		Filter: dto.ParamFilterDTO{
			UserName: "test",
		},
	}

	cacheKey := fmt.Sprintf("user_list_page_%d_limit_%d_%s",
		queryParams.Page,
		queryParams.Limit,
		queryParams.Filter.UserName,
	)

	t.Run("should return all users from cache", func(t *testing.T) {
		mockedResponse := &dto.PaginationResponseDTO{
			TotalRows:  10,
			TotalPages: 1,
			Page:       1,
			Limit:      10,
			Data:       mocks.UsersWithoutPassword,
		}
		cached, _ := json.Marshal(mockedResponse)
		repo.Cache.EXPECT().Get(ctx, cacheKey).Return(cached, nil).Times(1)

		resp, err := uc.GetAllUsers(ctx, queryParams)

		assert.NoError(t, err)
		assert.Equal(t, mockedResponse.TotalRows, resp.TotalRows)
		assert.Equal(t, mockedResponse.TotalPages, resp.TotalPages)
		assert.Equal(t, mockedResponse.Page, resp.Page)
		assert.Equal(t, mockedResponse.Limit, resp.Limit)
		assert.Equal(t, mockedResponse.Data, resp.Data)
	})

	t.Run("should return all users from database", func(t *testing.T) {
		repo.Cache.EXPECT().Get(ctx, cacheKey).Return(nil, nil).Times(1)
		repo.User.EXPECT().FindAll(ctx, queryParams).Return(mocks.Users, nil).Times(1)
		repo.User.EXPECT().Count(ctx, &queryParams.Filter).Return(int64(10), nil).Times(1)
		repo.Cache.EXPECT().Set(ctx, cacheKey, gomock.Any(), gomock.Any()).Return(nil).Times(1)

		resp, err := uc.GetAllUsers(ctx, queryParams)

		assert.NoError(t, err)
		assert.Equal(t, len(mocks.Users), len(resp.Data.([]*dto.UserResponseDTO)))
		assert.Equal(t, int64(10), resp.TotalRows)
	})

	t.Run("should return validation error", func(t *testing.T) {
		invalidQueryParams := &dto.PaginationDTO{
			Limit: -1, // Invalid limit
			Page:  1,
		}

		resp, err := uc.GetAllUsers(ctx, invalidQueryParams)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "limit must be greater than 0")
	})

	t.Run("should return error when cache is invalid", func(t *testing.T) {
		repo.Cache.EXPECT().Get(ctx, cacheKey).Return([]byte("invalid data"), nil).Times(1)

		resp, err := uc.GetAllUsers(ctx, queryParams)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("should return error when database call fails", func(t *testing.T) {
		repo.Cache.EXPECT().Get(ctx, cacheKey).Return(nil, nil).Times(1)
		repo.User.EXPECT().FindAll(ctx, queryParams).Return(nil, utils.NewInternalError("db error")).Times(1)

		resp, err := uc.GetAllUsers(ctx, queryParams)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "db error")
	})
}

func TestUserUsecase_UpdateUser(t *testing.T) {
	ids, mocks, dtos, repo, uc, ctx := UserUsecaseUtils(t)

	t.Run("should update user without password", func(t *testing.T) {
		// Mock find user by ID
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)
		// Mock update user
		repo.User.EXPECT().Update(ctx, ids.UserID, gomock.AssignableToTypeOf(&domain.User{})).Return(nil).Times(1)
		// Mock find updated user
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.UpdatedUser, nil).Times(1)
		// Mock cache delete
		repo.Cache.EXPECT().DeleteByPattern(ctx, "user").Return(nil).Times(1)

		resp, err := uc.UpdateUser(ctx, ids.UserID, "Admin", dtos.Update)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, mocks.UpdatedUser.ID, resp.ID)
		assert.Equal(t, dtos.Update.Name, resp.Name)
	})

	t.Run("should update user with password", func(t *testing.T) {
		// Mock find user by ID
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)
		// Mock hash password
		repo.Hash.EXPECT().HashPassword(dtos.UpdateWithPassword.Password).Return("hashed_password", nil).Times(1)
		// Mock update user
		repo.User.EXPECT().Update(ctx, ids.UserID, gomock.AssignableToTypeOf(&domain.User{})).Return(nil).Times(1)
		// Mock find updated user
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.UpdatedUser, nil).Times(1)
		// Mock cache delete
		repo.Cache.EXPECT().DeleteByPattern(ctx, "user").Return(nil).Times(1)

		resp, err := uc.UpdateUser(ctx, ids.UserID, "Admin", dtos.UpdateWithPassword)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, mocks.UpdatedUser.ID, resp.ID)
		assert.Equal(t, dtos.UpdateWithPassword.Name, resp.Name)
	})

	t.Run("should return validation error", func(t *testing.T) {
		resp, err := uc.UpdateUser(ctx, ids.UserID, "Admin", &dto.UserUpdateDTO{Name: "1", Email: "", Password: ""})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("should return error hashing password", func(t *testing.T) {
		// Mock find user by ID
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)
		// Mock hash password error
		repo.Hash.EXPECT().HashPassword(dtos.UpdateWithPassword.Password).Return("", errors.New("hashing error")).Times(1)

		resp, err := uc.UpdateUser(ctx, ids.UserID, "Admin", dtos.UpdateWithPassword)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "hashing error")
	})

	t.Run("should return error user not found", func(t *testing.T) {
		// Mock find user by ID returns not found error
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(nil, utils.NewNotFoundError("user not found")).Times(1)

		resp, err := uc.UpdateUser(ctx, ids.UserID, "Admin", dtos.UpdateWithPassword)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "user not found")
	})

	t.Run("should return error on finding updated user", func(t *testing.T) {
		// Mock find user by ID
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)
		// Mock hash password
		repo.Hash.EXPECT().HashPassword(dtos.UpdateWithPassword.Password).Return("hashed_password", nil).Times(1)
		// Mock update user
		repo.User.EXPECT().Update(ctx, ids.UserID, gomock.AssignableToTypeOf(&domain.User{})).Return(nil).Times(1)
		// Mock find updated user returns internal error
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(nil, utils.NewInternalError("user not found")).Times(1)

		resp, err := uc.UpdateUser(ctx, ids.UserID, "Admin", dtos.UpdateWithPassword)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "user not found")
	})

	t.Run("should return error on cache delete", func(t *testing.T) {
		// Mock find user by ID
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)
		// Mock hash password
		repo.Hash.EXPECT().HashPassword(dtos.UpdateWithPassword.Password).Return("hashed_password", nil).Times(1)
		// Mock update user
		repo.User.EXPECT().Update(ctx, ids.UserID, gomock.AssignableToTypeOf(&domain.User{})).Return(nil).Times(1)
		// Mock find updated user
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.UpdatedUser, nil).Times(1)
		// Mock cache delete error
		repo.Cache.EXPECT().DeleteByPattern(ctx, "user").Return(utils.NewInternalError("cache delete error")).Times(1)

		resp, err := uc.UpdateUser(ctx, ids.UserID, "Admin", dtos.UpdateWithPassword)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "cache delete error")
	})

	t.Run("should return error when forbidden", func(t *testing.T) {
		// Mock find user by ID
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)
		dtos.Update.RoleID = 1
		resp, err := uc.UpdateUser(ctx, ids.UserID, "Farmer", dtos.Update)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "forbidden")
	})
}

func TestUserUsecase_DeleteUser(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := UserUsecaseUtils(t)

	t.Run("should delete user successfully", func(t *testing.T) {
		// Mock find user by ID success
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)
		// Mock delete user success
		repo.User.EXPECT().Delete(ctx, ids.UserID).Return(nil).Times(1)
		// Mock cache delete success
		repo.Cache.EXPECT().DeleteByPattern(ctx, "user").Return(nil).Times(1)

		// Call the DeleteUser method
		err := uc.DeleteUser(ctx, ids.UserID)

		// Assert no error
		assert.NoError(t, err)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		// Mock find user by ID not found
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(nil, utils.NewNotFoundError("user not found")).Times(1)

		// Call the DeleteUser method
		err := uc.DeleteUser(ctx, ids.UserID)

		// Assert error and check the error message
		assert.Error(t, err)
		assert.EqualError(t, err, "user not found")
	})

	t.Run("should return error when delete user fails", func(t *testing.T) {
		// Mock find user by ID success
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)
		// Mock delete user error
		repo.User.EXPECT().Delete(ctx, ids.UserID).Return(utils.NewInternalError("delete user failed")).Times(1)

		// Call the DeleteUser method
		err := uc.DeleteUser(ctx, ids.UserID)

		// Assert error and check the error message
		assert.Error(t, err)
		assert.EqualError(t, err, "delete user failed")
	})

	t.Run("should return error when cache delete fails", func(t *testing.T) {
		// Mock find user by ID success
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)
		// Mock delete user success
		repo.User.EXPECT().Delete(ctx, ids.UserID).Return(nil).Times(1)
		// Mock cache delete error
		repo.Cache.EXPECT().DeleteByPattern(ctx, "user").Return(utils.NewInternalError("cache delete failed")).Times(1)

		// Call the DeleteUser method
		err := uc.DeleteUser(ctx, ids.UserID)

		// Assert error and check the error message
		assert.Error(t, err)
		assert.EqualError(t, err, "cache delete failed")
	})
}

func TestUserUsecase_RestoreUser(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := UserUsecaseUtils(t)

	t.Run("should restore user successfully", func(t *testing.T) {
		// Mock find deleted user success
		repo.User.EXPECT().FindDeletedByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)
		// Mock restore user success
		repo.User.EXPECT().Restore(ctx, ids.UserID).Return(nil).Times(1)
		// Mock find restored user success
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)
		// Mock cache delete success
		repo.Cache.EXPECT().DeleteByPattern(ctx, "user").Return(nil).Times(1)

		// Call the RestoreUser method
		resp, err := uc.RestoreUser(ctx, ids.UserID)

		// Assert no error
		assert.NoError(t, err)
		// Assert returned user data
		assert.Equal(t, mocks.User.ID, resp.ID)
		assert.Equal(t, mocks.User.Name, resp.Name)
		assert.Equal(t, mocks.User.Email, resp.Email)
	})

	t.Run("should return error not found when find deleted user", func(t *testing.T) {
		// Mock find deleted user not found
		repo.User.EXPECT().FindDeletedByID(ctx, ids.UserID).Return(nil, utils.NewNotFoundError("user not found")).Times(1)

		// Call the RestoreUser method
		resp, err := uc.RestoreUser(ctx, ids.UserID)

		// Assert error and check the error message
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "user not found")
	})

	t.Run("should return error internal error when restore user", func(t *testing.T) {
		// Mock find deleted user success
		repo.User.EXPECT().FindDeletedByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)
		// Mock restore user error
		repo.User.EXPECT().Restore(ctx, ids.UserID).Return(utils.NewInternalError("restore failed")).Times(1)

		// Call the RestoreUser method
		resp, err := uc.RestoreUser(ctx, ids.UserID)

		// Assert error and check the error message
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "restore failed")
	})

	t.Run("should return error internal error when find restored user", func(t *testing.T) {
		// Mock find deleted user success
		repo.User.EXPECT().FindDeletedByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)
		// Mock restore user success
		repo.User.EXPECT().Restore(ctx, ids.UserID).Return(nil).Times(1)
		// Mock find restored user error
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(nil, utils.NewInternalError("find restored user failed")).Times(1)

		// Call the RestoreUser method
		resp, err := uc.RestoreUser(ctx, ids.UserID)

		// Assert error and check the error message
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "find restored user failed")
	})

	t.Run("should return error internal error when cache delete fails", func(t *testing.T) {
		// Mock find deleted user success
		repo.User.EXPECT().FindDeletedByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)
		// Mock restore user success
		repo.User.EXPECT().Restore(ctx, ids.UserID).Return(nil).Times(1)
		// Mock find restored user success
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)
		// Mock cache delete error
		repo.Cache.EXPECT().DeleteByPattern(ctx, "user").Return(utils.NewInternalError("cache delete failed")).Times(1)

		// Call the RestoreUser method
		resp, err := uc.RestoreUser(ctx, ids.UserID)

		// Assert error and check the error message
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "cache delete failed")
	})
}
