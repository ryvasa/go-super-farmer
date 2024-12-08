package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/usecase/mock"
)

type responseRoleHandler struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    domain.Role `json:"data"`
	Errors  interface{} `json:"errors"`
}

func TestCreateRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := mock.NewMockRoleUsecase(ctrl)
	r := gin.Default()
	h := handler.NewRoleHandler(usecase)

	r.POST("/roles", h.CreateRole)

	t.Run("Test CreateRole, successfully", func(t *testing.T) {
		mockRole := domain.Role{
			ID:   1,
			Name: "Admin",
		}

		usecase.EXPECT().CreateRole(gomock.Any(), gomock.Any()).Return(&mockRole, nil).Times(1)

		reqBody := `{"name": "Admin"}`
		req, _ := http.NewRequest(http.MethodPost, "/roles", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Decode response
		var response responseRoleHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Validate response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, mockRole.Name, response.Data.Name) // Cocokkan dengan data di dalam field "Data"
	})

	t.Run("Test CreateRole, bind error", func(t *testing.T) {
		usecase.EXPECT().CreateRole(gomock.Any(), gomock.Any()).Times(0)

		req, _ := http.NewRequest(http.MethodPost, "/roles", bytes.NewReader([]byte(`invalid-json`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetAllRoles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := mock.NewMockRoleUsecase(ctrl)
	r := gin.Default()
	h := handler.NewRoleHandler(usecase)

	r.GET("/roles", h.GetAllRoles)

	t.Run("Test GetAllRoles, successfully", func(t *testing.T) {
		mockResponse := []domain.Role{
			{ID: 1, Name: "Admin"},
			{ID: 2, Name: "User"},
		}

		// Mock behavior
		usecase.EXPECT().GetAllRoles(gomock.Any()).Return(&mockResponse, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/roles", nil)
		w := httptest.NewRecorder()

		// Execute request
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		// Adjust unmarshaling based on the actual response structure
		var response struct {
			Data []domain.Role `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, len(mockResponse), len(response.Data))
		assert.Equal(t, mockResponse[0].Name, response.Data[0].Name)
	})

	t.Run("Test GetAllRoles, usecase error", func(t *testing.T) {
		usecase.EXPECT().GetAllRoles(gomock.Any()).Return(nil, assert.AnError).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/roles", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
