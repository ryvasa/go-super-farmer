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
	Users []*dto.UserResponseDTO
}

type UserHandlerIDs struct {
	UserID uuid.UUID
	RoleID int64
}
type UserUCMOcks struct {
	Auth *mock_usecase.MockAuthUsecase
	User *mock_usecase.MockUserUsecase
}

func UserHandlerSetup(t *testing.T) (*gin.Engine, handler_interface.UserHandler, *UserUCMOcks, UserHandlerIDs, UserHandlerMocks) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ucUser := mock_usecase.NewMockUserUsecase(ctrl)
	ucAut := mock_usecase.NewMockAuthUsecase(ctrl)
	h := handler_implementation.NewUserHandler(ucUser, ucAut)
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
		Users: []*dto.UserResponseDTO{
			{
				ID:    ids.UserID,
				Name:  "Test",
				Email: "test@example.com",
			},
		},
	}
	uc := &UserUCMOcks{
		Auth: ucAut,
		User: ucUser,
	}
	return r, h, uc, ids, mocks

}

func TestUserHandler_RegisterUser(t *testing.T) {
	r, h, uc, _, mocks := UserHandlerSetup(t)
	r.POST("/users", h.RegisterUser)

	t.Run("should register user successfully", func(t *testing.T) {
		uc.User.EXPECT().Register(gomock.Any(), gomock.Any()).Return(mocks.User, nil).Times(1)
		uc.Auth.EXPECT().SendOTP(gomock.Any(), gomock.Any()).Return(nil).Times(1)

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
		uc.User.EXPECT().Register(gomock.Any(), gomock.Any()).Times(0)
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

	t.Run("should return error when service_api error", func(t *testing.T) {
		uc.User.EXPECT().Register(gomock.Any(), gomock.Any()).Return(nil, utils.NewInternalError("service_api error")).Times(1)
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

	t.Run("should return error when send otp error", func(t *testing.T) {
		uc.User.EXPECT().Register(gomock.Any(), gomock.Any()).Return(mocks.User, nil).Times(1)
		uc.Auth.EXPECT().SendOTP(gomock.Any(), gomock.Any()).Return(utils.NewInternalError("service_api error")).Times(1)
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
		uc.User.EXPECT().GetUserByID(gomock.Any(), ids.UserID).Return(mocks.User, nil).Times(1)

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

		uc.User.EXPECT().GetUserByID(gomock.Any(), ids.UserID).Return(nil, utils.NewInternalError("service_api error")).Times(1)
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

	t.Run("should return all users successfully with default pagination", func(t *testing.T) {
		// Prepare expected response
		expectedResponse := &dto.PaginationResponseDTO{
			TotalRows:  1,
			TotalPages: 1,
			Page:       1,
			Limit:      10,
			Data:       mocks.Users,
		}

		// Setup mock
		uc.User.EXPECT().
			GetAllUsers(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx *gin.Context, p *dto.PaginationDTO) (*dto.PaginationResponseDTO, error) {
				// Verify default pagination
				assert.Equal(t, 1, p.Page)
				assert.Equal(t, 10, p.Limit)
				assert.Equal(t, "created_at desc", p.Sort)
				return expectedResponse, nil
			})

		// Make request
		req, _ := http.NewRequest(http.MethodGet, "/users", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)
		var response struct {
			Success bool                       `json:"success"`
			Data    *dto.PaginationResponseDTO `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, expectedResponse.TotalRows, response.Data.TotalRows)
		assert.Equal(t, expectedResponse.TotalPages, response.Data.TotalPages)
		assert.Equal(t, expectedResponse.Page, response.Data.Page)
		assert.Equal(t, expectedResponse.Limit, response.Data.Limit)
	})

	t.Run("should return users with custom pagination and filter", func(t *testing.T) {
		expectedResponse := &dto.PaginationResponseDTO{
			TotalRows:  1,
			TotalPages: 1,
			Page:       2,
			Limit:      5,
			Data:       mocks.Users,
		}

		// Setup mock
		uc.User.EXPECT().
			GetAllUsers(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx *gin.Context, p *dto.PaginationDTO) (*dto.PaginationResponseDTO, error) {
				assert.Equal(t, 2, p.Page)
				assert.Equal(t, 5, p.Limit)
				assert.Equal(t, "test", p.Filter.UserName)
				return expectedResponse, nil
			})

		// Make request with query params
		req, _ := http.NewRequest(http.MethodGet, "/users?page=2&limit=5&user_name=test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response struct {
			Success bool                       `json:"success"`
			Data    *dto.PaginationResponseDTO `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, expectedResponse.TotalRows, response.Data.TotalRows)
	})

	t.Run("should return error with invalid pagination params", func(t *testing.T) {
		// Make request with invalid page
		req, _ := http.NewRequest(http.MethodGet, "/users?page=-1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response struct {
			Success bool           `json:"success"`
			Message string         `json:"message"`
			Errors  response.Error `json:"errors"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "BAD_REQUEST", response.Errors.Code)
		assert.Contains(t, response.Errors.Message, "page must be greater than 0")
	})

	t.Run("should return error with too large limit", func(t *testing.T) {
		// Make request with invalid limit
		req, _ := http.NewRequest(http.MethodGet, "/users?limit=101", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response struct {
			Success bool           `json:"success"`
			Message string         `json:"message"`
			Errors  response.Error `json:"errors"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "BAD_REQUEST", response.Errors.Code)
		assert.Contains(t, response.Errors.Message, "limit must not exceed 100")
	})

	t.Run("should return error when usecase returns error", func(t *testing.T) {
		uc.User.EXPECT().
			GetAllUsers(gomock.Any(), gomock.Any()).
			Return(nil, utils.NewInternalError("service_api error"))

		req, _ := http.NewRequest(http.MethodGet, "/users", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response struct {
			Success bool        `json:"success"`
			Errors  interface{} `json:"errors"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUserHandler_UpdateUser(t *testing.T) {
	r, h, uc, ids, mocks := UserHandlerSetup(t)
	r.PATCH("/users/:id", h.UpdateUser)

	t.Run("should update user successfully", func(t *testing.T) {
		uc.User.EXPECT().UpdateUser(gomock.Any(), ids.UserID, gomock.Any()).Return(mocks.User, nil).Times(1)

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
		uc.User.EXPECT().UpdateUser(gomock.Any(), ids.UserID, gomock.Any()).Return(nil, utils.NewInternalError("service_api error")).Times(1)

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
		uc.User.EXPECT().UpdateUser(gomock.Any(), ids.UserID, gomock.Any()).Times(0)
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
		uc.User.EXPECT().DeleteUser(gomock.Any(), ids.UserID).Times(1)
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
		uc.User.EXPECT().DeleteUser(gomock.Any(), ids.UserID).Return(utils.NewInternalError("service_api error")).Times(1)
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
		uc.User.EXPECT().RestoreUser(gomock.Any(), ids.UserID).Return(&dto.UserResponseDTO{ID: ids.UserID, Email: "test@example.com"}, nil).Times(1)
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

		uc.User.EXPECT().RestoreUser(gomock.Any(), ids.UserID).Return(nil, utils.NewInternalError("service_api error")).Times(1)
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
