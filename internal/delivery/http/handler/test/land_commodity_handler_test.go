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
	"github.com/stretchr/testify/assert"
)

type responseLandCommodityHandler struct {
	Status  int                  `json:"status"`
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Data    domain.LandCommodity `json:"data"`
	Errors  interface{}          `json:"errors"`
}

func TestLandCommodityHandler_CreateLandCommodity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockLandCommodityUsecase(ctrl)
	h := handler.NewLandCommodityHandler(usecase)
	r := gin.Default()
	r.POST("/land_commodities", h.CreateLandCommodity)

	landID := uuid.New()
	commodityID := uuid.New()
	landCommodityID := uuid.New()

	t.Run("Test CreateLandCommodity, successfully", func(t *testing.T) {
		mockLandCommodity := &domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(100)}

		usecase.EXPECT().CreateLandCommodity(gomock.Any(), gomock.Any()).Return(mockLandCommodity, nil).Times(1)

		reqBody := `{"land_id":"` + landID.String() + `","commodity_id":"` + commodityID.String() + `","land_area":100}`
		req, _ := http.NewRequest(http.MethodPost, "/land_commodities", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

	})

	t.Run("Test CreateLandCommodity, bind error", func(t *testing.T) {
		reqBody := `{"land_id":"` + landID.String() + `","commodity_id":"` + commodityID.String() + `","land_area":"invalid"}`
		req, _ := http.NewRequest(http.MethodPost, "/land_commodities", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test CreateLandCommodity, usecase error", func(t *testing.T) {
		usecase.EXPECT().CreateLandCommodity(gomock.Any(), gomock.Any()).Return(nil, utils.NewInternalError("internal server error")).Times(1)

		reqBody := `{"land_id":"` + landID.String() + `","commodity_id":"` + commodityID.String() + `","land_area":100}`
		req, _ := http.NewRequest(http.MethodPost, "/land_commodities", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestLandCommodityHandler_GetLandCommodityByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockLandCommodityUsecase(ctrl)
	h := handler.NewLandCommodityHandler(usecase)
	r := gin.Default()
	r.GET("/land_commodities/:id", h.GetLandCommodityByID)

	landCommodityID := uuid.New()
	commodityID := uuid.New()
	landID := uuid.New()

	t.Run("Test GetLandCommodityByID, successfully", func(t *testing.T) {
		mockLandCommodity := &domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(100)}

		usecase.EXPECT().GetLandCommodityByID(gomock.Any(), landCommodityID).Return(mockLandCommodity, nil).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/land_commodities/"+landCommodityID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Test GetLandCommodityByID, database error", func(t *testing.T) {
		usecase.EXPECT().GetLandCommodityByID(gomock.Any(), landCommodityID).Return(nil, errors.New("internal error")).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/land_commodities/"+landCommodityID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Test GetLandCommodityByID, invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/land_commodities/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLandCommodityHandler_GetLandCommodityByLandID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockLandCommodityUsecase(ctrl)
	h := handler.NewLandCommodityHandler(usecase)
	r := gin.Default()
	r.GET("/land_commodities/land/:id", h.GetLandCommodityByLandID)

	landCommodityID := uuid.New()
	commodityID := uuid.New()
	landID := uuid.New()

	t.Run("Test GetLandCommodityByLandID, successfully", func(t *testing.T) {
		mockLandCommodity1 := &domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(100)}
		mockLandCommodity2 := &domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(200)}

		usecase.EXPECT().GetLandCommodityByLandID(gomock.Any(), landID).Return(&[]domain.LandCommodity{*mockLandCommodity1, *mockLandCommodity2}, nil).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/land_commodities/land/"+landID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Test GetLandCommodityByLandID, database error", func(t *testing.T) {
		usecase.EXPECT().GetLandCommodityByLandID(gomock.Any(), landID).Return(nil, errors.New("internal error")).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/land_commodities/land/"+landID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Test GetLandCommodityByLandID, invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/land_commodities/land/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLandCommodityHandler_GetAllLandCommodity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockLandCommodityUsecase(ctrl)
	h := handler.NewLandCommodityHandler(usecase)
	r := gin.Default()
	r.GET("/land_commodities", h.GetAllLandCommodity)

	landCommodityID := uuid.New()
	commodityID := uuid.New()
	landID := uuid.New()

	t.Run("Test GetAllLandCommodity, successfully", func(t *testing.T) {
		mockLandCommodity1 := &domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(100)}
		mockLandCommodity2 := &domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(200)}

		usecase.EXPECT().GetAllLandCommodity(gomock.Any()).Return(&[]domain.LandCommodity{*mockLandCommodity1, *mockLandCommodity2}, nil).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/land_commodities", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Test GetAllLandCommodity, database error", func(t *testing.T) {
		usecase.EXPECT().GetAllLandCommodity(gomock.Any()).Return(nil, errors.New("internal error")).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/land_commodities", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestLandCommodityHandler_GetLandCommodityByCommodityID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockLandCommodityUsecase(ctrl)
	h := handler.NewLandCommodityHandler(usecase)
	r := gin.Default()
	r.GET("/land_commodities/commodity/:id", h.GetLandCommodityByCommodityID)

	landCommodityID := uuid.New()
	commodityID := uuid.New()
	landID := uuid.New()

	t.Run("Test GetLandCommodityByCommodityID, successfully", func(t *testing.T) {
		mockLandCommodity1 := &domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(100)}
		mockLandCommodity2 := &domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(200)}

		usecase.EXPECT().GetLandCommodityByCommodityID(gomock.Any(), commodityID).Return(&[]domain.LandCommodity{*mockLandCommodity1, *mockLandCommodity2}, nil).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/land_commodities/commodity/"+commodityID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Test GetLandCommodityByCommodityID, internal server error", func(t *testing.T) {

		usecase.EXPECT().GetLandCommodityByCommodityID(gomock.Any(), commodityID).Return(nil, errors.New("Internal server error")).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/land_commodities/commodity/"+commodityID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Test GetLandCommodityByCommodityID, invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/land_commodities/commodity/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLandCommodityHandler_UpdateLandCommodity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockLandCommodityUsecase(ctrl)
	h := handler.NewLandCommodityHandler(usecase)
	r := gin.Default()
	r.PATCH("/land_commodities/:id", h.UpdateLandCommodity)

	landCommodityID := uuid.New()
	commodityID := uuid.New()
	landID := uuid.New()

	t.Run("Test UpdateLandCommodity, successfully", func(t *testing.T) {
		mockLandCommodity := &domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(100)}

		usecase.EXPECT().UpdateLandCommodity(gomock.Any(), landCommodityID, gomock.Any()).Return(mockLandCommodity, nil).Times(1)
		reqBody := `{"land_area":100}`
		req, _ := http.NewRequest(http.MethodPatch, "/land_commodities/"+landCommodityID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Test UpdateLandCommodity, bind error", func(t *testing.T) {
		reqBody := `{"land_area":"invalid"}`
		req, _ := http.NewRequest(http.MethodPatch, "/land_commodities/"+landCommodityID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test UpdateLandCommodity, usecase error", func(t *testing.T) {
		usecase.EXPECT().UpdateLandCommodity(gomock.Any(), landCommodityID, gomock.Any()).Return(nil, utils.NewInternalError("internal server error")).Times(1)

		reqBody := `{"land_area":100}`
		req, _ := http.NewRequest(http.MethodPatch, "/land_commodities/"+landCommodityID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Test UpdateLandCommodity, invalid id", func(t *testing.T) {
		reqBody := `{"land_area":100}`
		req, _ := http.NewRequest(http.MethodPatch, "/land_commodities/abc", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLandCommodityHandler_DeleteLandCommodity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockLandCommodityUsecase(ctrl)
	h := handler.NewLandCommodityHandler(usecase)
	r := gin.Default()
	r.DELETE("/land_commodities/:id", h.DeleteLandCommodity)

	landCommodityID := uuid.New()

	t.Run("Test DeleteLandCommodity, successfully", func(t *testing.T) {
		usecase.EXPECT().DeleteLandCommodity(gomock.Any(), landCommodityID).Return(nil).Times(1)
		req, _ := http.NewRequest(http.MethodDelete, "/land_commodities/"+landCommodityID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.Nil(t, err)
		assert.Equal(t, "Land commodity deleted successfully", response.Data.Message)
	})

	t.Run("Test DeleteLandCommodity, usecase error", func(t *testing.T) {
		usecase.EXPECT().DeleteLandCommodity(gomock.Any(), landCommodityID).Return(utils.NewInternalError("internal server error")).Times(1)
		req, _ := http.NewRequest(http.MethodDelete, "/land_commodities/"+landCommodityID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Test DeleteLandCommodity, invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/land_commodities/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLandCommodityHandler_RestoreLandCommodity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockLandCommodityUsecase(ctrl)
	h := handler.NewLandCommodityHandler(usecase)
	r := gin.Default()
	r.PATCH("/land_commodities/:id/restore", h.RestoreLandCommodity)

	landCommodityID := uuid.New()

	t.Run("Test RestoreLandCommodity, successfully", func(t *testing.T) {
		mockLandCommodity := &domain.LandCommodity{ID: landCommodityID, CommodityID: landCommodityID, LandID: landCommodityID, LandArea: float64(100)}
		usecase.EXPECT().RestoreLandCommodity(gomock.Any(), landCommodityID).Return(mockLandCommodity, nil).Times(1)
		req, _ := http.NewRequest(http.MethodPatch, "/land_commodities/"+landCommodityID.String()+"/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response responseLandCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.Nil(t, err)
		assert.Equal(t, response.Data.ID, landCommodityID)
		assert.Equal(t, response.Data.CommodityID, landCommodityID)
		assert.Equal(t, response.Data.LandID, landCommodityID)
		assert.Equal(t, response.Data.LandArea, float64(100))
	})

	t.Run("Test RestoreLandCommodity, usecase error", func(t *testing.T) {
		usecase.EXPECT().RestoreLandCommodity(gomock.Any(), landCommodityID).Return(nil, utils.NewInternalError("internal server error")).Times(1)
		req, _ := http.NewRequest(http.MethodPatch, "/land_commodities/"+landCommodityID.String()+"/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Test RestoreLandCommodity, invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, "/land_commodities/abc/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
