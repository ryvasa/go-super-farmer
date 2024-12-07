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
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/usecase/mock"
	"github.com/stretchr/testify/assert"
)

type ResponseAuthHandler struct {
	Status  int                 `json:"status"`
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    dto.AuthResponseDTO `json:"data"`
	Errors  interface{}         `json:"errors"`
}

func TestLoginAuth(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockAuthUsecase(ctrl)
	h := handler.NewAuthHandler(usecase)
	r := gin.Default()
	r.POST("/auth/login", h.Login)

	t.Run("Test Login, successfully", func(t *testing.T) {
		userID := uuid.New()
		user := &dto.UserResponseDTO{ID: userID, Email: "test@example.com"}
		mockResponse := &dto.AuthResponseDTO{User: user, Token: "token"}
		usecase.EXPECT().Login(gomock.Any(), gomock.Any()).Return(mockResponse, nil).Times(1)
		reqBody := `{"email":"test@example.com"}`
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response ResponseAuthHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Data.User.Email, "test@example.com")
	})

	t.Run("Test Login, bind error", func(t *testing.T) {
		usecase.EXPECT().Login(gomock.Any(), gomock.Any()).Times(0)
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte(`invalid-json`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
