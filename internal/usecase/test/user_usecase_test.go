package usecase_test

import (
	"testing"

	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindById(id int64) (*domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindAll() ([]domain.User, error) {
	args := m.Called()
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(id int64, user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) Restore(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestRegisterUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := usecase.NewUserUsecase(mockRepo)

	mockUser := &domain.User{Email: "test@example.com"}
	mockRepo.On("Create", mockUser).Return(nil)

	err := uc.Register(mockUser)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "Create", mockUser)
}

func TestGetUserByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := usecase.NewUserUsecase(mockRepo)

	userID := int64(1)
	mockUser := &domain.User{ID: userID, Email: "test@example.com"}
	mockRepo.On("FindById", userID).Return(mockUser, nil)

	user, err := uc.GetUserByID(userID)

	assert.NoError(t, err)
	assert.Equal(t, mockUser, user)
}
