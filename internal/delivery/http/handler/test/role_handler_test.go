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
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/test/response"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/usecase/mock"
	"github.com/ryvasa/go-super-farmer/utils"
)

type responseRoleHandler struct {
	Status  int            `json:"status"`
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    domain.Role    `json:"data"`
	Errors  response.Error `json:"errors"`
}

type responseRolesHandler struct {
	Status  int            `json:"status"`
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    []domain.Role  `json:"data"`
	Errors  response.Error `json:"errors"`
}

type RoleHandlerMocks struct {
	Role  *domain.Role
	Roles *[]domain.Role
}

type RoleHandlerIDs struct {
	RoleID int64
}

func RoleHandlerSetup(t *testing.T) (*gin.Engine, handler.RoleHandler, *mock.MockRoleUsecase, RoleHandlerIDs, RoleHandlerMocks) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mock.NewMockRoleUsecase(ctrl)
	r := gin.Default()
	h := handler.NewRoleHandler(uc)

	roleID := int64(1)
	ids := RoleHandlerIDs{
		RoleID: roleID,
	}

	mocks := RoleHandlerMocks{
		Role: &domain.Role{
			ID:   roleID,
			Name: "Admin",
		},
		Roles: &[]domain.Role{
			{
				ID:   roleID,
				Name: "Admin",
			},
		},
	}

	return r, h, uc, ids, mocks

}

func TestCreateRole(t *testing.T) {
	r, h, uc, _, mocks := RoleHandlerSetup(t)

	r.POST("/roles", h.CreateRole)

	t.Run("should create role successfully", func(t *testing.T) {

		uc.EXPECT().CreateRole(gomock.Any(), gomock.Any()).Return(mocks.Role, nil).Times(1)

		reqBody := `{"name": "Admin"}`
		req, _ := http.NewRequest(http.MethodPost, "/roles", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		var response responseRoleHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, mocks.Role.Name, response.Data.Name)
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		uc.EXPECT().CreateRole(gomock.Any(), gomock.Any()).Times(0)

		req, _ := http.NewRequest(http.MethodPost, "/roles", bytes.NewReader([]byte(`invalid-json`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		var response responseRoleHandler
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.EXPECT().CreateRole(gomock.Any(), gomock.Any()).Return(nil, utils.NewInternalError("internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodPost, "/roles", bytes.NewReader([]byte(`{"name": "Admin"}`)))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		var response responseRoleHandler
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetAllRoles(t *testing.T) {
	r, h, uc, _, mocks := RoleHandlerSetup(t)

	r.GET("/roles", h.GetAllRoles)

	t.Run("should return all roles successfully", func(t *testing.T) {
		uc.EXPECT().GetAllRoles(gomock.Any()).Return(mocks.Roles, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/roles", nil)
		w := httptest.NewRecorder()

		// Execute request
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		// Adjust unmarshaling based on the actual response structure
		var response responseRolesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, len(*mocks.Roles), len(response.Data))
		assert.Equal(t, (*mocks.Roles)[0].Name, response.Data[0].Name)
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.EXPECT().GetAllRoles(gomock.Any()).Return(nil, utils.NewInternalError("internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/roles", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		var response responseRolesHandler
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
