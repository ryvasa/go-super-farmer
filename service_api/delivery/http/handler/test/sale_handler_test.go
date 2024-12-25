package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	handler_implementation "github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler/implementation"
	handler_interface "github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler/interface"
	"github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler/test/response"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	"github.com/ryvasa/go-super-farmer/service_api/model/dto"
	mock_usecase "github.com/ryvasa/go-super-farmer/service_api/usecase/mock"
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
)

type responseSaleHandler struct {
	Status  int            `json:"status"`
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    domain.Sale    `json:"data"`
	Errors  response.Error `json:"errors"`
}

type responseSalesHandler struct {
	Status  int                       `json:"status"`
	Success bool                      `json:"success"`
	Message string                    `json:"message"`
	Data    dto.PaginationResponseDTO `json:"data"`
	Errors  response.Error            `json:"errors"`
}

type SaleHandlerDomain struct {
	Sale  *domain.Sale
	Sales []*domain.Sale
}

type SaleHandlerIDs struct {
	SaleID      uuid.UUID
	CommodityID uuid.UUID
	CityID      int64
}

type SaleHandlerUC struct {
	Sale *mock_usecase.MockSaleUsecase
}

func SaleHandlerSetup(t *testing.T) (*gin.Engine, handler_interface.SaleHandler, *SaleHandlerUC, SaleHandlerIDs, SaleHandlerDomain) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ucSale := mock_usecase.NewMockSaleUsecase(ctrl)
	h := handler_implementation.NewSaleHandler(ucSale)
	r := gin.Default()

	ids := SaleHandlerIDs{
		SaleID:      uuid.New(),
		CommodityID: uuid.New(),
		CityID:      1,
	}
	domain := SaleHandlerDomain{
		Sale: &domain.Sale{
			ID:          ids.SaleID,
			CommodityID: ids.CommodityID,
			CityID:      ids.CityID,
			Quantity:    1,
		},
		Sales: []*domain.Sale{
			{
				ID:          ids.SaleID,
				CommodityID: ids.CommodityID,
				CityID:      ids.CityID,
				Quantity:    1,
			},
		},
	}
	uc := &SaleHandlerUC{
		Sale: ucSale,
	}
	return r, h, uc, ids, domain
}

