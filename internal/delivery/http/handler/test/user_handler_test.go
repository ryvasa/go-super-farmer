package handler_test

import (
	"bytes"
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

func (m *MockUserUsecase) Register(req *dto.UserCreateDTO) (*dto.UserResponseDTO, error) {
	args := m.Called(req)
	return args.Get(0).(*dto.UserResponseDTO), args.Error(1)
}

func (m *MockUserUsecase) GetUserByID(id int64) (*dto.UserResponseDTO, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.UserResponseDTO), args.Error(1)
}

func (m *MockUserUsecase) GetAllUsers() (*[]dto.UserResponseDTO, error) {
	args := m.Called()
	return args.Get(0).(*[]dto.UserResponseDTO), args.Error(1)
}

func (m *MockUserUsecase) UpdateUser(id int64, req *dto.UserUpdateDTO) (*dto.UserResponseDTO, error) {
	args := m.Called(req)
	return args.Get(0).(*dto.UserResponseDTO), args.Error(0)
}

func (m *MockUserUsecase) DeleteUser(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserUsecase) RestoreUser(id int64) (*dto.UserResponseDTO, error) {
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
