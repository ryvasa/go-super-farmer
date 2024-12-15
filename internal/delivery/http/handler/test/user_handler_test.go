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
	handler_implementation "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/implementation"
	handler_interface "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/interface"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/test/response"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/usecase/mock"
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
)

type responseUserHandler struct {
	Status  int                 `json:"status"`
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    dto.UserResponseDTO `json:"data"`
	Errors  response.Error      `json:"errors"`
}

type responseUsersHandler struct {
	Status  int                   `json:"status"`
	Success bool                  `json:"success"`
	Message string                `json:"message"`
	Data    []dto.UserResponseDTO `json:"data"`
	Errors  response.Error        `json:"errors"`
}

type UserHandlerMocks struct {
	User  *dto.UserResponseDTO
	Users *[]dto.UserResponseDTO
}

type UserHandlerIDs struct {
	UserID uuid.UUID
	RoleID int64
}

func UserHandlerSetup(t *testing.T) (*gin.Engine, handler_interface.UserHandler, *mock.MockUserUsecase, UserHandlerIDs, UserHandlerMocks) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc := mock.NewMockUserUsecase(ctrl)
	h := handler_implementation.NewUserHandler(uc)
	r := gin.Default()

	ids := UserHandlerIDs{
		UserID: uuid.New(),
		RoleID: 1,
	}
	mocks := UserHandlerMocks{
		User: &dto.UserResponseDTO{
			ID:    ids.UserID,
			Name:  "Test",
			Email: "test@example.com",
		},
		Users: &[]dto.UserResponseDTO{
			{
				ID:    ids.UserID,
				Name:  "Test",
				Email: "test@example.com",
			},
		},
	}

	return r, h, uc, ids, mocks

}

func TestUserHandler_RegisterUser(t *testing.T) {
	r, h, uc, _, mocks := UserHandlerSetup(t)
	r.POST("/users", h.RegisterUser)

	t.Run("should register user successfully", func(t *testing.T) {
		uc.EXPECT().Register(gomock.Any(), gomock.Any()).Return(mocks.User, nil).Times(1)

		reqBody := `{"email":"test@example.com"}`
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseUserHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Data.Email, "test@example.com")
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		uc.EXPECT().Register(gomock.Any(), gomock.Any()).Times(0)
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(`invalid-json`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseUserHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().Register(gomock.Any(), gomock.Any()).Return(nil, utils.NewInternalError("internal error")).Times(1)
		reqBody := `{"email":"test@example.com"}`
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseUserHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUserHandler_GetOneUser(t *testing.T) {
	r, h, uc, ids, mocks := UserHandlerSetup(t)
	r.GET("/users/:id", h.GetOneUser)

	t.Run("should return user by id successfully", func(t *testing.T) {
		uc.EXPECT().GetUserByID(gomock.Any(), ids.UserID).Return(mocks.User, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/users/"+ids.UserID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response responseUserHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, mocks.User.ID, response.Data.ID)
		assert.Equal(t, mocks.User.Email, response.Data.Email)
	})

	t.Run("should return error when usecase error", func(t *testing.T) {

		uc.EXPECT().GetUserByID(gomock.Any(), ids.UserID).Return(nil, utils.NewInternalError("internal error")).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/users/"+ids.UserID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseUserHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseUserHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserHandler_GetAllUsers(t *testing.T) {
	r, h, uc, _, mocks := UserHandlerSetup(t)
	r.GET("/users", h.GetAllUsers)

	t.Run("should return all users successfully", func(t *testing.T) {
		uc.EXPECT().GetAllUsers(gomock.Any()).Return(mocks.Users, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/users", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseUsersHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Len(t, response.Data, len(*mocks.Users))
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.EXPECT().GetAllUsers(gomock.Any()).Return(nil, utils.NewInternalError("internal error")).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/users", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseUsersHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUserHandler_UpdateUser(t *testing.T) {
	r, h, uc, ids, mocks := UserHandlerSetup(t)
	r.PATCH("/users/:id", h.UpdateUser)

	t.Run("should update user successfully", func(t *testing.T) {
		uc.EXPECT().UpdateUser(gomock.Any(), ids.UserID, gomock.Any()).Return(mocks.User, nil).Times(1)

		reqBody := `{"name":"updated"}`
		req, _ := http.NewRequest(http.MethodPatch, "/users/"+ids.UserID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		var response responseUserHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.NoError(t, err)
		assert.Equal(t, mocks.User.Name, response.Data.Name)
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.EXPECT().UpdateUser(gomock.Any(), ids.UserID, gomock.Any()).Return(nil, utils.NewInternalError("internal error")).Times(1)

		// Gunakan payload valid
		reqBody := `{"name":"updated"}`
		req, _ := http.NewRequest(http.MethodPatch, "/users/"+ids.UserID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseUserHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		uc.EXPECT().UpdateUser(gomock.Any(), ids.UserID, gomock.Any()).Times(0)
		req, _ := http.NewRequest(http.MethodPatch, "/users/"+ids.UserID.String(), bytes.NewReader([]byte(`invalid-json`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseUserHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, "/users/abc", bytes.NewReader([]byte(`invalid-json`)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseUserHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

}

func TestUserHandler_DeleteUser(t *testing.T) {
	r, h, uc, ids, _ := UserHandlerSetup(t)
	r.DELETE("/users/:id", h.DeleteUser)

	t.Run("should delete user successfully", func(t *testing.T) {
		uc.EXPECT().DeleteUser(gomock.Any(), ids.UserID).Times(1)
		req, _ := http.NewRequest(http.MethodDelete, "/users/"+ids.UserID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User deleted successfully", response.Data.Message)
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.EXPECT().DeleteUser(gomock.Any(), ids.UserID).Return(utils.NewInternalError("internal error")).Times(1)
		req, _ := http.NewRequest(http.MethodDelete, "/users/"+ids.UserID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/users/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
	})
}

func TestUserHandler_RestoreUser(t *testing.T) {
	r, h, uc, ids, _ := UserHandlerSetup(t)
	r.POST("/users/:id/restore", h.RestoreUser)

	t.Run("should restore user successfully", func(t *testing.T) {
		uc.EXPECT().RestoreUser(gomock.Any(), ids.UserID).Return(&dto.UserResponseDTO{ID: ids.UserID, Email: "test@example.com"}, nil).Times(1)
		req, _ := http.NewRequest(http.MethodPost, "/users/"+ids.UserID.String()+"/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response responseUserHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Data.Email, "test@example.com")
	})

	t.Run("should return error when usecase error", func(t *testing.T) {

		uc.EXPECT().RestoreUser(gomock.Any(), ids.UserID).Return(nil, utils.NewInternalError("internal error")).Times(1)
		req, _ := http.NewRequest(http.MethodPost, "/users/"+ids.UserID.String()+"/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response responseUserHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/users/abc/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseUserHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
