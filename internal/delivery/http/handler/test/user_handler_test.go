package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/usecase/mock"
	"github.com/stretchr/testify/assert"
)

type ResponseUserHandler struct {
	Status  int                 `json:"status"`
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    dto.UserResponseDTO `json:"data"`
	Errors  interface{}         `json:"errors"`
}

type Message struct {
	Message string `json:"message"`
}

type ResponseUserHandlerMessage struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    Message     `json:"data"`
	Errors  interface{} `json:"errors"`
}

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockUserUsecase(ctrl)
	h := handler.NewUserHandler(usecase)
	r := gin.Default()
	r.POST("/users", h.RegisterUser)

	t.Run("Test RegisterUser, successfully", func(t *testing.T) {
		usecase.EXPECT().Register(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
		reqBody := `{"email":"test@example.com"}`
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Test RegisterUser, bind error", func(t *testing.T) {
		usecase.EXPECT().Register(gomock.Any(), gomock.Any()).Times(0)
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(`invalid-json`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetOneUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockUserUsecase(ctrl)
	h := handler.NewUserHandler(usecase)
	r := gin.Default()
	r.GET("/users/:id", h.GetOneUser)

	t.Run("Test GetOneUser, successfully", func(t *testing.T) {
		mockUser := &dto.UserResponseDTO{ID: 1, Email: "test@example.com"}
		usecase.EXPECT().GetUserByID(gomock.Any(), uint64(1)).Return(mockUser, nil).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response ResponseUserHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, mockUser.Email, response.Data.Email)
	})

	t.Run("Test GetOneUser, database error", func(t *testing.T) {
		usecase.EXPECT().GetUserByID(gomock.Any(), uint64(1)).Return(nil, errors.New("internal error")).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Test GetOneUser, invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetAllUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockUserUsecase(ctrl)
	h := handler.NewUserHandler(usecase)
	r := gin.Default()
	r.GET("/users", h.GetAllUsers)

	t.Run("Test GetAllUsers, successfully", func(t *testing.T) {
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-12-05 10:00:00")

		mockUsers := []dto.UserResponseDTO{
			{ID: 1, Email: "test@example.com", Name: "Test", CreatedAt: parsedTime, UpdatedAt: parsedTime},
			{ID: 2, Email: "test2@example.com", Name: "Test2", CreatedAt: parsedTime, UpdatedAt: parsedTime},
		}

		usecase.EXPECT().GetAllUsers(gomock.Any()).Return(&mockUsers, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/users", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response struct {
			Data []dto.UserResponseDTO `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Data, len(mockUsers))
	})

	t.Run("Test GetAllUsers, database error", func(t *testing.T) {
		usecase.EXPECT().GetAllUsers(gomock.Any()).Return(nil, errors.New("internal error")).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/users", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockUserUsecase(ctrl)
	h := handler.NewUserHandler(usecase)
	r := gin.Default()
	r.PATCH("/users/:id", h.UpdateUser)

	t.Run("Test UpdateUser, successfully", func(t *testing.T) {
		mockUser := &dto.UserResponseDTO{ID: 1, Email: "test@example.com"}

		usecase.EXPECT().UpdateUser(gomock.Any(), uint64(1), gomock.Any()).Return(mockUser, nil).Times(1)

		reqBody := `{"name":"updated"}`
		req, _ := http.NewRequest(http.MethodPatch, "/users/1", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response ResponseUserHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, mockUser.Name, response.Data.Name)
	})

	t.Run("Test UpdateUser, database error", func(t *testing.T) {
		usecase.EXPECT().UpdateUser(gomock.Any(), uint64(1), gomock.Any()).Return(nil, errors.New("internal error")).Times(1)

		// Gunakan payload valid
		reqBody := `{"name":"updated"}`
		req, _ := http.NewRequest(http.MethodPatch, "/users/1", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Test UpdateUser, bind error", func(t *testing.T) {
		usecase.EXPECT().UpdateUser(gomock.Any(), uint64(1), gomock.Any()).Times(0)
		req, _ := http.NewRequest(http.MethodPatch, "/users/1", bytes.NewReader([]byte(`invalid-json`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test UpdateUser, invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, "/users/abc", bytes.NewReader([]byte(`invalid-json`)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockUserUsecase(ctrl)
	h := handler.NewUserHandler(usecase)
	r := gin.Default()
	r.DELETE("/users/:id", h.DeleteUser)

	t.Run("Test DeleteUser, successfully", func(t *testing.T) {
		usecase.EXPECT().DeleteUser(gomock.Any(), uint64(1)).Times(1)
		req, _ := http.NewRequest(http.MethodDelete, "/users/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response ResponseUserHandlerMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User deleted successfully", response.Data.Message)
	})

	t.Run("Test DeleteUser, database error", func(t *testing.T) {
		usecase.EXPECT().DeleteUser(gomock.Any(), uint64(1)).Return(errors.New("internal error")).Times(1)
		req, _ := http.NewRequest(http.MethodDelete, "/users/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Test DeleteUser, invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/users/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestRestoreUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockUserUsecase(ctrl)
	h := handler.NewUserHandler(usecase)
	r := gin.Default()
	r.POST("/users/:id/restore", h.RestoreUser)

	t.Run("Test RestoreUser, successfully", func(t *testing.T) {
		usecase.EXPECT().RestoreUser(gomock.Any(), uint64(1)).Return(&dto.UserResponseDTO{ID: 1, Email: "test@example.com"}, nil).Times(1)
		req, _ := http.NewRequest(http.MethodPost, "/users/1/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response ResponseUserHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Data.Email, "test@example.com")
	})

	t.Run("Test RestoreUser, database error", func(t *testing.T) {
		usecase.EXPECT().RestoreUser(gomock.Any(), uint64(1)).Return(nil, errors.New("internal error")).Times(1)
		req, _ := http.NewRequest(http.MethodPost, "/users/1/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Test RestoreUser, invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/users/abc/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
