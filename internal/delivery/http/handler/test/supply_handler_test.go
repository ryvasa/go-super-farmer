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

type responseSupplyHandler struct {
	Status  int            `json:"status"`
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    domain.Supply  `json:"data"`
	Errors  response.Error `json:"errors"`
}

type responseSuppliesHandler struct {
	Status  int             `json:"status"`
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    []domain.Supply `json:"data"`
	Errors  response.Error  `json:"errors"`
}

type responseSupplyHistoryHandler struct {
	Status  int             `json:"status"`
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    []domain.Supply `json:"data"`
	Errors  response.Error  `json:"errors"`
}

type SupplyHandlerDomainMocks struct {
	Supply        *domain.Supply
	Supplies      *[]domain.Supply
	UpdatedSupply *domain.Supply
	SupplyHistory *[]domain.SupplyHistory
}

type SupplyHandlerIDs struct {
	SupplyID    uuid.UUID
	CommodityID uuid.UUID
	RegionID    uuid.UUID
}

func SupplyHandlerSetUp(t *testing.T) (*gin.Engine, handler_interface.SupplyHandler, *mock.MockSupplyUsecase, SupplyHandlerIDs, SupplyHandlerDomainMocks) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc := mock.NewMockSupplyUsecase(ctrl)
	h := handler_implementation.NewSupplyHandler(uc)
	r := gin.Default()

	supplyID := uuid.New()
	CommodityID := uuid.New()
	RegionID := uuid.New()

	ids := SupplyHandlerIDs{
		SupplyID:    supplyID,
		CommodityID: CommodityID,
		RegionID:    RegionID,
	}

	mocks := SupplyHandlerDomainMocks{
		Supply: &domain.Supply{
			ID:          supplyID,
			CommodityID: CommodityID,
			RegionID:    RegionID,
		},
		Supplies: &[]domain.Supply{
			{
				ID:          supplyID,
				CommodityID: CommodityID,
				RegionID:    RegionID,
			},
		},
		UpdatedSupply: &domain.Supply{
			ID:          supplyID,
			CommodityID: CommodityID,
			RegionID:    RegionID,
		},
		SupplyHistory: &[]domain.SupplyHistory{
			{
				ID:          supplyID,
				CommodityID: CommodityID,
				RegionID:    RegionID,
			},
		},
	}
	return r, h, uc, ids, mocks
}

