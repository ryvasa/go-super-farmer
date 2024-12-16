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

type responsePriceHandler struct {
	Status  int            `json:"status"`
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    domain.Price   `json:"data"`
	Errors  response.Error `json:"errors"`
}

type responsePricesHandler struct {
	Status  int            `json:"status"`
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    []domain.Price `json:"data"`
	Errors  response.Error `json:"errors"`
}

type PriceHandlerMocks struct {
	Price        *domain.Price
	Prices       *[]domain.Price
	UpdatePrice  *domain.Price
	PriceHistory *[]domain.PriceHistory
}

type PriceHandlerIDs struct {
	PriceID     uuid.UUID
	CommodityID uuid.UUID
	RegionID    uuid.UUID
}

func PriceHandlerSetUp(t *testing.T) (*gin.Engine, handler_interface.PriceHandler, *mock.MockPriceUsecase, PriceHandlerIDs, PriceHandlerMocks) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc := mock.NewMockPriceUsecase(ctrl)
	h := handler_implementation.NewPriceHandler(uc)
	r := gin.Default()

	priceID := uuid.New()
	CommodityID := uuid.New()
	RegionID := uuid.New()

	ids := PriceHandlerIDs{
		PriceID:     priceID,
		CommodityID: CommodityID,
		RegionID:    RegionID,
	}

	mocks := PriceHandlerMocks{
		Price: &domain.Price{
			ID:          priceID,
			CommodityID: CommodityID,
			RegionID:    RegionID,
		},
		Prices: &[]domain.Price{
			{
				ID:          priceID,
				CommodityID: CommodityID,
				RegionID:    RegionID,
			},
		},
		PriceHistory: &[]domain.PriceHistory{
			{
				ID:          priceID,
				CommodityID: CommodityID,
				RegionID:    RegionID,
			},
		},
		UpdatePrice: &domain.Price{
			ID:          priceID,
			CommodityID: CommodityID,
			RegionID:    RegionID,
		},
	}
	return r, h, uc, ids, mocks
}

