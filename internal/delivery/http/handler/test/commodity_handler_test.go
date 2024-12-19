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
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	mock_usecase "github.com/ryvasa/go-super-farmer/internal/usecase/mock"
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
)

type responseCommodityHandler struct {
	Status  int              `json:"status"`
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    domain.Commodity `json:"data"`
	Errors  response.Error   `json:"errors"`
}

type CommodityHandlerMocks struct {
	Commodity   *domain.Commodity
	Commodities []*domain.Commodity
}
type CommodityHandlerIDs struct {
	CommodityID uuid.UUID
}

func CommodityHandlerSetUp(t *testing.T) (*gin.Engine, handler_interface.CommodityHandler, *mock_usecase.MockCommodityUsecase, CommodityHandlerIDs, CommodityHandlerMocks) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc := mock_usecase.NewMockCommodityUsecase(ctrl)
	h := handler_implementation.NewCommodityHandler(uc)
	r := gin.Default()

	mocks := CommodityHandlerMocks{
		Commodity: &domain.Commodity{
			ID:          uuid.New(),
			Name:        "commodity",
			Description: "commodity description",
		},
		Commodities: []*domain.Commodity{
			{
				ID:          uuid.New(),
				Name:        "commodity",
				Description: "commodity description",
			},
		},
	}
	ids := CommodityHandlerIDs{
		CommodityID: uuid.New(),
	}

	return r, h, uc, ids, mocks
}

