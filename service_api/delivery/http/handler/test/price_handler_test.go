package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	handler_implementation "github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler/implementation"
	handler_interface "github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler/interface"
	"github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler/test/response"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	"github.com/ryvasa/go-super-farmer/service_api/model/dto"
	mock_usecase "github.com/ryvasa/go-super-farmer/service_api/usecase/mock"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
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
	Prices       []*domain.Price
	UpdatePrice  *domain.Price
	PriceHistory []*domain.PriceHistory
}

type PriceHandlerIDs struct {
	PriceID     uuid.UUID
	CommodityID uuid.UUID
	CityID      int64
}

type PriceHandlerDTOMocks struct {
	ParamsDTO           *dto.PriceParamsDTO
	ResponseDownloadDTO *dto.DownloadResponseDTO
}

func PriceHandlerSetUp(t *testing.T) (*gin.Engine, handler_interface.PriceHandler, *mock_usecase.MockPriceUsecase, PriceHandlerIDs, PriceHandlerMocks, PriceHandlerDTOMocks) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc := mock_usecase.NewMockPriceUsecase(ctrl)
	h := handler_implementation.NewPriceHandler(uc)
	r := gin.Default()

	priceID := uuid.New()
	commodityID := uuid.New()
	cityID := int64(1)

	ids := PriceHandlerIDs{
		PriceID:     priceID,
		CommodityID: commodityID,
		CityID:      cityID,
	}

	mocks := PriceHandlerMocks{
		Price: &domain.Price{
			ID:          priceID,
			CommodityID: commodityID,
			CityID:      cityID,
		},
		Prices: []*domain.Price{
			{
				ID:          priceID,
				CommodityID: commodityID,
				CityID:      cityID,
			},
		},
		PriceHistory: []*domain.PriceHistory{
			{
				ID:          priceID,
				CommodityID: commodityID,
				CityID:      cityID,
			},
		},
		UpdatePrice: &domain.Price{
			ID:          priceID,
			CommodityID: commodityID,
			CityID:      cityID,
		},
	}

	startDate, _ := time.Parse("2006-01-02", "2023-10-26")
	endDate, _ := time.Parse("2006-01-02", "2023-10-27")

	dtos := PriceHandlerDTOMocks{
		ParamsDTO: &dto.PriceParamsDTO{
			CommodityID: commodityID,
			CityID:      cityID,
			StartDate:   startDate,
			EndDate:     endDate,
		},
		ResponseDownloadDTO: &dto.DownloadResponseDTO{
			Message:     "Price history report generation in progress. Please check back in a few moments.",
			DownloadURL: "http://localhost:8080/api/prices/history/commodity/1/city/1/download/file?start_date=2023-10-26&end_date=2023-10-27",
		},
	}

	return r, h, uc, ids, mocks, dtos
}

