package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

type responseHarvestHandler struct {
	Status  int            `json:"status"`
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    domain.Harvest `json:"data"`
	Errors  response.Error `json:"errors"`
}

type HarvestHandlerMocks struct {
	Harvest        *domain.Harvest
	Harvests       *[]domain.Harvest
	UpdatedHarvest *domain.Harvest
}

type HarvestHandlerIDs struct {
	HarvestID       uuid.UUID
	LandCommodityID uuid.UUID
	RegionID        uuid.UUID
	LandID          uuid.UUID
	CommodityID     uuid.UUID
}

func HarvestHandlerSetUp(t *testing.T) (*gin.Engine, handler_interface.HarvestHandler, *mock.MockHarvestUsecase, HarvestHandlerIDs, HarvestHandlerMocks) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc := mock.NewMockHarvestUsecase(ctrl)
	h := handler_implementation.NewHarvestHandler(uc)
	r := gin.Default()

	harvestID := uuid.New()
	landCommodityID := uuid.New()
	regionID := uuid.New()
	landID := uuid.New()
	commodityID := uuid.New()
	quantity := float64(100)
	unit := "kg"
	date := "2022-01-01"
	harvestDate, _ := time.Parse("2006-01-02", date)
	ids := HarvestHandlerIDs{
		HarvestID:       harvestID,
		LandCommodityID: landCommodityID,
		RegionID:        regionID,
		LandID:          landID,
		CommodityID:     commodityID,
	}

	mocks := HarvestHandlerMocks{
		Harvest: &domain.Harvest{
			ID:              harvestID,
			LandCommodityID: landCommodityID,
			RegionID:        regionID,
			Quantity:        quantity,
			Unit:            unit,
			HarvestDate:     harvestDate,
		},
		Harvests: &[]domain.Harvest{
			{
				ID:              harvestID,
				LandCommodityID: landCommodityID,
				RegionID:        regionID,
				Quantity:        quantity,
				Unit:            unit,
				HarvestDate:     harvestDate,
			},
		},
		UpdatedHarvest: &domain.Harvest{
			ID:              harvestID,
			LandCommodityID: landCommodityID,
			RegionID:        regionID,
			Quantity:        quantity + 1,
			Unit:            unit,
			HarvestDate:     harvestDate,
		},
	}

	return r, h, uc, ids, mocks
}