func TestPriceHandler_CreatePrice(t *testing.T) {
	r, h, uc, ids, mockc := PriceHandlerSetUp(t)

	r.POST("/prices", h.CreatePrice)

	t.Run("should create price successfully", func(t *testing.T) {
		uc.EXPECT().CreatePrice(gomock.Any(), gomock.Any()).Return(mockc.Price, nil).Times(1)

		reqBody := `{"commodity_id":"` + ids.CommodityID.String() + `","region_id":"` + ids.RegionID.String() + `"}`
		req, _ := http.NewRequest(http.MethodPost, "/prices", bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responsePriceHandler
		assert.NoError(t, json.NewDecoder(w.Body).Decode(&response))

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, response.Data.ID, mockc.Price.ID)
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		uc.EXPECT().CreatePrice(gomock.Any(), gomock.Any()).Times(0)

		req, _ := http.NewRequest(http.MethodPost, "/prices", bytes.NewReader([]byte("invalid-json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responsePriceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().CreatePrice(gomock.Any(), gomock.Any()).Return(nil, utils.NewInternalError("Internal error"))

		reqBody := `{"commodity_id":"` + ids.CommodityID.String() + `","region_id":"` + ids.RegionID.String() + `"}`
		req, _ := http.NewRequest(http.MethodPost, "/prices", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responsePriceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPriceHandler_GetAllPrices(t *testing.T) {
	r, h, uc, _, mocks := PriceHandlerSetUp(t)

	r.GET("/prices", h.GetAllPrices)

	t.Run("should get all prices successfully", func(t *testing.T) {
		uc.EXPECT().GetAllPrices(gomock.Any()).Return(mocks.Prices, nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/prices", nil))

		var response responsePricesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, len(*mocks.Prices), len(response.Data))
		assert.Equal(t, response.Data[0].ID, (*mocks.Prices)[0].ID)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().GetAllPrices(gomock.Any()).Return(nil, utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/prices", nil))

		var response responsePricesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPriceHandler_GetPriceByID(t *testing.T) {
	r, h, uc, ids, mocks := PriceHandlerSetUp(t)

	r.GET("/prices/:id", h.GetPriceByID)

	t.Run("should get price by id successfully", func(t *testing.T) {
		uc.EXPECT().GetPriceByID(gomock.Any(), ids.PriceID).Return(mocks.Price, nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/prices/"+ids.PriceID.String(), nil))

		var response responsePriceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, response.Data.ID, mocks.Price.ID)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().GetPriceByID(gomock.Any(), ids.PriceID).Return(nil, utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/prices/"+ids.PriceID.String(), nil))

		var response responsePriceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when id is invalid", func(t *testing.T) {
		uc.EXPECT().GetPriceByID(gomock.Any(), uuid.Nil).Return(nil, utils.NewBadRequestError("ID is invalid"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/prices/aa", nil))

		var response responsePriceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPriceHandler_GetPricesByCommodityID(t *testing.T) {
	r, h, uc, ids, mocks := PriceHandlerSetUp(t)

	r.GET("/prices/commodity_id/:id", h.GetPricesByCommodityID)

	t.Run("should get prices by commodity id successfully", func(t *testing.T) {
		uc.EXPECT().GetPricesByCommodityID(gomock.Any(), ids.CommodityID).Return(mocks.Prices, nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/prices/commodity_id/"+ids.CommodityID.String(), nil))

		var response responsePricesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, len(*mocks.Prices), len(response.Data))
		assert.Equal(t, response.Data[0].ID, (*mocks.Prices)[0].ID)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().GetPricesByCommodityID(gomock.Any(), ids.CommodityID).Return(nil, utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/prices/commodity_id/"+ids.CommodityID.String(), nil))

		var response responsePricesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when id is invalid", func(t *testing.T) {
		uc.EXPECT().GetPricesByCommodityID(gomock.Any(), uuid.Nil).Return(nil, utils.NewBadRequestError("ID is invalid"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/prices/commodity_id/aa", nil))

		var response responsePricesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPriceHandler_GetPricesByRegionID(t *testing.T) {
	r, h, uc, ids, mocks := PriceHandlerSetUp(t)

	r.GET("/prices/region/:id", h.GetPricesByRegionID)

	t.Run("should get prices by region id successfully", func(t *testing.T) {
		uc.EXPECT().GetPricesByRegionID(gomock.Any(), ids.RegionID).Return(mocks.Prices, nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/prices/region/"+ids.RegionID.String(), nil))

		var response responsePricesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, len(*mocks.Prices), len(response.Data))
		assert.Equal(t, response.Data[0].ID, (*mocks.Prices)[0].ID)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().GetPricesByRegionID(gomock.Any(), ids.RegionID).Return(nil, utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/prices/region/"+ids.RegionID.String(), nil))

		var response responsePricesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when id is invalid", func(t *testing.T) {
		uc.EXPECT().GetPricesByRegionID(gomock.Any(), uuid.Nil).Return(nil, utils.NewBadRequestError("ID is invalid"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/prices/region/aa", nil))

		var response responsePricesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPriceHandler_UpdatePrice(t *testing.T) {
	r, h, uc, ids, mocks := PriceHandlerSetUp(t)

	r.PUT("/prices/:id", h.UpdatePrice)

	t.Run("should update price successfully", func(t *testing.T) {
		uc.EXPECT().UpdatePrice(gomock.Any(), ids.PriceID, gomock.Any()).Return(mocks.UpdatePrice, nil).Times(1)

		reqBody := `{"commodity_id":"` + ids.CommodityID.String() + `","region_id":"` + ids.RegionID.String() + `"}`
		req, _ := http.NewRequest(http.MethodPut, "/prices/"+ids.PriceID.String(), bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responsePriceHandler
		assert.NoError(t, json.NewDecoder(w.Body).Decode(&response))

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, response.Data.ID, mocks.UpdatePrice.ID)
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		uc.EXPECT().UpdatePrice(gomock.Any(), ids.PriceID, gomock.Any()).Times(0)

		req, _ := http.NewRequest(http.MethodPut, "/prices/"+ids.PriceID.String(), bytes.NewReader([]byte("invalid-json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responsePriceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().UpdatePrice(gomock.Any(), ids.PriceID, gomock.Any()).Return(nil, utils.NewInternalError("Internal error"))

		reqBody := `{"commodity_id":"` + ids.CommodityID.String() + `","region_id":"` + ids.RegionID.String() + `"}`
		req, _ := http.NewRequest(http.MethodPut, "/prices/"+ids.PriceID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responsePriceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when id is invalid", func(t *testing.T) {
		uc.EXPECT().UpdatePrice(gomock.Any(), uuid.Nil, gomock.Any()).Return(nil, utils.NewBadRequestError("ID is invalid"))

		reqBody := `{"commodity_id":"` + ids.CommodityID.String() + `","region_id":"` + ids.RegionID.String() + `"}`
		req, _ := http.NewRequest(http.MethodPut, "/prices/aa", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responsePriceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPriceHandler_DeletePrice(t *testing.T) {
	r, h, uc, ids, _ := PriceHandlerSetUp(t)

	r.DELETE("/prices/:id", h.DeletePrice)

	t.Run("should delete price successfully", func(t *testing.T) {
		uc.EXPECT().DeletePrice(gomock.Any(), ids.PriceID).Return(nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/prices/"+ids.PriceID.String(), nil))

		var response response.ResponseMessage
		assert.NoError(t, json.NewDecoder(w.Body).Decode(&response))

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, response.Data.Message, "Price deleted successfully")
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().DeletePrice(gomock.Any(), ids.PriceID).Return(utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/prices/"+ids.PriceID.String(), nil))

		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when id is invalid", func(t *testing.T) {
		uc.EXPECT().DeletePrice(gomock.Any(), uuid.Nil).Return(utils.NewBadRequestError("ID is invalid"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/prices/aa", nil))

		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPriceHandler_RestorePrice(t *testing.T) {
	r, h, uc, ids, mocks := PriceHandlerSetUp(t)

	r.PATCH("/prices/:id/restore", h.RestorePrice)

	t.Run("should restore price successfully", func(t *testing.T) {
		uc.EXPECT().RestorePrice(gomock.Any(), ids.PriceID).Return(mocks.Price, nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodPatch, "/prices/"+ids.PriceID.String()+"/restore", nil))

		var response responsePriceHandler
		assert.NoError(t, json.NewDecoder(w.Body).Decode(&response))

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, response.Data.ID, ids.PriceID)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().RestorePrice(gomock.Any(), ids.PriceID).Return(nil, utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodPatch, "/prices/"+ids.PriceID.String()+"/restore", nil))

		var response responsePriceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when id is invalid", func(t *testing.T) {
		uc.EXPECT().RestorePrice(gomock.Any(), uuid.Nil).Return(nil, utils.NewBadRequestError("ID is invalid"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodPatch, "/prices/aa/restore", nil))

		var response responsePriceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPriceHandler_GetPriceByCommodityIDAndRegionID(t *testing.T) {
	r, h, uc, ids, mocks := PriceHandlerSetUp(t)

	r.GET("/prices/commodity_id/:commodity_id/region/:region_id", h.GetPriceByCommodityIDAndRegionID)

	t.Run("should get price by commodity id and region id successfully", func(t *testing.T) {
		uc.EXPECT().GetPriceByCommodityIDAndRegionID(gomock.Any(), ids.CommodityID, ids.RegionID).Return(mocks.Price, nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/prices/commodity_id/"+ids.CommodityID.String()+"/region/"+ids.RegionID.String(), nil))

		var response responsePriceHandler
		assert.NoError(t, json.NewDecoder(w.Body).Decode(&response))

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, response.Data.ID, mocks.Price.ID)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().GetPriceByCommodityIDAndRegionID(gomock.Any(), ids.CommodityID, ids.RegionID).Return(nil, utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/prices/commodity_id/"+ids.CommodityID.String()+"/region/"+ids.RegionID.String(), nil))

		var response responsePriceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when id is invalid", func(t *testing.T) {

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/prices/commodity_id/aa/region/bb", nil))

		var response responsePriceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPriceHandler_GetPricesHistoryByCommodityIDAndRegionID(t *testing.T) {
	r, h, uc, ids, mocks := PriceHandlerSetUp(t)

	r.GET("/prices/commodity_id/:commodity_id/region/:region_id/history", h.GetPricesHistoryByCommodityIDAndRegionID)

	t.Run("should get price history by commodity id and region id successfully", func(t *testing.T) {
		uc.EXPECT().GetPriceHistoryByCommodityIDAndRegionID(gomock.Any(), ids.CommodityID, ids.RegionID).Return(mocks.PriceHistory, nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/prices/commodity_id/"+ids.CommodityID.String()+"/region/"+ids.RegionID.String()+"/history", nil))

		var response responsePricesHandler
		assert.NoError(t, json.NewDecoder(w.Body).Decode(&response))

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().GetPriceHistoryByCommodityIDAndRegionID(gomock.Any(), ids.CommodityID, ids.RegionID).Return(nil, utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/prices/commodity_id/"+ids.CommodityID.String()+"/region/"+ids.RegionID.String()+"/history", nil))

		var response responsePricesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when id is invalid", func(t *testing.T) {
		uc.EXPECT().GetPriceHistoryByCommodityIDAndRegionID(gomock.Any(), uuid.Nil, uuid.Nil).Return(nil, utils.NewBadRequestError("ID is invalid"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/prices/commodity_id/aa/region/bb/history", nil))

		var response responsePricesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
