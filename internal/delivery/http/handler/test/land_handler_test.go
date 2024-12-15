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
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/usecase/mock"
	"github.com/ryvasa/go-super-farmer/utils"
	mockAuthUtil "github.com/ryvasa/go-super-farmer/utils/mock"
	"github.com/stretchr/testify/assert"
)

type responseLandHandler struct {
	Status  int            `json:"status"`
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    domain.Land    `json:"data"`
	Errors  response.Error `json:"errors"`
}

type responseLandsHandler struct {
	Status  int            `json:"status"`
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    []domain.Land  `json:"data"`
	Errors  response.Error `json:"errors"`
}

type LandHandlerMocks struct {
	Land  *domain.Land
	Lands *[]domain.Land
}

type LandHandlerIDs struct {
	LandID uuid.UUID
	UserID uuid.UUID
}

func LandHandlerSetup(t *testing.T) (*gin.Engine, handler_interface.LandHandler, *mock.MockLandUsecase, *mockAuthUtil.MockAuthUtil, LandHandlerIDs, LandHandlerMocks) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc := mock.NewMockLandUsecase(ctrl)
	utils := mockAuthUtil.NewMockAuthUtil(ctrl)
	h := handler_implementation.NewLandHandler(uc, utils)
	r := gin.Default()

	landID := uuid.New()
	userID := uuid.New()
	ids := LandHandlerIDs{
		LandID: landID,
		UserID: userID,
	}

	mocks := LandHandlerMocks{
		Land: &domain.Land{
			ID:          landID,
			UserID:      ids.UserID,
			LandArea:    float64(1000),
			Certificate: "certificate",
		},
		Lands: &[]domain.Land{
			{
				ID:          landID,
				UserID:      ids.UserID,
				LandArea:    float64(1000),
				Certificate: "certificate",
			},
		},
	}

	return r, h, uc, utils, ids, mocks
}