func TestCommodityHandler_CreateCommodity(t *testing.T) {
	r, h, uc, _, mocks := CommodityHandlerSetUp(t)

	r.POST("/commodities", h.CreateCommodity)

	t.Run("should create commodity successfully", func(t *testing.T) {

		uc.EXPECT().CreateCommodity(gomock.Any(), gomock.Any()).Return(mocks.Commodity, nil).Times(1)

		reqBody := `{"name":"commodity","description":"commodity description"}`
		req, _ := http.NewRequest(http.MethodPost, "/commodities", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, response.Data.Name, "commodity")
		assert.Equal(t, response.Data.Description, "commodity description")
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		uc.EXPECT().CreateCommodity(gomock.Any(), gomock.Any()).Times(0)

		req, _ := http.NewRequest(http.MethodPost, "/commodities", bytes.NewReader([]byte("invalid-json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().CreateCommodity(gomock.Any(), gomock.Any()).Return(nil, utils.NewInternalError("Internal error"))

		reqBody := `{"name":"commodity","description":"commodity description"}`
		req, _ := http.NewRequest(http.MethodPost, "/commodities", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestCommodityHandler_GetAllCommodities(t *testing.T) {
	r, h, uc, _, mocks := CommodityHandlerSetUp(t)
	r.GET("/commodities", h.GetAllCommodities)

	t.Run("should return all commodities successfully with default pagination", func(t *testing.T) {
		// Prepare expected pagination
		expectedPagination := &dto.PaginationDTO{
			Page:  1,
			Limit: 10,
			Sort:  "created_at desc",
		}

		// Prepare expected response
		expectedResponse := &dto.PaginationResponseDTO{
			TotalRows:  1,
			TotalPages: 1,
			Page:       1,
			Limit:      10,
			Data:       mocks.Commodities,
		}

		// Setup mock
		uc.EXPECT().
			GetAllCommodities(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx *gin.Context, p *dto.PaginationDTO) (*dto.PaginationResponseDTO, error) {
				// Verify pagination params
				assert.Equal(t, expectedPagination.Page, p.Page)
				assert.Equal(t, expectedPagination.Limit, p.Limit)
				assert.Equal(t, expectedPagination.Sort, p.Sort)
				return expectedResponse, nil
			})

		// Make request
		req, _ := http.NewRequest(http.MethodGet, "/commodities", nil)
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

	t.Run("should return commodities with custom pagination", func(t *testing.T) {
		// Custom pagination params
		expectedPagination := &dto.PaginationDTO{
			Page:  2,
			Limit: 5,
			Filter: dto.PaginationFilterDTO{
				CommodityName: "test",
			},
		}

		expectedResponse := &dto.PaginationResponseDTO{
			TotalRows:  10,
			TotalPages: 2,
			Page:       2,
			Limit:      5,
			Data:       mocks.Commodities,
		}

		// Setup mock
		uc.EXPECT().
			GetAllCommodities(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx *gin.Context, p *dto.PaginationDTO) (*dto.PaginationResponseDTO, error) {
				assert.Equal(t, expectedPagination.Page, p.Page)
				assert.Equal(t, expectedPagination.Limit, p.Limit)
				assert.Equal(t, expectedPagination.Filter.CommodityName, p.Filter.CommodityName)
				return expectedResponse, nil
			})

		// Make request with query params
		req, _ := http.NewRequest(http.MethodGet, "/commodities?page=2&limit=5&commodity_name=test", nil)
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

	t.Run("should return error when usecase returns error", func(t *testing.T) {
		uc.EXPECT().
			GetAllCommodities(gomock.Any(), gomock.Any()).
			Return(nil, utils.NewInternalError("internal error"))

		req, _ := http.NewRequest(http.MethodGet, "/commodities", nil)
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

	t.Run("should return error with invalid pagination params", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/commodities?page=-1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response struct {
			Success bool           `json:"success"`
			Errors  response.Error `json:"errors"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "BAD_REQUEST", response.Errors.Code)
	})

	t.Run("should return error with too large limit", func(t *testing.T) {
		// Make request with invalid limit
		req, _ := http.NewRequest(http.MethodGet, "/commodities?limit=101", nil)
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

	// Test case untuk memastikan default values bekerja
	t.Run("should use default pagination values when not provided", func(t *testing.T) {
		expectedResponse := &dto.PaginationResponseDTO{
			TotalRows:  1,
			TotalPages: 1,
			Page:       1,  // default page
			Limit:      10, // default limit
			Data:       mocks.Commodities,
		}

		// Setup mock
		uc.EXPECT().
			GetAllCommodities(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx *gin.Context, p *dto.PaginationDTO) (*dto.PaginationResponseDTO, error) {
				// Verify default values
				assert.Equal(t, 1, p.Page)
				assert.Equal(t, 10, p.Limit)
				assert.Equal(t, "created_at desc", p.Sort)
				return expectedResponse, nil
			})

		// Make request without pagination params
		req, _ := http.NewRequest(http.MethodGet, "/commodities", nil)
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
		assert.Equal(t, expectedResponse.Page, response.Data.Page)
		assert.Equal(t, expectedResponse.Limit, response.Data.Limit)
	})
}

func TestCommodityHandler_GetCommodityById(t *testing.T) {
	r, h, uc, ids, mocks := CommodityHandlerSetUp(t)

	r.GET("/commodities/:id", h.GetCommodityById)

	t.Run("should return commodity by id successfully", func(t *testing.T) {

		uc.EXPECT().GetCommodityById(gomock.Any(), ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/commodities/"+ids.CommodityID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, mocks.Commodity.Name, response.Data.Name)
		assert.Equal(t, mocks.Commodity.Description, response.Data.Description)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return error when internal error", func(t *testing.T) {

		uc.EXPECT().GetCommodityById(gomock.Any(), ids.CommodityID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/commodities/"+ids.CommodityID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/commodities/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCommodityHandler_UpdateCommodity(t *testing.T) {
	r, h, uc, ids, mocks := CommodityHandlerSetUp(t)

	r.PATCH("/commodities/:id", h.UpdateCommodity)
	t.Run("should update commodity successfully", func(t *testing.T) {

		uc.EXPECT().UpdateCommodity(gomock.Any(), ids.CommodityID, gomock.Any()).Return(mocks.Commodity, nil).Times(1)

		reqBody := `{"name":"updated"}`
		req, _ := http.NewRequest(http.MethodPatch, "/commodities/"+ids.CommodityID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response responseCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, mocks.Commodity.Name, response.Data.Name)
		assert.Equal(t, mocks.Commodity.Description, response.Data.Description)
	})

	t.Run("should return error when internal error", func(t *testing.T) {

		uc.EXPECT().UpdateCommodity(gomock.Any(), ids.CommodityID, gomock.Any()).Return(nil, utils.NewInternalError("internal error")).Times(1)

		reqBody := `{"name":"updated"}`
		req, _ := http.NewRequest(http.MethodPatch, "/commodities/"+ids.CommodityID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when bind error", func(t *testing.T) {

		uc.EXPECT().UpdateCommodity(gomock.Any(), ids.CommodityID, gomock.Any()).Times(0)

		req, _ := http.NewRequest(http.MethodPatch, "/commodities/"+ids.CommodityID.String(), bytes.NewReader([]byte(`invalid-json`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, "/commodities/abc", bytes.NewReader([]byte(`invalid-json`)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when not found", func(t *testing.T) {

		uc.EXPECT().UpdateCommodity(gomock.Any(), ids.CommodityID, gomock.Any()).Return(nil, utils.NewNotFoundError("commodity not found")).Times(1)

		reqBody := `{"name":"updated"}`
		req, _ := http.NewRequest(http.MethodPatch, "/commodities/"+ids.CommodityID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "NOT_FOUND")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestCommodityHandler_DeleteCommodity(t *testing.T) {
	r, h, uc, ids, _ := CommodityHandlerSetUp(t)

	r.DELETE("/commodities/:id", h.DeleteCommodity)

	t.Run("should delete commodity successfully", func(t *testing.T) {

		uc.EXPECT().DeleteCommodity(gomock.Any(), ids.CommodityID).Times(1)

		req, _ := http.NewRequest(http.MethodDelete, "/commodities/"+ids.CommodityID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Commodity deleted successfully", response.Data.Message)
	})

	t.Run("should return error when internal error", func(t *testing.T) {

		uc.EXPECT().DeleteCommodity(gomock.Any(), ids.CommodityID).Return(utils.NewInternalError("internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodDelete, "/commodities/"+ids.CommodityID.String(), nil)
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
		req, _ := http.NewRequest(http.MethodDelete, "/commodities/abc", nil)
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

		uc.EXPECT().DeleteCommodity(gomock.Any(), ids.CommodityID).Return(utils.NewNotFoundError("commodity not found")).Times(1)

		req, _ := http.NewRequest(http.MethodDelete, "/commodities/"+ids.CommodityID.String(), nil)
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

func TestCommodityHandler_RestoreCommodity(t *testing.T) {
	r, h, uc, ids, mocks := CommodityHandlerSetUp(t)

	r.PATCH("/commodities/:id/restore", h.RestoreCommodity)

	t.Run("should restore commodity successfully", func(t *testing.T) {

		uc.EXPECT().RestoreCommodity(gomock.Any(), ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		req, _ := http.NewRequest(http.MethodPatch, "/commodities/"+ids.CommodityID.String()+"/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response responseCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, mocks.Commodity.Name, response.Data.Name)
		assert.Equal(t, mocks.Commodity.Description, response.Data.Description)
	})

	t.Run("should return error when internal error", func(t *testing.T) {

		uc.EXPECT().RestoreCommodity(gomock.Any(), ids.CommodityID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodPatch, "/commodities/"+ids.CommodityID.String()+"/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, "/commodities/abc/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when not found", func(t *testing.T) {

		uc.EXPECT().RestoreCommodity(gomock.Any(), ids.CommodityID).Return(nil, utils.NewNotFoundError("commodity not found")).Times(1)

		req, _ := http.NewRequest(http.MethodPatch, "/commodities/"+ids.CommodityID.String()+"/restore", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseCommodityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "NOT_FOUND")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