func TestHarvestHandler_CreateHarvest(t *testing.T) {
	r, h, uc, _, mocks := HarvestHandlerSetUp(t)
	r.POST("/harvests", h.CreateHarvest)

	t.Run("should create harvest successfully", func(t *testing.T) {
		uc.EXPECT().CreateHarvest(gomock.Any(), gomock.Any()).Return(mocks.Harvest, nil).Times(1)

		reqBody := `{"land_commodity_id":"` + mocks.Harvest.LandCommodityID.String() + `","region_id":"` + mocks.Harvest.RegionID.String() + `","quantity":100,"unit":"kg","harvest_date":"2022-01-01"}`

		req, _ := http.NewRequest(http.MethodPost, "/harvests", bytes.NewReader([]byte(reqBody)))

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, response.Data.Quantity, float64(100))
		assert.Equal(t, response.Data.Unit, "kg")
		assert.Equal(t, response.Data.HarvestDate, mocks.Harvest.HarvestDate)
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		uc.EXPECT().CreateHarvest(gomock.Any(), gomock.Any()).Times(0)

		req, _ := http.NewRequest(http.MethodPost, "/harvests", bytes.NewReader([]byte(`invalid-json`)))

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().CreateHarvest(gomock.Any(), gomock.Any()).Return(nil, utils.NewInternalError("Internal error"))

		reqBody := `{"land_commodity_id":"` + mocks.Harvest.LandCommodityID.String() + `","region_id":"` + mocks.Harvest.RegionID.String() + `","quantity":100,"unit":"kg","harvest_date":"2022-01-01"}`

		req, _ := http.NewRequest(http.MethodPost, "/harvests", bytes.NewReader([]byte(reqBody)))

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHarvestHandler_GetAllHarvest(t *testing.T) {
	r, h, uc, _, mocks := HarvestHandlerSetUp(t)
	r.GET("/harvests", h.GetAllHarvest)
	t.Run("should return all harvests successfully", func(t *testing.T) {
		uc.EXPECT().GetAllHarvest(gomock.Any()).Return(mocks.Harvests, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/harvests", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var response struct {
			Data []domain.Harvest `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Data, len(*mocks.Harvests))
	})
	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().GetAllHarvest(gomock.Any()).Return(nil, utils.NewInternalError("Internal error"))
		req, _ := http.NewRequest(http.MethodGet, "/harvests", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHarvestHandler_GetHarvestByID(t *testing.T) {
	r, h, uc, ids, mocks := HarvestHandlerSetUp(t)

	r.GET("/harvests/:id", h.GetHarvestByID)

	t.Run("should return harvest by id successfully", func(t *testing.T) {
		uc.EXPECT().GetHarvestByID(gomock.Any(), ids.HarvestID).Return(mocks.Harvest, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/harvests/"+ids.HarvestID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, response.Data.Quantity, mocks.Harvest.Quantity)
		assert.Equal(t, response.Data.Unit, mocks.Harvest.Unit)
		assert.Equal(t, response.Data.HarvestDate, mocks.Harvest.HarvestDate)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().GetHarvestByID(gomock.Any(), ids.HarvestID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/harvests/"+ids.HarvestID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/harvests/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when not found", func(t *testing.T) {
		uc.EXPECT().GetHarvestByID(gomock.Any(), ids.HarvestID).Return(nil, utils.NewNotFoundError("harvest not found")).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/harvests/"+ids.HarvestID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "NOT_FOUND")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestHarvestHandler_GetHarvestByCommodityID(t *testing.T) {
	r, h, uc, ids, mocks := HarvestHandlerSetUp(t)

	r.GET("/harvests/commodity/:id", h.GetHarvestByCommodityID)

	t.Run("should return harvests by commodity id successfully", func(t *testing.T) {
		uc.EXPECT().GetHarvestByCommodityID(gomock.Any(), ids.CommodityID).Return(mocks.Harvests, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/harvests/commodity/"+ids.CommodityID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var response struct {
			Data []domain.Harvest `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Data, len(*mocks.Harvests))
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().GetHarvestByCommodityID(gomock.Any(), ids.CommodityID).Return(nil, utils.NewInternalError("Internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/harvests/commodity/"+ids.CommodityID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/harvests/commodity/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestHarvestHandler_GetHarvestByLandID(t *testing.T) {
	r, h, uc, ids, mocks := HarvestHandlerSetUp(t)

	r.GET("/harvests/land/:id", h.GetHarvestByLandID)

	t.Run("should return harvests by land id successfully", func(t *testing.T) {
		uc.EXPECT().GetHarvestByLandID(gomock.Any(), ids.LandID).Return(mocks.Harvests, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/harvests/land/"+ids.LandID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var response struct {
			Data []domain.Harvest `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Data, len(*mocks.Harvests))
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().GetHarvestByLandID(gomock.Any(), ids.LandID).Return(nil, utils.NewInternalError("Internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/harvests/land/"+ids.LandID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/harvests/land/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestHarvestHandler_GetHarvestByLandCommodityID(t *testing.T) {
	r, h, uc, ids, mocks := HarvestHandlerSetUp(t)

	r.GET("/harvests/land_commodity/:id", h.GetHarvestByLandCommodityID)

	t.Run("should return harvests by land commodity id successfully", func(t *testing.T) {
		uc.EXPECT().GetHarvestByLandCommodityID(gomock.Any(), ids.LandCommodityID).Return(mocks.Harvests, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/harvests/land_commodity/"+ids.LandCommodityID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var response struct {
			Data []domain.Harvest `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Data, len(*mocks.Harvests))
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().GetHarvestByLandCommodityID(gomock.Any(), ids.LandCommodityID).Return(nil, utils.NewInternalError("Internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/harvests/land_commodity/"+ids.LandCommodityID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/harvests/land_commodity/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestHarvestHandler_GetHarvestByRegionID(t *testing.T) {
	r, h, uc, ids, mocks := HarvestHandlerSetUp(t)

	r.GET("/harvests/region/:id", h.GetHarvestByRegionID)

	t.Run("should return harvests by region id successfully", func(t *testing.T) {
		uc.EXPECT().GetHarvestByRegionID(gomock.Any(), ids.RegionID).Return(mocks.Harvests, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/harvests/region/"+ids.RegionID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var response struct {
			Data []domain.Harvest `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Data, len(*mocks.Harvests))
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().GetHarvestByRegionID(gomock.Any(), ids.RegionID).Return(nil, utils.NewInternalError("Internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/harvests/region/"+ids.RegionID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/harvests/region/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestHarvestHandler_UpdateHarvest(t *testing.T) {
	r, h, uc, ids, mocks := HarvestHandlerSetUp(t)

	r.PATCH("/harvests/:id", h.UpdateHarvest)

	t.Run("should update harvest successfully", func(t *testing.T) {
		uc.EXPECT().UpdateHarvest(gomock.Any(), ids.HarvestID, gomock.Any()).Return(mocks.Harvest, nil).Times(1)

		reqBody := `{"quantity":100,"unit":"kg","harvest_date":"2022-01-01"}`
		req, _ := http.NewRequest(http.MethodPatch, "/harvests/"+ids.HarvestID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, mocks.Harvest.Quantity, response.Data.Quantity)
		assert.Equal(t, mocks.Harvest.Unit, response.Data.Unit)
		assert.Equal(t, mocks.Harvest.HarvestDate, response.Data.HarvestDate)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().UpdateHarvest(gomock.Any(), ids.HarvestID, gomock.Any()).Return(nil, utils.NewInternalError("internal error"))

		reqBody := `{"quantity":100,"unit":"kg","harvest_date":"2022-01-01"}`
		req, _ := http.NewRequest(http.MethodPatch, "/harvests/"+ids.HarvestID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		uc.EXPECT().UpdateHarvest(gomock.Any(), ids.HarvestID, gomock.Any()).Times(0)

		req, _ := http.NewRequest(http.MethodPatch, "/harvests/"+ids.HarvestID.String(), bytes.NewReader([]byte(`invalid-json`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, "/harvests/abc", bytes.NewReader([]byte(`invalid-json`)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when not found", func(t *testing.T) {
		uc.EXPECT().UpdateHarvest(gomock.Any(), ids.HarvestID, gomock.Any()).Return(nil, utils.NewNotFoundError("harvest not found"))

		reqBody := `{"quantity":100,"unit":"kg","harvest_date":"2022-01-01"}`
		req, _ := http.NewRequest(http.MethodPatch, "/harvests/"+ids.HarvestID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "NOT_FOUND")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestHarvestHandler_DeleteHarvest(t *testing.T) {
	r, h, uc, ids, _ := HarvestHandlerSetUp(t)

	r.DELETE("/harvests/:id", h.DeleteHarvest)

	t.Run("should delete harvest successfully", func(t *testing.T) {
		uc.EXPECT().DeleteHarvest(gomock.Any(), ids.HarvestID).Return(nil).Times(1)

		req, _ := http.NewRequest(http.MethodDelete, "/harvests/"+ids.HarvestID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Harvest deleted successfully", response.Data.Message)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().DeleteHarvest(gomock.Any(), ids.HarvestID).Return(utils.NewInternalError("internal error"))

		req, _ := http.NewRequest(http.MethodDelete, "/harvests/"+ids.HarvestID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/harvests/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when not found", func(t *testing.T) {
		uc.EXPECT().DeleteHarvest(gomock.Any(), ids.HarvestID).Return(utils.NewNotFoundError("harvest not found"))

		req, _ := http.NewRequest(http.MethodDelete, "/harvests/"+ids.HarvestID.String(), nil)
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

func TestHarvestHandler_RestoreHarvest(t *testing.T) {
	r, h, uc, ids, mocks := HarvestHandlerSetUp(t)

	r.POST("/harvests/:id/restore", h.RestoreHarvest)

	t.Run("should restore harvest successfully", func(t *testing.T) {
		uc.EXPECT().RestoreHarvest(gomock.Any(), ids.HarvestID).Return(mocks.Harvest, nil).Times(1)

		req, _ := http.NewRequest(http.MethodPost, "/harvests/"+ids.HarvestID.String()+"/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, mocks.Harvest.ID, response.Data.ID)
		assert.Equal(t, mocks.Harvest.Quantity, response.Data.Quantity)
		assert.Equal(t, mocks.Harvest.Unit, response.Data.Unit)
		assert.Equal(t, mocks.Harvest.HarvestDate, response.Data.HarvestDate)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().RestoreHarvest(gomock.Any(), ids.HarvestID).Return(nil, utils.NewInternalError("Internal error"))

		req, _ := http.NewRequest(http.MethodPost, "/harvests/"+ids.HarvestID.String()+"/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/harvests/abc/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when not found", func(t *testing.T) {
		uc.EXPECT().RestoreHarvest(gomock.Any(), ids.HarvestID).Return(nil, utils.NewNotFoundError("harvest not found"))

		req, _ := http.NewRequest(http.MethodPost, "/harvests/"+ids.HarvestID.String()+"/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "NOT_FOUND")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestHarvestHandler_GetAllDeletedHarvest(t *testing.T) {
	r, h, uc, _, mocks := HarvestHandlerSetUp(t)
	r.GET("/harvests/deleted", h.GetAllDeletedHarvest)
	t.Run("should return all deleted harvests successfully", func(t *testing.T) {
		uc.EXPECT().GetAllDeletedHarvest(gomock.Any()).Return(mocks.Harvests, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/harvests/deleted", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var response struct {
			Data []domain.Harvest `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Data, len(*mocks.Harvests))
	})
	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().GetAllDeletedHarvest(gomock.Any()).Return(nil, utils.NewInternalError("Internal error")).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/harvests/deleted", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHarvestHandler_GetHarvestDeletedByID(t *testing.T) {
	r, h, uc, ids, mocks := HarvestHandlerSetUp(t)

	r.GET("/harvests/deleted/:id", h.GetHarvestDeletedByID)

	t.Run("should return deleted harvest by id successfully", func(t *testing.T) {
		uc.EXPECT().GetHarvestDeletedByID(gomock.Any(), ids.HarvestID).Return(mocks.Harvest, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/harvests/deleted/"+ids.HarvestID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, response.Data.Quantity, mocks.Harvest.Quantity)
		assert.Equal(t, response.Data.Unit, mocks.Harvest.Unit)
		assert.Equal(t, response.Data.HarvestDate, mocks.Harvest.HarvestDate)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().GetHarvestDeletedByID(gomock.Any(), ids.HarvestID).Return(nil, utils.NewInternalError("Internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/harvests/deleted/"+ids.HarvestID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/harvests/deleted/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when not found", func(t *testing.T) {
		uc.EXPECT().GetHarvestDeletedByID(gomock.Any(), ids.HarvestID).Return(nil, utils.NewNotFoundError("harvest not found")).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/harvests/deleted/"+ids.HarvestID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseHarvestHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "NOT_FOUND")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
