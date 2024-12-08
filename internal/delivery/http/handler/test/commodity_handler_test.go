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

type responseCommodityHandler struct {
	Status  int              `json:"status"`
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    domain.Commodity `json:"data"`
	Errors  interface{}      `json:"errors"`
}

func TestCreateCommodity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockCommodityUsecase(ctrl)
	h := handler.NewCommodityHandler(usecase)
	r := gin.Default()

	r.POST("/commodities", h.CreateCommodity)

	t.Run("Test CreateCommodity, successfully", func(t *testing.T) {
		mockResCommodity := &domain.Commodity{Name: "commodity", Description: "commodity description"}

		usecase.EXPECT().CreateCommodity(gomock.Any(), gomock.Any()).Return(mockResCommodity, nil).Times(1)

		reqBody := `{"name":"commodity","description":"commodity description"}`
		req, _ := http.NewRequest(http.MethodPost, "/commodities", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Test CreateCommodity, bind error", func(t *testing.T) {
		usecase.EXPECT().CreateCommodity(gomock.Any(), gomock.Any()).Times(0)

		req, _ := http.NewRequest(http.MethodPost, "/commodities", bytes.NewReader([]byte("invalid-json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test CreateCommodity, internal error", func(t *testing.T) {
		usecase.EXPECT().CreateCommodity(gomock.Any(), gomock.Any()).Return(nil, errors.New("internal error")).Times(1)

		reqBody := `{"name":"commodity","description":"commodity description"}`
		req, _ := http.NewRequest(http.MethodPost, "/commodities", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetAllCommodities(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockCommodityUsecase(ctrl)
	h := handler.NewCommodityHandler(usecase)
	r := gin.Default()

	r.GET("/commodities", h.GetAllCommodities)

	t.Run("Test GetAllCommodities, successfully", func(t *testing.T) {
		mockCommodities := []domain.Commodity{{Name: "commodity", Description: "commodity description"}}

		usecase.EXPECT().GetAllCommodities(gomock.Any()).Return(&mockCommodities, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/commodities", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response struct {
			Data []domain.Commodity `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Data, len(mockCommodities))
	})

	t.Run("Test GetAllCommodities, internal error", func(t *testing.T) {
		usecase.EXPECT().GetAllCommodities(gomock.Any()).Return(nil, errors.New("internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/commodities", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetCommodityById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockCommodityUsecase(ctrl)
	h := handler.NewCommodityHandler(usecase)
	r := gin.Default()

	r.GET("/commodities/:id", h.GetCommodityById)

	t.Run("Test GetCommodityById, successfully", func(t *testing.T) {
		commodityID := uuid.New()
		mockCommodity := &domain.Commodity{ID: commodityID, Name: "commodity", Description: "commodity description"}

		usecase.EXPECT().GetCommodityById(gomock.Any(), commodityID).Return(mockCommodity, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/commodities/"+commodityID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response responseCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, mockCommodity.Name, response.Data.Name)
		assert.Equal(t, mockCommodity.Description, response.Data.Description)
	})

	t.Run("Test GetCommodityById, database error", func(t *testing.T) {
		commodityID := uuid.New()

		usecase.EXPECT().GetCommodityById(gomock.Any(), commodityID).Return(nil, errors.New("internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/commodities/"+commodityID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Test GetCommodityById, invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/commodities/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUpdateCommodity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockCommodityUsecase(ctrl)
	h := handler.NewCommodityHandler(usecase)
	r := gin.Default()

	r.PATCH("/commodities/:id", h.UpdateCommodity)
	t.Run("Test UpdateCommodity, successfully", func(t *testing.T) {
		commodityID := uuid.New()
		mockCommodity := &domain.Commodity{ID: commodityID, Name: "commodity", Description: "commodity description"}

		usecase.EXPECT().UpdateCommodity(gomock.Any(), commodityID, gomock.Any()).Return(mockCommodity, nil).Times(1)

		reqBody := `{"name":"updated"}`
		req, _ := http.NewRequest(http.MethodPatch, "/commodities/"+commodityID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response responseCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, mockCommodity.Name, response.Data.Name)
		assert.Equal(t, mockCommodity.Description, response.Data.Description)
	})

	t.Run("Test UpdateCommodity, database error", func(t *testing.T) {
		commodityID := uuid.New()

		usecase.EXPECT().UpdateCommodity(gomock.Any(), commodityID, gomock.Any()).Return(nil, errors.New("internal error")).Times(1)

		reqBody := `{"name":"updated"}`
		req, _ := http.NewRequest(http.MethodPatch, "/commodities/"+commodityID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Test UpdateCommodity, bind error", func(t *testing.T) {
		commodityID := uuid.New()

		usecase.EXPECT().UpdateCommodity(gomock.Any(), commodityID, gomock.Any()).Times(0)

		req, _ := http.NewRequest(http.MethodPatch, "/commodities/"+commodityID.String(), bytes.NewReader([]byte(`invalid-json`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test UpdateCommodity, invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, "/commodities/abc", bytes.NewReader([]byte(`invalid-json`)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test UpdateCommodity, not found", func(t *testing.T) {
		commodityID := uuid.New()

		usecase.EXPECT().UpdateCommodity(gomock.Any(), commodityID, gomock.Any()).Return(nil, utils.NewNotFoundError("commodity not found")).Times(1)

		reqBody := `{"name":"updated"}`
		req, _ := http.NewRequest(http.MethodPatch, "/commodities/"+commodityID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestDeleteCommodity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockCommodityUsecase(ctrl)
	h := handler.NewCommodityHandler(usecase)
	r := gin.Default()

	r.DELETE("/commodities/:id", h.DeleteCommodity)

	t.Run("Test DeleteCommodity, successfully", func(t *testing.T) {
		commodityID := uuid.New()

		usecase.EXPECT().DeleteCommodity(gomock.Any(), commodityID).Times(1)

		req, _ := http.NewRequest(http.MethodDelete, "/commodities/"+commodityID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Commodity deleted successfully", response.Data.Message)
	})

	t.Run("Test DeleteCommodity, database error", func(t *testing.T) {
		commodityID := uuid.New()

		usecase.EXPECT().DeleteCommodity(gomock.Any(), commodityID).Return(errors.New("internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodDelete, "/commodities/"+commodityID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Test DeleteCommodity, invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/commodities/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test DeleteCommodity, not found", func(t *testing.T) {
		commodityID := uuid.New()

		usecase.EXPECT().DeleteCommodity(gomock.Any(), commodityID).Return(utils.NewNotFoundError("commodity not found")).Times(1)

		req, _ := http.NewRequest(http.MethodDelete, "/commodities/"+commodityID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestRestoreCommodity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockCommodityUsecase(ctrl)
	h := handler.NewCommodityHandler(usecase)
	r := gin.Default()

	r.PATCH("/commodities/:id/restore", h.RestoreCommodity)

	t.Run("Test RestoreCommodity, successfully", func(t *testing.T) {
		commodityID := uuid.New()
		mockCommodity := &domain.Commodity{ID: commodityID, Name: "commodity", Description: "commodity description"}

		usecase.EXPECT().RestoreCommodity(gomock.Any(), commodityID).Return(mockCommodity, nil).Times(1)

		req, _ := http.NewRequest(http.MethodPatch, "/commodities/"+commodityID.String()+"/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response responseCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, mockCommodity.Name, response.Data.Name)
		assert.Equal(t, mockCommodity.Description, response.Data.Description)
	})

	t.Run("Test RestoreCommodity, database error", func(t *testing.T) {
		commodityID := uuid.New()

		usecase.EXPECT().RestoreCommodity(gomock.Any(), commodityID).Return(nil, errors.New("internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodPatch, "/commodities/"+commodityID.String()+"/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Test RestoreCommodity, invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, "/commodities/abc/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test RestoreCommodity, not found", func(t *testing.T) {
		commodityID := uuid.New()

		usecase.EXPECT().RestoreCommodity(gomock.Any(), commodityID).Return(nil, utils.NewNotFoundError("commodity not found")).Times(1)

		req, _ := http.NewRequest(http.MethodPatch, "/commodities/"+commodityID.String()+"/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