func TestPriceHandler_CreatePrice(t *testing.T) {
	r, h, uc, ids, mockc, _ := PriceHandlerSetUp(t)

	r.POST("/prices", h.CreatePrice)

	t.Run("should create price successfully", func(t *testing.T) {
		uc.EXPECT().CreatePrice(gomock.Any(), gomock.Any()).Return(mockc.Price, nil).Times(1)

		reqBody := `{"commodity_id":"` + ids.CommodityID.String() + `","city_id":` + strconv.FormatInt(ids.CityID, 10) + `}`
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

	t.Run("should return error when service_api error", func(t *testing.T) {
		uc.EXPECT().CreatePrice(gomock.Any(), gomock.Any()).Return(nil, utils.NewInternalError("Internal error"))

		reqBody := `{"commodity_id":"` + ids.CommodityID.String() + `","city_id":` + strconv.FormatInt(ids.CityID, 10) + `}`
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
	r, h, uc, _, mocks, _ := PriceHandlerSetUp(t)

	r.GET("/prices", h.GetAllPrices)

	t.Run("should return all prices successfully with default pagination", func(t *testing.T) {
		// Prepare expected response
		expectedResponse := &dto.PaginationResponseDTO{
			TotalRows:  1,
			TotalPages: 1,
			Page:       1,
			Limit:      10,
			Data:       mocks.Prices,
		}

		// Setup mock
		uc.EXPECT().
			GetAllPrices(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx *gin.Context, p *dto.PaginationDTO) (*dto.PaginationResponseDTO, error) {
				// Verify default pagination
				assert.Equal(t, 1, p.Page)
				assert.Equal(t, 10, p.Limit)
				assert.Equal(t, "created_at desc", p.Sort)
				return expectedResponse, nil
			})

		// Make request
		req, _ := http.NewRequest(http.MethodGet, "/prices", nil)
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

	t.Run("should return prices with custom pagination and filter", func(t *testing.T) {
		expectedResponse := &dto.PaginationResponseDTO{
			TotalRows:  1,
			TotalPages: 1,
			Page:       2,
			Limit:      5,
			Data:       mocks.Prices,
		}

		// Setup mock
		uc.EXPECT().
			GetAllPrices(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx *gin.Context, p *dto.PaginationDTO) (*dto.PaginationResponseDTO, error) {
				assert.Equal(t, 2, p.Page)
				assert.Equal(t, 5, p.Limit)
				return expectedResponse, nil
			})

		// Make request with query params
		req, _ := http.NewRequest(http.MethodGet, "/prices?page=2&limit=5", nil)
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
		req, _ := http.NewRequest(http.MethodGet, "/prices?page=-1", nil)
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
		req, _ := http.NewRequest(http.MethodGet, "/prices?limit=101", nil)
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
		uc.EXPECT().
			GetAllPrices(gomock.Any(), gomock.Any()).
			Return(nil, utils.NewInternalError("service_api error"))

		req, _ := http.NewRequest(http.MethodGet, "/prices", nil)
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

func TestPriceHandler_GetPriceByID(t *testing.T) {
	r, h, uc, ids, mocks, _ := PriceHandlerSetUp(t)

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

	t.Run("should return error when service_api error", func(t *testing.T) {
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
	r, h, uc, ids, mocks, _ := PriceHandlerSetUp(t)

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
		assert.Equal(t, len(mocks.Prices), len(response.Data))
		assert.Equal(t, response.Data[0].ID, (mocks.Prices)[0].ID)
	})

	t.Run("should return error when service_api error", func(t *testing.T) {
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

func TestPriceHandler_GetPricesByCityID(t *testing.T) {
	r, h, uc, ids, mocks, _ := PriceHandlerSetUp(t)

	r.GET("/prices/city/:id", h.GetPricesByCityID)

	t.Run("should get prices by city id successfully", func(t *testing.T) {
		uc.EXPECT().GetPricesByCityID(gomock.Any(), ids.CityID).Return(mocks.Prices, nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/prices/city/"+strconv.FormatInt(ids.CityID, 10), nil))

		var response responsePricesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, len(mocks.Prices), len(response.Data))
		assert.Equal(t, response.Data[0].ID, (mocks.Prices)[0].ID)
	})

	t.Run("should return error when service_api error", func(t *testing.T) {
		uc.EXPECT().GetPricesByCityID(gomock.Any(), ids.CityID).Return(nil, utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/prices/city/"+strconv.FormatInt(ids.CityID, 10), nil))

		var response responsePricesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when id is invalid", func(t *testing.T) {
		uc.EXPECT().GetPricesByCityID(gomock.Any(), uuid.Nil).Return(nil, utils.NewBadRequestError("ID is invalid"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/prices/city/aa", nil))

		var response responsePricesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPriceHandler_UpdatePrice(t *testing.T) {
	r, h, uc, ids, mocks, _ := PriceHandlerSetUp(t)

	r.PUT("/prices/:id", h.UpdatePrice)

	t.Run("should update price successfully", func(t *testing.T) {
		uc.EXPECT().UpdatePrice(gomock.Any(), ids.PriceID, gomock.Any()).Return(mocks.UpdatePrice, nil).Times(1)

		reqBody := `{"commodity_id":"` + ids.CommodityID.String() + `","city_id":` + strconv.FormatInt(ids.CityID, 10) + `}`
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

	t.Run("should return error when service_api error", func(t *testing.T) {
		uc.EXPECT().UpdatePrice(gomock.Any(), ids.PriceID, gomock.Any()).Return(nil, utils.NewInternalError("Internal error"))

		reqBody := `{"commodity_id":"` + ids.CommodityID.String() + `","city_id":` + strconv.FormatInt(ids.CityID, 10) + `}`
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

		reqBody := `{"commodity_id":"` + ids.CommodityID.String() + `","city_id":` + strconv.FormatInt(ids.CityID, 10) + `}`
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
	r, h, uc, ids, _, _ := PriceHandlerSetUp(t)

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

	t.Run("should return error when service_api error", func(t *testing.T) {
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
	r, h, uc, ids, mocks, _ := PriceHandlerSetUp(t)

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

	t.Run("should return error when service_api error", func(t *testing.T) {
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

func TestPriceHandler_GetPriceByCommodityIDAndCityID(t *testing.T) {
	r, h, uc, ids, mocks, _ := PriceHandlerSetUp(t)

	r.GET("/prices/commodity_id/:commodity_id/city/:city_id", h.GetPriceByCommodityIDAndCityID)

	t.Run("should get price by commodity id and city id successfully", func(t *testing.T) {
		uc.EXPECT().GetPriceByCommodityIDAndCityID(gomock.Any(), ids.CommodityID, ids.CityID).Return(mocks.Price, nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/prices/commodity_id/"+ids.CommodityID.String()+"/city/"+strconv.FormatInt(ids.CityID, 10), nil))

		var response responsePriceHandler
		assert.NoError(t, json.NewDecoder(w.Body).Decode(&response))

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, response.Data.ID, mocks.Price.ID)
	})

	t.Run("should return error when service_api error", func(t *testing.T) {
		uc.EXPECT().GetPriceByCommodityIDAndCityID(gomock.Any(), ids.CommodityID, ids.CityID).Return(nil, utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/prices/commodity_id/"+ids.CommodityID.String()+"/city/"+strconv.FormatInt(ids.CityID, 10), nil))

		var response responsePriceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when commodity id is invalid", func(t *testing.T) {

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, fmt.Sprintf("/prices/commodity_id/aa/city/%d", ids.CityID), nil))

		var response responsePriceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when city id is invalid", func(t *testing.T) {

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, fmt.Sprintf("/prices/commodity_id/%s/city/qq", ids.CommodityID), nil))

		var response responsePriceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPriceHandler_GetPricesHistoryByCommodityIDAndCityID(t *testing.T) {
	r, h, uc, ids, mocks, _ := PriceHandlerSetUp(t)

	r.GET("/prices/commodity_id/:commodity_id/city/:city_id/history", h.GetPricesHistoryByCommodityIDAndCityID)

	t.Run("should get price history by commodity id and city id successfully", func(t *testing.T) {
		uc.EXPECT().GetPriceHistoryByCommodityIDAndCityID(gomock.Any(), ids.CommodityID, ids.CityID).Return(mocks.PriceHistory, nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/prices/commodity_id/"+ids.CommodityID.String()+"/city/"+strconv.FormatInt(ids.CityID, 10)+"/history", nil))

		var response responsePricesHandler
		assert.NoError(t, json.NewDecoder(w.Body).Decode(&response))

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
	})

	t.Run("should return error when service_api error", func(t *testing.T) {
		uc.EXPECT().GetPriceHistoryByCommodityIDAndCityID(gomock.Any(), ids.CommodityID, ids.CityID).Return(nil, utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/prices/commodity_id/"+ids.CommodityID.String()+"/city/"+strconv.FormatInt(ids.CityID, 10)+"/history", nil))

		var response responsePricesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when commodity id is invalid", func(t *testing.T) {
		uc.EXPECT().GetPriceHistoryByCommodityIDAndCityID(gomock.Any(), uuid.Nil, uuid.Nil).Return(nil, utils.NewBadRequestError("ID is invalid"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, fmt.Sprintf("/prices/commodity_id/aa/city/%d/history", ids.CityID), nil))

		var response responsePricesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when city id is invalid", func(t *testing.T) {
		uc.EXPECT().GetPriceHistoryByCommodityIDAndCityID(gomock.Any(), uuid.Nil, uuid.Nil).Return(nil, utils.NewBadRequestError("ID is invalid"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, fmt.Sprintf("/prices/commodity_id/%s/city/bb/history", ids.CommodityID), nil))

		var response responsePricesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPriceHandler_DownloadPriceByLandCommodityID(t *testing.T) {
	r, h, uc, ids, _, dtos := PriceHandlerSetUp(t) // No usecase mocking needed here either.
	r.GET("/prices/history/commodity/:commodity_id/city/:city_id/download", h.DownloadPricesHistoryByCommodityIDAndCityID)

	t.Run("should return success response and download URL", func(t *testing.T) {
		uc.EXPECT().DownloadPriceHistoryByCommodityIDAndCityID(gomock.Any(), dtos.ParamsDTO).Return(dtos.ResponseDownloadDTO, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/prices/history/commodity/%s/city/%d/download?start_date=2023-10-26&end_date=2023-10-27", ids.CommodityID, ids.CityID), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response response.ResponseDownload
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.Data.DownloadURL)
		assert.Equal(t, "Price history report generation in progress. Please check back in a few moments.", response.Data.Message)
	})

	t.Run("should return error when invalid commodity id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/prices/history/commodity/abc/city/%d/download", ids.CityID), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response response.ResponseDownload
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when invalid city id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/prices/history/commodity/%s/city/abc/download", ids.CommodityID), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response response.ResponseDownload
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when invalid start date format", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/prices/history/commodity/%s/city/%d/download?start_date=invalid-date&end_date=2023-10-27", ids.CommodityID, ids.CityID), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response response.ResponseDownload
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when invalid end date format", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/prices/history/commodity/%s/city/%d/download?start_date=2023-10-26&end_date=invalid-date", ids.CommodityID, ids.CityID), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response response.ResponseDownload
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when land commodity usecase returns error", func(t *testing.T) {
		uc.EXPECT().DownloadPriceHistoryByCommodityIDAndCityID(gomock.Any(), dtos.ParamsDTO).Return(nil, utils.NewInternalError("Internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/prices/history/commodity/%s/city/%d/download?start_date=2023-10-26&end_date=2023-10-27", ids.CommodityID, ids.CityID), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response response.ResponseDownload
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPriceHandler_GetPriceExcelFile(t *testing.T) {
	r, h, uc, ids, _, dtos := PriceHandlerSetUp(t)
	r.GET("/prices/history/commodity/:commodity_id/city/:city_id/download/file", h.GetPriceHistoryExcelFile)

	// Create a temporary directory for test files.  This is crucial for cleanup.
	tempDir, err := os.MkdirTemp("", "price_reports")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after the test

	// Create a dummy Excel file (replace with your actual file creation if needed)
	dummyFilePath := filepath.Join(tempDir, "prices_dummy.xlsx")
	err = os.WriteFile(dummyFilePath, []byte("Dummy Excel content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create dummy Excel file: %v", err)
	}

	t.Run("should return excel file", func(t *testing.T) {
		// Modify the expectation to match the dummy file we just created.  Important!
		// We need to return the correct file path.
		expectedFilePath := dummyFilePath
		logrus.Log.Info(expectedFilePath)
		uc.EXPECT().GetPriceExcelFile(gomock.Any(), dtos.ParamsDTO).Return(&expectedFilePath, nil).Times(1)

		logrus.Log.Info(ids.CommodityID, ids.CityID)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/prices/history/commodity/%s/city/%d/download/file?start_date=2023-10-26&end_date=2023-10-27", ids.CommodityID, ids.CityID), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response response.ResponseDownload
		err := json.Unmarshal(w.Body.Bytes(), &response)
		logrus.Log.Info(response.Data.DownloadURL)
		logrus.Log.Info(err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Disposition"), "filename=prices_dummy.xlsx")                              // Check filename in header
		assert.Equal(t, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", w.Header().Get("Content-Type")) //Check content type
	})

	t.Run("should return error when invalid commodity id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/prices/history/commodity/abc/city/%d/download/file?start_date=2023-10-26&end_date=2023-10-27", ids.CityID), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response response.ResponseDownload
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when invalid city id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/prices/history/commodity/%s/city/abc/download/file?start_date=2023-10-26&end_date=2023-10-27", ids.CommodityID), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response response.ResponseDownload
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when invalid start date format", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/prices/history/commodity/%s/city/%d/download/file?start_date=invalid-date&end_date=2023-10-27", ids.CommodityID, ids.CityID), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response response.ResponseDownload
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when invalid end date format", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/prices/history/commodity/%s/city/%d/download/file?start_date=2023-10-26&end_date=invalid-date", ids.CommodityID, ids.CityID), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response response.ResponseDownload
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return 404 when file not found", func(t *testing.T) {
		// This case now has to be modified to reflect that a file is NOT present
		uc.EXPECT().GetPriceExcelFile(gomock.Any(), dtos.ParamsDTO).Return(nil, utils.NewNotFoundError("Report file not found")).Times(1)
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/prices/history/commodity/%s/city/%d/download/file?start_date=2023-10-26&end_date=2023-10-27", ids.CommodityID, ids.CityID), nil) // Use existing ID
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var response responsePriceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "NOT_FOUND")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should return 500 when usecase returns an error", func(t *testing.T) {
		uc.EXPECT().GetPriceExcelFile(gomock.Any(), dtos.ParamsDTO).Return(nil, utils.NewInternalError("Simulated file system error")).Times(1)
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/prices/history/commodity/%s/city/%d/download/file?start_date=2023-10-26&end_date=2023-10-27", ids.CommodityID, ids.CityID), nil) // Use existing ID
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