func TestSupplyHandler_CreateSupply(t *testing.T) {
	r, h, uc, ids, mocks := SupplyHandlerSetUp(t)

	r.POST("/supplies", h.CreateSupply)

	t.Run("should create supply successfully", func(t *testing.T) {
		uc.EXPECT().CreateSupply(gomock.Any(), gomock.Any()).Return(mocks.Supply, nil).Times(1)

		reqBody := `{"commodity_id":"` + ids.CommodityID.String() + `","region_id":"` + ids.RegionID.String() + `"}`
		req, _ := http.NewRequest(http.MethodPost, "/supplies", bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseSupplyHandler
		assert.NoError(t, json.NewDecoder(w.Body).Decode(&response))

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, response.Data.ID, mocks.Supply.ID)
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		uc.EXPECT().CreateSupply(gomock.Any(), gomock.Any()).Times(0)

		req, _ := http.NewRequest(http.MethodPost, "/supplies", bytes.NewReader([]byte("invalid-json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseSupplyHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().CreateSupply(gomock.Any(), gomock.Any()).Return(nil, utils.NewInternalError("Internal error"))

		reqBody := `{"commodity_id":"` + ids.CommodityID.String() + `","region_id":"` + ids.RegionID.String() + `"}`
		req, _ := http.NewRequest(http.MethodPost, "/supplies", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseSupplyHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestSupplyHandler_GetAllSupply(t *testing.T) {
	r, h, uc, _, mocks := SupplyHandlerSetUp(t)

	r.GET("/supplies", h.GetAllSupply)

	t.Run("should get all supplies successfully", func(t *testing.T) {
		uc.EXPECT().GetAllSupply(gomock.Any()).Return(mocks.Supplies, nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/supplies", nil))

		var response responseSuppliesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, len(*mocks.Supplies), len(response.Data))
		assert.Equal(t, response.Data[0].ID, (*mocks.Supplies)[0].ID)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().GetAllSupply(gomock.Any()).Return(nil, utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/supplies", nil))

		var response responseSupplyHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestSupplyHandler_GetSupplyByID(t *testing.T) {
	r, h, uc, ids, mocks := SupplyHandlerSetUp(t)

	r.GET("/supplies/:id", h.GetSupplyByID)

	t.Run("should get supply by id successfully", func(t *testing.T) {
		uc.EXPECT().GetSupplyByID(gomock.Any(), ids.SupplyID).Return(mocks.Supply, nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/supplies/"+ids.SupplyID.String(), nil))

		var response responseSupplyHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, response.Data.ID, mocks.Supply.ID)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().GetSupplyByID(gomock.Any(), ids.SupplyID).Return(nil, utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/supplies/"+ids.SupplyID.String(), nil))

		var response responseSupplyHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when id is invalid", func(t *testing.T) {
		uc.EXPECT().GetSupplyByID(gomock.Any(), uuid.Nil).Return(nil, utils.NewBadRequestError("ID is invalid"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/supplies/aa", nil))

		var response responseSupplyHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestSupplyHandler_GetSupplyByCommodityID(t *testing.T) {
	r, h, uc, ids, mocks := SupplyHandlerSetUp(t)

	r.GET("/supplies/commodity/:commodity_id", h.GetSupplyByCommodityID)

	t.Run("should get supplies by commodity id successfully", func(t *testing.T) {
		uc.EXPECT().GetSupplyByCommodityID(gomock.Any(), ids.CommodityID).Return(mocks.Supplies, nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/supplies/commodity/"+ids.CommodityID.String(), nil))

		var response responseSuppliesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, len(*mocks.Supplies), len(response.Data))
		assert.Equal(t, response.Data[0].ID, (*mocks.Supplies)[0].ID)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().GetSupplyByCommodityID(gomock.Any(), ids.CommodityID).Return(nil, utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/supplies/commodity/"+ids.CommodityID.String(), nil))

		var response responseSupplyHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when id is invalid", func(t *testing.T) {
		uc.EXPECT().GetSupplyByCommodityID(gomock.Any(), uuid.Nil).Return(nil, utils.NewBadRequestError("ID is invalid"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/supplies/commodity/aa", nil))

		var response responseSupplyHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestSupplyHandler_GetSupplyByRegionID(t *testing.T) {
	r, h, uc, ids, mocks := SupplyHandlerSetUp(t)

	r.GET("/supplies/region/:id", h.GetSupplyByRegionID)

	t.Run("should get supplies by region id successfully", func(t *testing.T) {
		uc.EXPECT().GetSupplyByRegionID(gomock.Any(), ids.RegionID).Return(mocks.Supplies, nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/supplies/region/"+ids.RegionID.String(), nil))

		var response responseSuppliesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, len(*mocks.Supplies), len(response.Data))
		assert.Equal(t, response.Data[0].ID, (*mocks.Supplies)[0].ID)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().GetSupplyByRegionID(gomock.Any(), ids.RegionID).Return(nil, utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/supplies/region/"+ids.RegionID.String(), nil))

		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when id is invalid", func(t *testing.T) {
		uc.EXPECT().GetSupplyByRegionID(gomock.Any(), uuid.Nil).Return(nil, utils.NewBadRequestError("ID is invalid"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/supplies/region/aa", nil))

		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestSupplyHandler_UpdateSupply(t *testing.T) {
	r, h, uc, ids, mocks := SupplyHandlerSetUp(t)

	r.PATCH("/supplies/:id", h.UpdateSupply)

	t.Run("should update supply successfully", func(t *testing.T) {
		uc.EXPECT().UpdateSupply(gomock.Any(), ids.SupplyID, gomock.Any()).Return(mocks.UpdatedSupply, nil).Times(1)

		reqBody := `{"commodity_id":"` + ids.CommodityID.String() + `","region_id":"` + ids.RegionID.String() + `"}`
		req, _ := http.NewRequest(http.MethodPatch, "/supplies/"+ids.SupplyID.String(), bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseSupplyHandler
		assert.NoError(t, json.NewDecoder(w.Body).Decode(&response))

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, response.Data.ID, mocks.UpdatedSupply.ID)
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		uc.EXPECT().UpdateSupply(gomock.Any(), ids.SupplyID, gomock.Any()).Times(0)

		req, _ := http.NewRequest(http.MethodPatch, "/supplies/"+ids.SupplyID.String(), bytes.NewReader([]byte("invalid-json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseSupplyHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().UpdateSupply(gomock.Any(), ids.SupplyID, gomock.Any()).Return(nil, utils.NewInternalError("Internal error"))

		reqBody := `{"commodity_id":"` + ids.CommodityID.String() + `","region_id":"` + ids.RegionID.String() + `"}`
		req, _ := http.NewRequest(http.MethodPatch, "/supplies/"+ids.SupplyID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseSupplyHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when id is invalid", func(t *testing.T) {

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodPatch, "/supplies/aa", nil))

		var response responseSupplyHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestSupplyHandler_DeleteSupply(t *testing.T) {
	r, h, uc, ids, _ := SupplyHandlerSetUp(t)

	r.DELETE("/supplies/:id", h.DeleteSupply)

	t.Run("should delete supply successfully", func(t *testing.T) {
		uc.EXPECT().DeleteSupply(gomock.Any(), ids.SupplyID).Return(nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/supplies/"+ids.SupplyID.String(), nil))

		var response response.ResponseMessage
		assert.NoError(t, json.NewDecoder(w.Body).Decode(&response))

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, response.Data.Message, "Supply deleted successfully")
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().DeleteSupply(gomock.Any(), ids.SupplyID).Return(utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/supplies/"+ids.SupplyID.String(), nil))

		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when id is invalid", func(t *testing.T) {
		uc.EXPECT().DeleteSupply(gomock.Any(), uuid.Nil).Return(utils.NewBadRequestError("ID is invalid"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/supplies/aa", nil))

		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestSupplyHandler_GetSupplyHistoryByCommodityIDAndRegionID(t *testing.T) {
	r, h, uc, ids, mocks := SupplyHandlerSetUp(t)

	r.GET("/supplies/commodity/:commodity_id/region/:region_id", h.GetSupplyHistoryByCommodityIDAndRegionID)

	t.Run("should get supply history by commodity id and region id successfully", func(t *testing.T) {
		uc.EXPECT().GetSupplyHistoryByCommodityIDAndRegionID(gomock.Any(), ids.CommodityID, ids.RegionID).Return(mocks.SupplyHistory, nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/supplies/commodity/"+ids.CommodityID.String()+"/region/"+ids.RegionID.String(), nil))

		var response responseSupplyHistoryHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, len(*mocks.SupplyHistory), len(response.Data))
		assert.Equal(t, response.Data[0].ID, (*mocks.SupplyHistory)[0].ID)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().GetSupplyHistoryByCommodityIDAndRegionID(gomock.Any(), ids.CommodityID, ids.RegionID).Return(nil, utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/supplies/commodity/"+ids.CommodityID.String()+"/region/"+ids.RegionID.String(), nil))

		var response responseSupplyHistoryHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when id is invalid", func(t *testing.T) {
		uc.EXPECT().GetSupplyHistoryByCommodityIDAndRegionID(gomock.Any(), uuid.Nil, uuid.Nil).Return(nil, utils.NewBadRequestError("ID is invalid"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/supplies/commodity/aa/region/bb", nil))

		var response responseSupplyHistoryHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
