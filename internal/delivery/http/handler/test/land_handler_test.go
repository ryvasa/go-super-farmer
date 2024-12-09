package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/test/response"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/usecase/mock"
	"github.com/ryvasa/go-super-farmer/utils"
	mockAuthUtil "github.com/ryvasa/go-super-farmer/utils/mock"
	"github.com/stretchr/testify/assert"
)

type responseLandHandler struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    domain.Land `json:"data"`
	Errors  interface{} `json:"errors"`
}

func TestCreateLand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUtil := mockAuthUtil.NewMockAuthUtil(ctrl)
	usecase := mock.NewMockLandUsecase(ctrl)
	h := handler.NewLandHandler(usecase, authUtil)
	r := gin.Default()
	r.POST("/lands", h.CreateLand)

	t.Run("Test CreateLand, successfully", func(t *testing.T) {
		userID := uuid.New()
		landID := uuid.New()
		mockLand := &domain.Land{ID: landID, UserID: userID, LandArea: 10, Certificate: "test"}

		authUtil.EXPECT().GetAuthUserID(gomock.Any()).Return(userID, nil).Times(1)
		usecase.EXPECT().CreateLand(gomock.Any(), userID, gomock.Any()).Return(mockLand, nil).Times(1)

		reqBody := `{"land_area":10,"certificate":"test"}`
		req, _ := http.NewRequest(http.MethodPost, "/lands", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Test CreateLand, when GetAuthUserID fails", func(t *testing.T) {
		authUtil.EXPECT().GetAuthUserID(gomock.Any()).Return(uuid.UUID{}, utils.NewUnauthorizedError("unauthorized")).Times(1)

		reqBody := `{"land_area":10,"certificate":"test"}`
		req, _ := http.NewRequest(http.MethodPost, "/lands", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Test CreateLand, when ShouldBindJSON fails", func(t *testing.T) {
		userID := uuid.New()
		authUtil.EXPECT().GetAuthUserID(gomock.Any()).Return(userID, nil).Times(1)

		reqBody := `{"land_area":"invalid","certificate":"test"}`
		req, _ := http.NewRequest(http.MethodPost, "/lands", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test CreateLand, when usecase.CreateLand fails", func(t *testing.T) {
		userID := uuid.New()

		authUtil.EXPECT().GetAuthUserID(gomock.Any()).Return(userID, nil).Times(1)
		usecase.EXPECT().CreateLand(gomock.Any(), userID, gomock.Any()).Return(nil, utils.NewInternalError("internal server error")).Times(1)

		reqBody := `{"land_area":10,"certificate":"test"}`
		req, _ := http.NewRequest(http.MethodPost, "/lands", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetOneLand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUtil := mockAuthUtil.NewMockAuthUtil(ctrl)
	usecase := mock.NewMockLandUsecase(ctrl)
	h := handler.NewLandHandler(usecase, authUtil)
	r := gin.Default()
	r.GET("/lands/:id", h.GetLandByID)

	landID := uuid.New()
	userID := uuid.New()
	t.Run("Test GetLandByID, successfully", func(t *testing.T) {
		mockLand := &domain.Land{ID: landID, UserID: userID, LandArea: 10, Certificate: "test"}

		usecase.EXPECT().GetLandByID(gomock.Any(), landID).Return(mockLand, nil).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/lands/"+landID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Test GetLandByID, database error", func(t *testing.T) {
		usecase.EXPECT().GetLandByID(gomock.Any(), landID).Return(nil, errors.New("internal error")).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/lands/"+landID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Test GetLandByID, invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/lands/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetOneLandByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUtil := mockAuthUtil.NewMockAuthUtil(ctrl)
	usecase := mock.NewMockLandUsecase(ctrl)
	h := handler.NewLandHandler(usecase, authUtil)
	r := gin.Default()
	r.GET("/lands/user/:id", h.GetLandByUserID)

	userID := uuid.New()

	t.Run("Test GetLandByUserID, successfully", func(t *testing.T) {
		mockLand1 := &domain.Land{ID: uuid.New(), UserID: uuid.New(), LandArea: 10, Certificate: "test"}
		mockLand2 := &domain.Land{ID: uuid.New(), UserID: uuid.New(), LandArea: 20, Certificate: "test2"}

		usecase.EXPECT().GetLandByUserID(gomock.Any(), userID).Return(&[]domain.Land{*mockLand1, *mockLand2}, nil).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/lands/user/"+userID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Test GetLandByUserID, database error", func(t *testing.T) {
		usecase.EXPECT().GetLandByUserID(gomock.Any(), userID).Return(nil, errors.New("internal error")).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/lands/user/"+userID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Test GetLandByUserID, invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/lands/user/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetAllLands(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUtil := mockAuthUtil.NewMockAuthUtil(ctrl)
	usecase := mock.NewMockLandUsecase(ctrl)
	h := handler.NewLandHandler(usecase, authUtil)
	r := gin.Default()
	r.GET("/lands", h.GetAllLands)

	t.Run("Test GetAllLands, successfully", func(t *testing.T) {
		mockLand1 := &domain.Land{ID: uuid.New(), UserID: uuid.New(), LandArea: 10, Certificate: "test"}
		mockLand2 := &domain.Land{ID: uuid.New(), UserID: uuid.New(), LandArea: 20, Certificate: "test2"}

		usecase.EXPECT().GetAllLands(gomock.Any()).Return(&[]domain.Land{*mockLand1, *mockLand2}, nil).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/lands", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Test GetAllLands, database error", func(t *testing.T) {
		usecase.EXPECT().GetAllLands(gomock.Any()).Return(nil, errors.New("internal error")).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/lands", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUpdateLand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUtil := mockAuthUtil.NewMockAuthUtil(ctrl)
	usecase := mock.NewMockLandUsecase(ctrl)
	h := handler.NewLandHandler(usecase, authUtil)
	r := gin.Default()
	r.PATCH("/lands/:id", h.UpdateLand)

	t.Run("Test UpdateLand, successfully", func(t *testing.T) {
		userID := uuid.New()
		landID := uuid.New()
		mockLand := &domain.Land{ID: landID, UserID: userID, LandArea: 10, Certificate: "test"}

		authUtil.EXPECT().GetAuthUserID(gomock.Any()).Return(userID, nil).Times(1)
		usecase.EXPECT().UpdateLand(gomock.Any(), userID, landID, gomock.Any()).Return(mockLand, nil).Times(1)

		reqBody := `{"land_area":10,"certificate":"test"}`
		req, _ := http.NewRequest(http.MethodPatch, "/lands/"+landID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response responseLandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, mockLand.LandArea, response.Data.LandArea)
		assert.Equal(t, mockLand.Certificate, response.Data.Certificate)
	})

	t.Run("Test UpdateLand, database error", func(t *testing.T) {
		userID := uuid.New()
		landID := uuid.New()

		authUtil.EXPECT().GetAuthUserID(gomock.Any()).Return(userID, nil).Times(1)
		usecase.EXPECT().UpdateLand(gomock.Any(), userID, landID, gomock.Any()).Return(nil, utils.NewInternalError("internal server error")).Times(1)

		reqBody := `{"land_area":10,"certificate":"test"}`
		req, _ := http.NewRequest(http.MethodPatch, "/lands/"+landID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Test UpdateLand, invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, "/lands/abc", bytes.NewReader([]byte(`invalid-json`)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test UpdateLand, bind error", func(t *testing.T) {
		userID := uuid.New()
		landID := uuid.New()
		authUtil.EXPECT().GetAuthUserID(gomock.Any()).Return(userID, nil).Times(1)
		usecase.EXPECT().UpdateLand(gomock.Any(), userID, landID, gomock.Any()).Times(0)
		req, _ := http.NewRequest(http.MethodPatch, "/lands/"+landID.String(), bytes.NewReader([]byte(`invalid-json`)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test UpdateLand, not found", func(t *testing.T) {
		userID := uuid.New()
		landID := uuid.New()

		authUtil.EXPECT().GetAuthUserID(gomock.Any()).Return(userID, nil).Times(1)
		usecase.EXPECT().UpdateLand(gomock.Any(), userID, landID, gomock.Any()).Return(nil, utils.NewNotFoundError("land not found")).Times(1)

		reqBody := `{"land_area":10,"certificate":"test"}`
		req, _ := http.NewRequest(http.MethodPatch, "/lands/"+landID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestDeleteLand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUtil := mockAuthUtil.NewMockAuthUtil(ctrl)
	usecase := mock.NewMockLandUsecase(ctrl)
	h := handler.NewLandHandler(usecase, authUtil)
	r := gin.Default()
	r.DELETE("/lands/:id", h.DeleteLand)

	t.Run("Test DeleteLand, successfully", func(t *testing.T) {
		landID := uuid.New()

		usecase.EXPECT().DeleteLand(gomock.Any(), landID).Return(nil).Times(1)

		req, _ := http.NewRequest(http.MethodDelete, "/lands/"+landID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Land deleted successfully", response.Data.Message)
	})

	t.Run("Test DeleteLand, database error", func(t *testing.T) {
		landID := uuid.New()

		usecase.EXPECT().DeleteLand(gomock.Any(), landID).Return(errors.New("internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodDelete, "/lands/"+landID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Test DeleteLand, invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/lands/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test DeleteLand, not found", func(t *testing.T) {
		landID := uuid.New()

		usecase.EXPECT().DeleteLand(gomock.Any(), landID).Return(utils.NewNotFoundError("land not found")).Times(1)

		req, _ := http.NewRequest(http.MethodDelete, "/lands/"+landID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestRestoreLand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUtil := mockAuthUtil.NewMockAuthUtil(ctrl)
	usecase := mock.NewMockLandUsecase(ctrl)
	h := handler.NewLandHandler(usecase, authUtil)
	r := gin.Default()
	r.PATCH("/lands/:id/restore", h.RestoreLand)

	t.Run("Test RestoreLand, successfully", func(t *testing.T) {
		landID := uuid.New()
		userID := uuid.New()
		mockLand := &domain.Land{ID: landID, UserID: userID, LandArea: 10, Certificate: "test"}

		usecase.EXPECT().RestoreLand(gomock.Any(), landID).Return(mockLand, nil).Times(1)

		req, _ := http.NewRequest(http.MethodPatch, "/lands/"+landID.String()+"/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response responseLandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, mockLand.LandArea, response.Data.LandArea)
		assert.Equal(t, mockLand.Certificate, response.Data.Certificate)
		assert.Equal(t, mockLand.UserID, response.Data.UserID)
	})

	t.Run("Test RestoreLand, database error", func(t *testing.T) {
		landID := uuid.New()

		usecase.EXPECT().RestoreLand(gomock.Any(), landID).Return(nil, errors.New("internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodPatch, "/lands/"+landID.String()+"/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Test RestoreLand, invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, "/lands/abc/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test RestoreLand, not found", func(t *testing.T) {
		landID := uuid.New()

		usecase.EXPECT().RestoreLand(gomock.Any(), landID).Return(nil, utils.NewNotFoundError("land not found")).Times(1)

		req, _ := http.NewRequest(http.MethodPatch, "/lands/"+landID.String()+"/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
