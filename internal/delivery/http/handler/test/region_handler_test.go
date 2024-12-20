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
	mock_usecase "github.com/ryvasa/go-super-farmer/internal/usecase/mock"
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
)

type responseRegionHandler struct {
	Status  int            `json:"status"`
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    domain.Region  `json:"data"`
	Errors  response.Error `json:"errors"`
}

type responseRegionsHandler struct {
	Status  int             `json:"status"`
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    []domain.Region `json:"data"`
	Errors  response.Error  `json:"errors"`
}

type RegionHandlerMocks struct {
	Region  *domain.Region
	Regions []*domain.Region
}

type RegionHandlerIDs struct {
	RegionID   uuid.UUID
	CityID     int64
	ProvinceID int64
}

func RegionHandlerSetUp(t *testing.T) (*gin.Engine, handler_interface.RegionHandler, *mock_usecase.MockRegionUsecase, RegionHandlerIDs, RegionHandlerMocks) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc := mock_usecase.NewMockRegionUsecase(ctrl)
	h := handler_implementation.NewRegionHandler(uc)
	r := gin.Default()

	regionID := uuid.New()
	CityID := int64(1)
	ProvinceID := int64(1)

	ids := RegionHandlerIDs{
		RegionID:   regionID,
		CityID:     CityID,
		ProvinceID: ProvinceID,
	}

	mocks := RegionHandlerMocks{
		Region: &domain.Region{
			ID:         regionID,
			CityID:     CityID,
			ProvinceID: ProvinceID,
		},
		Regions: []*domain.Region{
			{
				ID:         regionID,
				CityID:     CityID,
				ProvinceID: ProvinceID,
			},
		},
	}

	return r, h, uc, ids, mocks
}

func TestRegionHandler_CreateRegion(t *testing.T) {
	r, h, uc, _, mocks := RegionHandlerSetUp(t)

	r.POST("/regions", h.CreateRegion)

	t.Run("should create region successfully", func(t *testing.T) {
		uc.EXPECT().CreateRegion(gomock.Any(), gomock.Any()).Return(mocks.Region, nil).Times(1)

		reqBody := `{"city_id":1,"province_id":1}`
		req, _ := http.NewRequest(http.MethodPost, "/regions", bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseRegionHandler
		assert.NoError(t, json.NewDecoder(w.Body).Decode(&response))

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, response.Data.ID, mocks.Region.ID)
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		uc.EXPECT().CreateRegion(gomock.Any(), gomock.Any()).Times(0)

		req, _ := http.NewRequest(http.MethodPost, "/regions", bytes.NewReader([]byte("invalid-json")))
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
		uc.EXPECT().CreateRegion(gomock.Any(), gomock.Any()).Return(nil, utils.NewInternalError("Internal error"))

		reqBody := `{"city_id":1,"province_id":1}`
		req, _ := http.NewRequest(http.MethodPost, "/regions", bytes.NewReader([]byte(reqBody)))
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

func TestRegionHandler_GetAllRegions(t *testing.T) {
	r, h, uc, _, mocks := RegionHandlerSetUp(t)

	r.GET("/regions", h.GetAllRegions)

	t.Run("should get all regions successfully", func(t *testing.T) {
		uc.EXPECT().GetAllRegions(gomock.Any()).Return(mocks.Regions, nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/regions", nil))

		var response responseRegionsHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, len(mocks.Regions), len(response.Data))
		assert.Equal(t, response.Data[0].ID, (mocks.Regions)[0].ID)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().GetAllRegions(gomock.Any()).Return(nil, utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/regions", nil))

		var response responseRegionsHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestRegionHandler_GetRegionByID(t *testing.T) {
	r, h, uc, ids, mocks := RegionHandlerSetUp(t)

	r.GET("/regions/:id", h.GetRegionByID)

	t.Run("should get region by id successfully", func(t *testing.T) {
		uc.EXPECT().GetRegionByID(gomock.Any(), ids.RegionID).Return(mocks.Region, nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/regions/"+ids.RegionID.String(), nil))

		var response responseRegionHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, response.Data.ID, mocks.Region.ID)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().GetRegionByID(gomock.Any(), ids.RegionID).Return(nil, utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/regions/"+ids.RegionID.String(), nil))

		var response responseRegionHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when id is invalid", func(t *testing.T) {
		uc.EXPECT().GetRegionByID(gomock.Any(), uuid.Nil).Return(nil, utils.NewBadRequestError("ID is invalid"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/regions/aa", nil))

		var response responseRegionHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