func TestSaleHandler_CreateSale(t *testing.T) {
	r, h, uc, ids, domain := SaleHandlerSetup(t)
	r.POST("/sales", h.CreateSale)

	t.Run("should create sale successfully", func(t *testing.T) {
		uc.Sale.EXPECT().CreateSale(gomock.Any(), gomock.Any()).Return(domain.Sale, nil).Times(1)

		reqBody := `{"commodity_id":"` + ids.CommodityID.String() + `","city_id":1,"quantity":1}`
		req, _ := http.NewRequest(http.MethodPost, "/sales", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseSaleHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, response.Data.ID, ids.SaleID)
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.Sale.EXPECT().CreateSale(gomock.Any(), gomock.Any()).Return(domain.Sale, utils.NewInternalError("service_api error")).Times(1)

		reqBody := `{"commodity_id":"` + ids.CommodityID.String() + `","city_id":1,"quantity":1}`
		req, _ := http.NewRequest(http.MethodPost, "/sales", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseSaleHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		uc.Sale.EXPECT().CreateSale(gomock.Any(), gomock.Any()).Times(0)
		req, _ := http.NewRequest(http.MethodPost, "/sales", bytes.NewReader([]byte(`invalid-json`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseSaleHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestSaleHandler_GetAllSales(t *testing.T) {
	r, h, uc, _, domain := SaleHandlerSetup(t)
	r.GET("/sales", h.GetAllSales)

	t.Run("should return all sales successfully with default pagination", func(t *testing.T) {
		// Prepare expected response
		expectedResponse := &dto.PaginationResponseDTO{
			TotalRows:  1,
			TotalPages: 1,
			Page:       1,
			Limit:      10,
			Data:       domain.Sales,
		}

		// Setup mock
		uc.Sale.EXPECT().
			GetAllSales(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx *gin.Context, p *dto.PaginationDTO) (*dto.PaginationResponseDTO, error) {
				// Verify default pagination
				assert.Equal(t, 1, p.Page)
				assert.Equal(t, 10, p.Limit)
				assert.Equal(t, "created_at desc", p.Sort)
				return expectedResponse, nil
			})

		// Make request
		req, _ := http.NewRequest(http.MethodGet, "/sales", nil)
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

	t.Run("should return sales with custom pagination and filter", func(t *testing.T) {
		expectedResponse := &dto.PaginationResponseDTO{
			TotalRows:  1,
			TotalPages: 1,
			Page:       2,
			Limit:      5,
			Data:       domain.Sales,
		}

		// Setup mock
		uc.Sale.EXPECT().
			GetAllSales(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx *gin.Context, p *dto.PaginationDTO) (*dto.PaginationResponseDTO, error) {
				assert.Equal(t, 2, p.Page)
				assert.Equal(t, 5, p.Limit)
				return expectedResponse, nil
			})

		// Make request with query params
		req, _ := http.NewRequest(http.MethodGet, "/sales?page=2&limit=5&sale_name=test", nil)
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
		// Tidak perlu mock GetAllSales karena tidak dipanggil jika ada error validasi

		req, _ := http.NewRequest(http.MethodGet, "/sales?page=-1&limit=10", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code) // Assert status code
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
		// Tidak perlu mock GetAllSales karena tidak dipanggil jika ada error validasi

		req, _ := http.NewRequest(http.MethodGet, "/sales?limit=101", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code) // Assert status code
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
		uc.Sale.EXPECT().
			GetAllSales(gomock.Any(), gomock.Any()).
			Return(nil, utils.NewInternalError("service_api error"))

		req, _ := http.NewRequest(http.MethodGet, "/sales", nil)
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

func TestSaleHandler_GetSaleByID(t *testing.T) {
	r, h, uc, ids, domain := SaleHandlerSetup(t)
	r.GET("/sales/:id", h.GetSaleByID)

	t.Run("should return sale by id successfully", func(t *testing.T) {
		uc.Sale.EXPECT().GetSaleByID(gomock.Any(), ids.SaleID).Return(domain.Sale, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/sales/"+ids.SaleID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response responseSaleHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, domain.Sale.ID, response.Data.ID)
		assert.Equal(t, domain.Sale.CommodityID, response.Data.CommodityID)
		assert.Equal(t, domain.Sale.CityID, response.Data.CityID)
		assert.Equal(t, domain.Sale.Quantity, response.Data.Quantity)
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.Sale.EXPECT().GetSaleByID(gomock.Any(), ids.SaleID).Return(nil, utils.NewInternalError("service_api error")).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/sales/"+ids.SaleID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseSaleHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/sales/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseSaleHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestSaleHandler_GetSalesByCommodityID(t *testing.T) {
	r, h, uc, ids, domain := SaleHandlerSetup(t)
	r.GET("/sales/commodity/:id", h.GetSalesByCommodityID)

	t.Run("should return all sales successfully with default pagination", func(t *testing.T) {
		// Prepare expected response
		expectedResponse := &dto.PaginationResponseDTO{
			TotalRows:  1,
			TotalPages: 1,
			Page:       1,
			Limit:      10,
			Data:       domain.Sales,
		}

		// Setup mock
		uc.Sale.EXPECT().
			GetSalesByCommodityID(gomock.Any(), gomock.Any(), ids.CommodityID).
			DoAndReturn(func(ctx *gin.Context, p *dto.PaginationDTO, id uuid.UUID) (*dto.PaginationResponseDTO, error) {
				// Verify default pagination
				assert.Equal(t, 1, p.Page)
				assert.Equal(t, 10, p.Limit)
				assert.Equal(t, "created_at desc", p.Sort)
				return expectedResponse, nil
			})

		// Make request
		req, _ := http.NewRequest(http.MethodGet, "/sales/commodity/"+ids.CommodityID.String(), nil)
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

	t.Run("should return sales with custom pagination and filter", func(t *testing.T) {
		expectedResponse := &dto.PaginationResponseDTO{
			TotalRows:  1,
			TotalPages: 1,
			Page:       2,
			Limit:      5,
			Data:       domain.Sales,
		}

		// Setup mock
		uc.Sale.EXPECT().
			GetSalesByCommodityID(gomock.Any(), gomock.Any(), ids.CommodityID).
			DoAndReturn(func(ctx *gin.Context, p *dto.PaginationDTO, id uuid.UUID) (*dto.PaginationResponseDTO, error) {
				assert.Equal(t, 2, p.Page)
				assert.Equal(t, 5, p.Limit)
				return expectedResponse, nil
			})

		// Make request with query params
		req, _ := http.NewRequest(http.MethodGet, "/sales/commodity/"+ids.CommodityID.String()+"?page=2&limit=5&sale_name=test", nil)
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
		// Tidak perlu mock GetAllSales karena tidak dipanggil jika ada error validasi

		req, _ := http.NewRequest(http.MethodGet, "/sales/commodity/"+ids.CommodityID.String()+"?page=-1&limit=10", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code) // Assert status code
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
		// Tidak perlu mock GetAllSales karena tidak dipanggil jika ada error validasi

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/sales/commodity/%s?limit=101", ids.CommodityID.String()), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code) // Assert status code
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
		uc.Sale.EXPECT().
			GetSalesByCommodityID(gomock.Any(), gomock.Any(), ids.CommodityID).
			Return(nil, utils.NewInternalError("service_api error"))

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/sales/commodity/%s", ids.CommodityID.String()), nil)
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

func TestSaleHandler_GetSalesByCityID(t *testing.T) {
	r, h, uc, ids, domain := SaleHandlerSetup(t)
	r.GET("/sales/city/:id", h.GetSalesByCityID)

	t.Run("should return all sales successfully with default pagination", func(t *testing.T) {
		// Prepare expected response
		expectedResponse := &dto.PaginationResponseDTO{
			TotalRows:  1,
			TotalPages: 1,
			Page:       1,
			Limit:      10,
			Data:       domain.Sales,
		}

		// Setup mock
		uc.Sale.EXPECT().
			GetSalesByCityID(gomock.Any(), gomock.Any(), ids.CityID).
			DoAndReturn(func(ctx *gin.Context, p *dto.PaginationDTO, id int64) (*dto.PaginationResponseDTO, error) {
				// Verify default pagination
				assert.Equal(t, 1, p.Page)
				assert.Equal(t, 10, p.Limit)
				assert.Equal(t, "created_at desc", p.Sort)
				return expectedResponse, nil
			})

		// Make request
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/sales/city/%d", ids.CityID), nil)
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

	t.Run("should return sales with custom pagination and filter", func(t *testing.T) {
		expectedResponse := &dto.PaginationResponseDTO{
			TotalRows:  1,
			TotalPages: 1,
			Page:       2,
			Limit:      5,
			Data:       domain.Sales,
		}

		// Setup mock
		uc.Sale.EXPECT().
			GetSalesByCityID(gomock.Any(), gomock.Any(), ids.CityID).
			DoAndReturn(func(ctx *gin.Context, p *dto.PaginationDTO, id int64) (*dto.PaginationResponseDTO, error) {
				assert.Equal(t, 2, p.Page)
				assert.Equal(t, 5, p.Limit)
				return expectedResponse, nil
			})

		// Make request with query params
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/sales/city/%d?page=2&limit=5&sale_name=test", ids.CityID), nil)
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
		// Tidak perlu mock GetAllSales karena tidak dipanggil jika ada error validasi

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/sales/city/%d?page=-1&limit=10", ids.CityID), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code) // Assert status code
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
		// Tidak perlu mock GetAllSales karena tidak dipanggil jika ada error validasi

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/sales/city/%d?limit=101", ids.CityID), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code) // Assert status code
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
		uc.Sale.EXPECT().
			GetSalesByCityID(gomock.Any(), gomock.Any(), ids.CityID).
			Return(nil, utils.NewInternalError("service_api error"))

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/sales/city/%d", ids.CityID), nil)
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

func TestSaleHandler_UpdateSale(t *testing.T) {
	r, h, uc, ids, domain := SaleHandlerSetup(t)
	r.PATCH("/sales/:id", h.UpdateSale)

	t.Run("should update sale successfully", func(t *testing.T) {
		uc.Sale.EXPECT().UpdateSale(gomock.Any(), ids.SaleID, gomock.Any()).Return(domain.Sale, nil).Times(1)

		reqBody := `{"name":"updated"}`
		req, _ := http.NewRequest(http.MethodPatch, "/sales/"+ids.SaleID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		var response responseSaleHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.NoError(t, err)
		assert.Equal(t, domain.Sale.ID, response.Data.ID)
		assert.Equal(t, domain.Sale.CommodityID, response.Data.CommodityID)
		assert.Equal(t, domain.Sale.CityID, response.Data.CityID)
		assert.Equal(t, domain.Sale.Quantity, response.Data.Quantity)
	})

	t.Run("should return error when service_api error", func(t *testing.T) {
		uc.Sale.EXPECT().UpdateSale(gomock.Any(), ids.SaleID, gomock.Any()).Return(nil, utils.NewInternalError("service_api error")).Times(1)

		reqBody := `{"name":"updated"}`
		req, _ := http.NewRequest(http.MethodPatch, "/sales/"+ids.SaleID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseSaleHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		uc.Sale.EXPECT().UpdateSale(gomock.Any(), ids.SaleID, gomock.Any()).Times(0)
		req, _ := http.NewRequest(http.MethodPatch, "/sales/"+ids.SaleID.String(), bytes.NewReader([]byte(`invalid-json`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseSaleHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, "/sales/abc", bytes.NewReader([]byte(`invalid-json`)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseSaleHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when not found", func(t *testing.T) {

		uc.Sale.EXPECT().UpdateSale(gomock.Any(), ids.SaleID, gomock.Any()).Return(nil, utils.NewNotFoundError("sale not found")).Times(1)

		reqBody := `{"commodity_id":"` + ids.CommodityID.String() + `","city_id":1,"quantity":1}`
		req, _ := http.NewRequest(http.MethodPatch, "/sales/"+ids.SaleID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseSaleHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "NOT_FOUND")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestSaleHandler_DeleteSale(t *testing.T) {
	r, h, uc, ids, _ := SaleHandlerSetup(t)
	r.DELETE("/sales/:id", h.DeleteSale)

	t.Run("should delete sale successfully", func(t *testing.T) {
		uc.Sale.EXPECT().DeleteSale(gomock.Any(), ids.SaleID).Times(1)
		req, _ := http.NewRequest(http.MethodDelete, "/sales/"+ids.SaleID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Sale deleted successfully", response.Data.Message)
	})

	t.Run("should return error when service_api error", func(t *testing.T) {
		uc.Sale.EXPECT().DeleteSale(gomock.Any(), ids.SaleID).Return(utils.NewInternalError("service_api error")).Times(1)
		req, _ := http.NewRequest(http.MethodDelete, "/sales/"+ids.SaleID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/sales/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
	})

	t.Run("should return error when not found", func(t *testing.T) {

		uc.Sale.EXPECT().DeleteSale(gomock.Any(), ids.SaleID).Return(utils.NewNotFoundError("sale not found")).Times(1)
		req, _ := http.NewRequest(http.MethodDelete, "/sales/"+ids.SaleID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "NOT_FOUND")
	})
}

func TestSaleHandler_RestoreSale(t *testing.T) {
	r, h, uc, ids, domain := SaleHandlerSetup(t)
	r.PATCH("/sales/:id/restore", h.RestoreSale)

	t.Run("should restore sale successfully", func(t *testing.T) {
		uc.Sale.EXPECT().RestoreSale(gomock.Any(), ids.SaleID).Return(domain.Sale, nil).Times(1)
		req, _ := http.NewRequest(http.MethodPatch, "/sales/"+ids.SaleID.String()+"/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response responseSaleHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, domain.Sale.ID, response.Data.ID)
		assert.Equal(t, domain.Sale.CommodityID, response.Data.CommodityID)
		assert.Equal(t, domain.Sale.CityID, response.Data.CityID)
		assert.Equal(t, domain.Sale.Quantity, response.Data.Quantity)
	})

	t.Run("should return error when service_api error", func(t *testing.T) {
		uc.Sale.EXPECT().RestoreSale(gomock.Any(), ids.SaleID).Return(nil, utils.NewInternalError("service_api error")).Times(1)

		req, _ := http.NewRequest(http.MethodPatch, "/sales/"+ids.SaleID.String()+"/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseSaleHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, "/sales/abc/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseSaleHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when not found", func(t *testing.T) {

		uc.Sale.EXPECT().RestoreSale(gomock.Any(), ids.SaleID).Return(nil, utils.NewNotFoundError("sale not found")).Times(1)
		req, _ := http.NewRequest(http.MethodPatch, "/sales/"+ids.SaleID.String()+"/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseSaleHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "NOT_FOUND")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestSaleHandler_GetAllDeletedSales(t *testing.T) {
}

func TestSaleHandler_GetDeletedSaleByID(t *testing.T) {
	r, h, uc, ids, domain := SaleHandlerSetup(t)
	r.GET("/sales/deleted/:id", h.GetDeletedSaleByID)

	t.Run("should return deleted sale by id successfully", func(t *testing.T) {
		uc.Sale.EXPECT().GetDeletedSaleByID(gomock.Any(), ids.SaleID).Return(domain.Sale, nil).Times(1)
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/sales/deleted/%s", ids.SaleID.String()), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response responseSaleHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, domain.Sale.ID, response.Data.ID)
		assert.Equal(t, domain.Sale.CommodityID, response.Data.CommodityID)
		assert.Equal(t, domain.Sale.CityID, response.Data.CityID)
		assert.Equal(t, domain.Sale.Quantity, response.Data.Quantity)
	})

	t.Run("should return error when service_api error", func(t *testing.T) {
		uc.Sale.EXPECT().GetDeletedSaleByID(gomock.Any(), ids.SaleID).Return(nil, utils.NewInternalError("service_api error")).Times(1)
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/sales/deleted/%s", ids.SaleID.String()), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseSaleHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/sales/deleted/abc"), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseSaleHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when not found", func(t *testing.T) {

		uc.Sale.EXPECT().GetDeletedSaleByID(gomock.Any(), ids.SaleID).Return(nil, utils.NewNotFoundError("sale not found")).Times(1)
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/sales/deleted/%s", ids.SaleID.String()), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseSaleHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "NOT_FOUND")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
