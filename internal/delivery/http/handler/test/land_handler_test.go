package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/usecase/mock"
	"github.com/ryvasa/go-super-farmer/utils"
	mockAuthUtil "github.com/ryvasa/go-super-farmer/utils/mock"
	"github.com/stretchr/testify/assert"
)

type ResponseLandhHandler struct {
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
	r.GET("/lands:id", h.GetLandByID)

	t.Run("Test GetLandByID, successfully", func(t *testing.T) {

	})

	t.Run("Test GetLandByID, database error", func(t *testing.T) {

	})

	t.Run("Test GetLandByID, invalid id", func(t *testing.T) {

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

	t.Run("Test GetLandByUserID, successfully", func(t *testing.T) {

	})

	t.Run("Test GetLandByUserID, database error", func(t *testing.T) {

	})

	t.Run("Test GetLandByUserID, invalid id", func(t *testing.T) {

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

	})

	t.Run("Test GetAllLands, database error", func(t *testing.T) {

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

	})

	t.Run("Test UpdateLand, database error", func(t *testing.T) {

	})

	t.Run("Test UpdateLand, invalid id", func(t *testing.T) {

	})

	t.Run("Test UpdateLand, bind error", func(t *testing.T) {

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

	})

	t.Run("Test DeleteLand, database error", func(t *testing.T) {

	})

	t.Run("Test DeleteLand, invalid id", func(t *testing.T) {

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

	})

	t.Run("Test RestoreLand, database error", func(t *testing.T) {

	})

	t.Run("Test RestoreLand, invalid id", func(t *testing.T) {

	})
}
