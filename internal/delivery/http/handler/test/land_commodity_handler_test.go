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
	"github.com/stretchr/testify/assert"
)

type responseLandCommodityHandler struct {
	Status  int                  `json:"status"`
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Data    domain.LandCommodity `json:"data"`
	Errors  response.Error       `json:"errors"`
}

type responseLandCommoditiesHandler struct {
	Status  int                    `json:"status"`
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    []domain.LandCommodity `json:"data"`
	Errors  response.Error         `json:"errors"`
}

type LandCommodityHandlerMocks struct {
	LandCommodity   *domain.LandCommodity
	LandCommodities *[]domain.LandCommodity
}

type LandCommodityHandlerIDs struct {
	LandCommodityID uuid.UUID
	CommodityID     uuid.UUID
	LandID          uuid.UUID
}

func LandCommodityHandlerSetUp(t *testing.T) (*gin.Engine, handler_interface.LandCommodityHandler, *mock.MockLandCommodityUsecase, LandCommodityHandlerIDs, LandCommodityHandlerMocks) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc := mock.NewMockLandCommodityUsecase(ctrl)
	h := handler_implementation.NewLandCommodityHandler(uc)
	r := gin.Default()

	landCommodityID := uuid.New()
	commodityID := uuid.New()
	landID := uuid.New()
	mocks := LandCommodityHandlerMocks{
		LandCommodity: &domain.LandCommodity{
			ID:          landCommodityID,
			CommodityID: commodityID,
			LandID:      landID,
			LandArea:    float64(100),
		},
		LandCommodities: &[]domain.LandCommodity{
			{
				ID:          landCommodityID,
				CommodityID: commodityID,
				LandID:      landID,
				LandArea:    float64(100),
			},
		},
	}
	ids := LandCommodityHandlerIDs{
		LandCommodityID: landCommodityID,
		CommodityID:     commodityID,
		LandID:          landID,
	}

	return r, h, uc, ids, mocks
}

