package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) Register(ctx context.Context, req *dto.UserCreateDTO) (*dto.UserResponseDTO, error) {
	args := m.Called(req)
	return args.Get(0).(*dto.UserResponseDTO), args.Error(1)
}

func (m *MockUserUsecase) GetUserByID(ctx context.Context, id uint64) (*dto.UserResponseDTO, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.UserResponseDTO), args.Error(1)
}

func (m *MockUserUsecase) GetAllUsers(ctx context.Context) (*[]dto.UserResponseDTO, error) {
	args := m.Called()
	return args.Get(0).(*[]dto.UserResponseDTO), args.Error(1)
}

func (m *MockUserUsecase) UpdateUser(ctx context.Context, id uint64, req *dto.UserUpdateDTO) (*dto.UserResponseDTO, error) {
	args := m.Called(req)
	return args.Get(0).(*dto.UserResponseDTO), args.Error(0)
}

func (m *MockUserUsecase) DeleteUser(ctx context.Context, id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserUsecase) RestoreUser(ctx context.Context, id uint64) (*dto.UserResponseDTO, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.UserResponseDTO), args.Error(0)
}

func TestRegisterUserHandler(t *testing.T) {
	mockUsecase := new(MockUserUsecase)
	h := handler.NewUserHandler(mockUsecase)

	router := gin.Default()
	router.POST("/users", h.RegisterUser)

	mockUser := &domain.User{Email: "test@example.com"}
	mockUsecase.On("Register", mockUser).Return(nil)

	reqBody := `{"email":"test@example.com"}`
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockUsecase.AssertCalled(t, "Register", mockUser)
}
