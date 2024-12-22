package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	handler_implementation "github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler/implementation"
	handler_interface "github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler/interface"
	"github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler/test/response"
	"github.com/ryvasa/go-super-farmer/service_api/model/dto"
	mock_usecase "github.com/ryvasa/go-super-farmer/service_api/usecase/mock"
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
)

type responseAuthHandler struct {
	Status  int                 `json:"status"`
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    dto.AuthResponseDTO `json:"data"`
	Errors  response.Error      `json:"errors"`
}

type AuthHandlerMocks struct {
	AuthResponse *dto.AuthResponseDTO
}
type AuthHandlerIDs struct {
	UserID uuid.UUID
}

func AuthHandlerSetUp(t *testing.T) (*gin.Engine, handler_interface.AuthHandler, *mock_usecase.MockAuthUsecase, AuthHandlerIDs, AuthHandlerMocks) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc := mock_usecase.NewMockAuthUsecase(ctrl)
	h := handler_implementation.NewAuthHandler(uc)
	r := gin.Default()

	mocks := AuthHandlerMocks{
		AuthResponse: &dto.AuthResponseDTO{
			Token: "token",
			User: &dto.UserResponseDTO{
				ID:    uuid.New(),
				Email: "test@example.com",
				Name:  "name",
			},
		},
	}
	ids := AuthHandlerIDs{
		UserID: uuid.New(),
	}

	return r, h, uc, ids, mocks
}

func TestAuthHandler_Login(t *testing.T) {
	r, h, uc, _, mocks := AuthHandlerSetUp(t)

	r.POST("/auth/login", h.Login)

	t.Run("should login successfully", func(t *testing.T) {
		uc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(mocks.AuthResponse, nil).Times(1)

		reqBody := `{"email":"test@example.com"}`
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response responseAuthHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Data.User.Email, "test@example.com")
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		uc.EXPECT().Login(gomock.Any(), gomock.Any()).Times(0)
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte(`invalid-json`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseAuthHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when service_api error", func(t *testing.T) {
		uc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(nil, utils.NewInternalError("Internal error"))

		reqBody := `{"email":"test@example.com"}`
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseAuthHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAuthHandler_SendOTP(t *testing.T) {
	r, h, uc, _, _ := AuthHandlerSetUp(t)
	r.POST("/auth/send", h.SendOTP)

	t.Run("should send otp successfully", func(t *testing.T) {
		uc.EXPECT().SendOTP(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		reqBody := `{"email":"test@example.com"}`
		req, _ := http.NewRequest(http.MethodPost, "/auth/send", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response responseAuthHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		uc.EXPECT().SendOTP(gomock.Any(), gomock.Any()).Times(0)

		req, _ := http.NewRequest(http.MethodPost, "/auth/send", bytes.NewReader([]byte(`invalid-json`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseAuthHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, "BAD_REQUEST", response.Errors.Code)
	})

	t.Run("should return error when service_api error", func(t *testing.T) {
		uc.EXPECT().SendOTP(gomock.Any(), gomock.Any()).Return(utils.NewInternalError("Internal error"))

		reqBody := `{"email":"test@example.com"}`
		req, _ := http.NewRequest(http.MethodPost, "/auth/send", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseAuthHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, "INTERNAL_ERROR", response.Errors.Code)
	})
}

func TestAuthHandler_VerifyOTP(t *testing.T) {
	r, h, uc, _, _ := AuthHandlerSetUp(t)
	r.POST("/auth/verify", h.VerifyOTP)

	t.Run("should verify otp successfully", func(t *testing.T) {
		uc.EXPECT().VerifyOTP(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		reqBody := `{"email":"test@example.com","otp":"123456"}`
		req, _ := http.NewRequest(http.MethodPost, "/auth/verify", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response responseAuthHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		uc.EXPECT().VerifyOTP(gomock.Any(), gomock.Any()).Times(0)

		req, _ := http.NewRequest(http.MethodPost, "/auth/verify", bytes.NewReader([]byte(`invalid-json`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseAuthHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, "BAD_REQUEST", response.Errors.Code)
	})

	t.Run("should return error when invalid otp", func(t *testing.T) {
		uc.EXPECT().VerifyOTP(gomock.Any(), gomock.Any()).Return(utils.NewBadRequestError("Invalid OTP"))

		reqBody := `{"email":"test@example.com","otp":"123456"}`
		req, _ := http.NewRequest(http.MethodPost, "/auth/verify", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseAuthHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, "BAD_REQUEST", response.Errors.Code)
	})
}