func TestLandCommodityHandler_CreateLandCommodity(t *testing.T) {
	r, h, uc, ids, mocks := LandCommodityHandlerSetUp(t)
	r.POST("/land_commodities", h.CreateLandCommodity)

	t.Run("should create land commodity successfully", func(t *testing.T) {
		uc.EXPECT().CreateLandCommodity(gomock.Any(), gomock.Any()).Return(mocks.LandCommodity, nil).Times(1)

		reqBody := `{"land_id":"` + ids.LandID.String() + `","commodity_id":"` + ids.CommodityID.String() + `","land_area":100}`
		req, _ := http.NewRequest(http.MethodPost, "/land_commodities", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, response.Data.ID, ids.LandCommodityID)
		assert.Equal(t, response.Data.LandID, ids.LandID)
		assert.Equal(t, response.Data.CommodityID, ids.CommodityID)
		assert.Equal(t, response.Data.LandArea, float64(100))
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		reqBody := `{"land_id":"` + ids.LandID.String() + `","commodity_id":"` + ids.CommodityID.String() + `","land_area":"invalid"}`
		req, _ := http.NewRequest(http.MethodPost, "/land_commodities", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.EXPECT().CreateLandCommodity(gomock.Any(), gomock.Any()).Return(nil, utils.NewInternalError("internal server error")).Times(1)

		reqBody := `{"land_id":"` + ids.LandID.String() + `","commodity_id":"` + ids.CommodityID.String() + `","land_area":100}`
		req, _ := http.NewRequest(http.MethodPost, "/land_commodities", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestLandCommodityHandler_GetLandCommodityByID(t *testing.T) {
	r, h, uc, ids, mocks := LandCommodityHandlerSetUp(t)

	r.GET("/land_commodities/:id", h.GetLandCommodityByID)
	t.Run("should return land commodity by id successfully", func(t *testing.T) {
		uc.EXPECT().GetLandCommodityByID(gomock.Any(), ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/land_commodities/"+ids.LandCommodityID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, response.Data.ID, ids.LandCommodityID)
		assert.Equal(t, response.Data.CommodityID, ids.CommodityID)
		assert.Equal(t, response.Data.LandID, ids.LandID)
		assert.Equal(t, response.Data.LandArea, float64(100))
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.EXPECT().GetLandCommodityByID(gomock.Any(), ids.LandCommodityID).Return(nil, utils.NewInternalError("internal error")).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/land_commodities/"+ids.LandCommodityID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/land_commodities/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
	})
}

func TestLandCommodityHandler_GetLandCommodityByLandID(t *testing.T) {
	r, h, uc, ids, mocks := LandCommodityHandlerSetUp(t)
	r.GET("/land_commodities/land/:id", h.GetLandCommodityByLandID)

	t.Run("should return land commodity by land id successfully", func(t *testing.T) {

		uc.EXPECT().GetLandCommodityByLandID(gomock.Any(), ids.LandID).Return(mocks.LandCommodities, nil).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/land_commodities/land/"+ids.LandID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandCommoditiesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Data, len(*mocks.LandCommodities))
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.EXPECT().GetLandCommodityByLandID(gomock.Any(), ids.LandID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/land_commodities/land/"+ids.LandID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandCommoditiesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/land_commodities/land/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandCommoditiesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLandCommodityHandler_GetAllLandCommodity(t *testing.T) {
	r, h, uc, _, mocks := LandCommodityHandlerSetUp(t)
	r.GET("/land_commodities", h.GetAllLandCommodity)

	t.Run("should return all land commodities successfully", func(t *testing.T) {
		uc.EXPECT().GetAllLandCommodity(gomock.Any()).Return(mocks.LandCommodities, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/land_commodities", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandCommoditiesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Data, len(*mocks.LandCommodities))
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.EXPECT().GetAllLandCommodity(gomock.Any()).Return(nil, utils.NewInternalError("internal error")).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/land_commodities", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandCommoditiesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestLandCommodityHandler_GetLandCommodityByCommodityID(t *testing.T) {
	r, h, uc, ids, mocks := LandCommodityHandlerSetUp(t)
	r.GET("/land_commodities/commodity/:id", h.GetLandCommodityByCommodityID)

	t.Run("should return land commodity by commodity id successfully", func(t *testing.T) {
		uc.EXPECT().GetLandCommodityByCommodityID(gomock.Any(), ids.CommodityID).Return(mocks.LandCommodities, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/land_commodities/commodity/"+ids.CommodityID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandCommoditiesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Data, len(*mocks.LandCommodities))
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.EXPECT().GetLandCommodityByCommodityID(gomock.Any(), ids.CommodityID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/land_commodities/commodity/"+ids.CommodityID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandCommoditiesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/land_commodities/commodity/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandCommoditiesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLandCommodityHandler_UpdateLandCommodity(t *testing.T) {
	r, h, uc, ids, mocks := LandCommodityHandlerSetUp(t)
	r.PATCH("/land_commodities/:id", h.UpdateLandCommodity)

	t.Run("should update land commodity successfully", func(t *testing.T) {
		uc.EXPECT().UpdateLandCommodity(gomock.Any(), ids.LandCommodityID, gomock.Any()).Return(mocks.LandCommodity, nil).Times(1)

		reqBody := `{"land_area":100}`
		req, _ := http.NewRequest(http.MethodPatch, "/land_commodities/"+ids.LandCommodityID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, response.Data.ID, ids.LandCommodityID)
		assert.Equal(t, response.Data.LandID, ids.LandID)
		assert.Equal(t, response.Data.CommodityID, ids.CommodityID)
		assert.Equal(t, response.Data.LandArea, float64(100))
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		reqBody := `{"land_area":"invalid"}`
		req, _ := http.NewRequest(http.MethodPatch, "/land_commodities/"+ids.LandCommodityID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.EXPECT().UpdateLandCommodity(gomock.Any(), ids.LandCommodityID, gomock.Any()).Return(nil, utils.NewInternalError("internal server error")).Times(1)

		reqBody := `{"land_area":100}`
		req, _ := http.NewRequest(http.MethodPatch, "/land_commodities/"+ids.LandCommodityID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		reqBody := `{"land_area":100}`
		req, _ := http.NewRequest(http.MethodPatch, "/land_commodities/abc", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLandCommodityHandler_DeleteLandCommodity(t *testing.T) {
	r, h, uc, ids, _ := LandCommodityHandlerSetUp(t)
	r.DELETE("/land_commodities/:id", h.DeleteLandCommodity)

	t.Run("should delete land commodity successfully", func(t *testing.T) {
		uc.EXPECT().DeleteLandCommodity(gomock.Any(), ids.LandCommodityID).Return(nil).Times(1)
		req, _ := http.NewRequest(http.MethodDelete, "/land_commodities/"+ids.LandCommodityID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.Nil(t, err)
		assert.Equal(t, "Land commodity deleted successfully", response.Data.Message)
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.EXPECT().DeleteLandCommodity(gomock.Any(), ids.LandCommodityID).Return(utils.NewInternalError("internal server error")).Times(1)
		req, _ := http.NewRequest(http.MethodDelete, "/land_commodities/"+ids.LandCommodityID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		response := response.ResponseMessage{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/land_commodities/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLandCommodityHandler_RestoreLandCommodity(t *testing.T) {
	r, h, uc, ids, mocks := LandCommodityHandlerSetUp(t)
	r.PATCH("/land_commodities/:id/restore", h.RestoreLandCommodity)

	t.Run("should restore land commodity successfully", func(t *testing.T) {
		uc.EXPECT().RestoreLandCommodity(gomock.Any(), ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)

		req, _ := http.NewRequest(http.MethodPatch, "/land_commodities/"+ids.LandCommodityID.String()+"/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response responseLandCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.Nil(t, err)
		assert.Equal(t, response.Data.ID, ids.LandCommodityID)
		assert.Equal(t, response.Data.CommodityID, ids.CommodityID)
		assert.Equal(t, response.Data.LandID, ids.LandID)
		assert.Equal(t, response.Data.LandArea, float64(100))
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.EXPECT().RestoreLandCommodity(gomock.Any(), ids.LandCommodityID).Return(nil, utils.NewInternalError("internal server error")).Times(1)
		req, _ := http.NewRequest(http.MethodPatch, "/land_commodities/"+ids.LandCommodityID.String()+"/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, "/land_commodities/abc/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseLandCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
