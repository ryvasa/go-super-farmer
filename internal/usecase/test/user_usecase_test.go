package usecase_test

import (
	"context"
	"testing"

	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint64) (*domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindAll(ctx context.Context) (*[]domain.User, error) {
	args := m.Called()
	return args.Get(0).(*[]domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, id uint64, user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) Restore(ctx context.Context, id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) FindDeletedByID(ctx context.Context, id uint64) (*domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestRegisterUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := usecase.NewUserUsecase(mockRepo)

	mockUser := &dto.UserCreateDTO{Name: "John Doe", Email: "test@example.com"}
	mockRepo.On("Create", mockUser).Return(nil)

	var ctx context.Context

	createdUser, err := uc.Register(ctx, mockUser)

	assert.NoError(t, err)
	assert.Equal(t, mockUser, createdUser)
	mockRepo.AssertCalled(t, "Create", mockUser)
}

func TestGetUserByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := usecase.NewUserUsecase(mockRepo)

	userID := uint64(1)
	mockUser := &domain.User{ID: userID, Email: "test@example.com"}
	mockRepo.On("FindByID", userID).Return(mockUser, nil)

	var ctx context.Context

	user, err := uc.GetUserByID(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, mockUser, user)
}