func TestLandHandler_CreateLand(t *testing.T) {
	r, h, uc, authUtil, ids, mocks := LandHandlerSetup(t)
	r.POST("/lands", h.CreateLand)

	t.Run("should create land successfully", func(t *testing.T) {
		authUtil.EXPECT().GetAuthUserID(gomock.Any()).Return(ids.UserID, nil).Times(1)
		uc.EXPECT().CreateLand(gomock.Any(), ids.UserID, gomock.Any()).Return(mocks.Land, nil).Times(1)

		reqBody := `{"land_area":10,"certificate":"test"}`
		req, _ := http.NewRequest(http.MethodPost, "/lands", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Data.ID, ids.LandID)
		assert.Equal(t, response.Data.Certificate, "certificate")
		assert.Equal(t, response.Data.LandArea, float64(1000))
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("should return error when authorization error", func(t *testing.T) {
		authUtil.EXPECT().GetAuthUserID(gomock.Any()).Return(uuid.UUID{}, utils.NewUnauthorizedError("unauthorized")).Times(1)

		reqBody := `{"land_area":10,"certificate":"test"}`
		req, _ := http.NewRequest(http.MethodPost, "/lands", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "UNAUTHORIZED")
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		authUtil.EXPECT().GetAuthUserID(gomock.Any()).Return(ids.UserID, nil).Times(1)

		reqBody := `{"land_area":"invalid","certificate":"test"}`
		req, _ := http.NewRequest(http.MethodPost, "/lands", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		authUtil.EXPECT().GetAuthUserID(gomock.Any()).Return(ids.UserID, nil).Times(1)
		uc.EXPECT().CreateLand(gomock.Any(), ids.UserID, gomock.Any()).Return(nil, utils.NewInternalError("internal server error")).Times(1)

		reqBody := `{"land_area":10,"certificate":"test"}`
		req, _ := http.NewRequest(http.MethodPost, "/lands", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestLandHandler_GetOneLand(t *testing.T) {
	r, h, uc, _, ids, mocks := LandHandlerSetup(t)

	r.GET("/lands/:id", h.GetLandByID)
	t.Run("should return land by id successfully", func(t *testing.T) {
		uc.EXPECT().GetLandByID(gomock.Any(), ids.LandID).Return(mocks.Land, nil).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/lands/"+ids.LandID.String(), nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Data.ID, ids.LandID)
		assert.Equal(t, response.Data.Certificate, "certificate")
		assert.Equal(t, response.Data.LandArea, float64(1000))
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.EXPECT().GetLandByID(gomock.Any(), ids.LandID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/lands/"+ids.LandID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return bad request when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/lands/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLandHandler_GetOneLandByUserID(t *testing.T) {
	r, h, uc, _, ids, mocks := LandHandlerSetup(t)

	r.GET("/lands/user/:id", h.GetLandByUserID)

	t.Run("should return all lands when successfully", func(t *testing.T) {
		uc.EXPECT().GetLandByUserID(gomock.Any(), ids.UserID).Return(mocks.Lands, nil).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/lands/user/"+ids.UserID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandsHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Data, len(*mocks.Lands))
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.EXPECT().GetLandByUserID(gomock.Any(), ids.UserID).Return(nil, utils.NewInternalError("internal error")).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/lands/user/"+ids.UserID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandsHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return bad request when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/lands/user/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandsHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLandHandler_GetAllLands(t *testing.T) {
	r, h, uc, _, _, mocks := LandHandlerSetup(t)

	r.GET("/lands", h.GetAllLands)

	t.Run("should return all lands successfully", func(t *testing.T) {

		uc.EXPECT().GetAllLands(gomock.Any()).Return(mocks.Lands, nil).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/lands", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandsHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Data, len(*mocks.Lands))
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.EXPECT().GetAllLands(gomock.Any()).Return(nil, utils.NewInternalError("internal error")).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/lands", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandsHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestLandHandler_UpdateLand(t *testing.T) {
	r, h, uc, authUtil, ids, mocks := LandHandlerSetup(t)

	r.PATCH("/lands/:id", h.UpdateLand)

	t.Run("should update land successfully", func(t *testing.T) {
		authUtil.EXPECT().GetAuthUserID(gomock.Any()).Return(ids.UserID, nil).Times(1)
		uc.EXPECT().UpdateLand(gomock.Any(), ids.UserID, gomock.Any(), gomock.Any()).Return(mocks.Land, nil).Times(1)

		reqBody := `{"land_area":10,"certificate":"test"}`
		req, _ := http.NewRequest(http.MethodPatch, "/lands/"+ids.LandID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response responseLandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, mocks.Land.LandArea, response.Data.LandArea)
		assert.Equal(t, mocks.Land.Certificate, response.Data.Certificate)
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		authUtil.EXPECT().GetAuthUserID(gomock.Any()).Return(ids.UserID, nil).Times(1)
		uc.EXPECT().UpdateLand(gomock.Any(), ids.UserID, ids.LandID, gomock.Any()).Return(nil, utils.NewInternalError("internal server error")).Times(1)

		reqBody := `{"land_area":10,"certificate":"test"}`
		req, _ := http.NewRequest(http.MethodPatch, "/lands/"+ids.LandID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, "/lands/abc", bytes.NewReader([]byte(`invalid-json`)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		authUtil.EXPECT().GetAuthUserID(gomock.Any()).Return(ids.UserID, nil).Times(1)
		uc.EXPECT().UpdateLand(gomock.Any(), ids.UserID, ids.LandID, gomock.Any()).Return(nil, utils.NewNotFoundError("land not found")).Times(1)

		reqBody := `{"land_area":10,"certificate":"test"}`
		req, _ := http.NewRequest(http.MethodPatch, "/lands/"+ids.LandID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "NOT_FOUND")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestLandHandler_DeleteLand(t *testing.T) {
	r, h, uc, _, ids, _ := LandHandlerSetup(t)

	r.DELETE("/lands/:id", h.DeleteLand)

	t.Run("should delete land successfully", func(t *testing.T) {
		uc.EXPECT().DeleteLand(gomock.Any(), ids.LandID).Return(nil).Times(1)

		req, _ := http.NewRequest(http.MethodDelete, "/lands/"+ids.LandID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Land deleted successfully", response.Data.Message)
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.EXPECT().DeleteLand(gomock.Any(), ids.LandID).Return(utils.NewInternalError("internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodDelete, "/lands/"+ids.LandID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response response.ResponseMessage
		json.NewDecoder(w.Body).Decode(&response)
		assert.Equal(t, "internal error", response.Errors.Message)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when id invalid", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/lands/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when land not found", func(t *testing.T) {
		uc.EXPECT().DeleteLand(gomock.Any(), ids.LandID).Return(utils.NewNotFoundError("land not found")).Times(1)

		req, _ := http.NewRequest(http.MethodDelete, "/lands/"+ids.LandID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "NOT_FOUND")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestLandHandler_RestoreLand(t *testing.T) {
	r, h, uc, _, ids, mocks := LandHandlerSetup(t)

	r.PATCH("/lands/:id/restore", h.RestoreLand)

	t.Run("should restore land successfully", func(t *testing.T) {

		uc.EXPECT().RestoreLand(gomock.Any(), ids.LandID).Return(mocks.Land, nil).Times(1)

		req, _ := http.NewRequest(http.MethodPatch, "/lands/"+ids.LandID.String()+"/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response responseLandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, mocks.Land.LandArea, response.Data.LandArea)
		assert.Equal(t, mocks.Land.Certificate, response.Data.Certificate)
		assert.Equal(t, mocks.Land.UserID, response.Data.UserID)
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.EXPECT().RestoreLand(gomock.Any(), ids.LandID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodPatch, "/lands/"+ids.LandID.String()+"/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when id invalid", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, "/lands/abc/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when land not found", func(t *testing.T) {
		uc.EXPECT().RestoreLand(gomock.Any(), ids.LandID).Return(nil, utils.NewNotFoundError("land not found")).Times(1)

		req, _ := http.NewRequest(http.MethodPatch, "/lands/"+ids.LandID.String()+"/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "NOT_FOUND")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
