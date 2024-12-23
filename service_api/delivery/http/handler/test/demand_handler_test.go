package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	handler_implementation "github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler/implementation"
	handler_interface "github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler/interface"
	"github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler/test/response"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	mock_usecase "github.com/ryvasa/go-super-farmer/service_api/usecase/mock"
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
)

type responseDemandHandler struct {
	Status  int            `json:"status"`
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    domain.Demand  `json:"data"`
	Errors  response.Error `json:"errors"`
}

type responseDemandsHandler struct {
	Status  int             `json:"status"`
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    []domain.Demand `json:"data"`
	Errors  response.Error  `json:"errors"`
}

type responseDemandHistoryHandler struct {
	Status  int             `json:"status"`
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    []domain.Demand `json:"data"`
	Errors  response.Error  `json:"errors"`
}

type DemandHandlerDomainMocks struct {
	Demand        *domain.Demand
	Demands       []*domain.Demand
	UpdatedDemand *domain.Demand
	DemandHistory []*domain.DemandHistory
}

type DemandHandlerIDs struct {
	DemandID    uuid.UUID
	CommodityID uuid.UUID
	CityID      int64
	CityIDStr   string
}

func DemandHandlerSetUp(t *testing.T) (*gin.Engine, handler_interface.DemandHandler, *mock_usecase.MockDemandUsecase, DemandHandlerIDs, DemandHandlerDomainMocks) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc := mock_usecase.NewMockDemandUsecase(ctrl)
	h := handler_implementation.NewDemandHandler(uc)
	r := gin.Default()

	demandID := uuid.New()
	commodityID := uuid.New()
	cityID := int64(1)

	cityIDStr := strconv.FormatInt(cityID, 10)

	ids := DemandHandlerIDs{
		DemandID:    demandID,
		CommodityID: commodityID,
		CityID:      cityID,
		CityIDStr:   cityIDStr,
	}

	mocks := DemandHandlerDomainMocks{
		Demand: &domain.Demand{
			ID:          demandID,
			CommodityID: commodityID,
			CityID:      cityID,
		},
		Demands: []*domain.Demand{
			{
				ID:          demandID,
				CommodityID: commodityID,
				CityID:      cityID,
			},
		},
		UpdatedDemand: &domain.Demand{
			ID:          demandID,
			CommodityID: commodityID,
			CityID:      cityID,
		},
		DemandHistory: []*domain.DemandHistory{
			{
				ID:          demandID,
				CommodityID: commodityID,
				CityID:      cityID,
			},
		},
	}
	return r, h, uc, ids, mocks
}

func TestDemandHandler_CreateDemand(t *testing.T) {
	r, h, uc, ids, mocks := DemandHandlerSetUp(t)

	r.POST("/demands", h.CreateDemand)

	t.Run("should create demand successfully", func(t *testing.T) {
		uc.EXPECT().CreateDemand(gomock.Any(), gomock.Any()).Return(mocks.Demand, nil).Times(1)

		reqBody := `{"commodity_id":"` + ids.CommodityID.String() + `","city_id":` + ids.CityIDStr + `}`
		req, _ := http.NewRequest(http.MethodPost, "/demands", bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseDemandHandler
		assert.NoError(t, json.NewDecoder(w.Body).Decode(&response))

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, response.Data.ID, mocks.Demand.ID)
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		uc.EXPECT().CreateDemand(gomock.Any(), gomock.Any()).Times(0)

		req, _ := http.NewRequest(http.MethodPost, "/demands", bytes.NewReader([]byte("invalid-json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseDemandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
	})

	t.Run("should return error when service_api error", func(t *testing.T) {
		uc.EXPECT().CreateDemand(gomock.Any(), gomock.Any()).Return(nil, utils.NewInternalError("Internal error"))

		reqBody := `{"commodity_id":"` + ids.CommodityID.String() + `","city_id":` + ids.CityIDStr + `}`
		req, _ := http.NewRequest(http.MethodPost, "/demands", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseDemandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestDemandHandler_GetAllDemands(t *testing.T) {
	r, h, uc, _, mocks := DemandHandlerSetUp(t)

	r.GET("/demands", h.GetAllDemands)

	t.Run("should get all demands successfully", func(t *testing.T) {
		uc.EXPECT().GetAllDemands(gomock.Any()).Return(mocks.Demands, nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/demands", nil))

		var response responseDemandsHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, len(mocks.Demands), len(response.Data))
		assert.Equal(t, response.Data[0].ID, (mocks.Demands)[0].ID)
	})

	t.Run("should return error when service_api error", func(t *testing.T) {
		uc.EXPECT().GetAllDemands(gomock.Any()).Return(nil, utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/demands", nil))

		var response responseDemandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestDemandHandler_GetDemandByID(t *testing.T) {
	r, h, uc, ids, mocks := DemandHandlerSetUp(t)

	r.GET("/demands/:id", h.GetDemandByID)

	t.Run("should get demand by id successfully", func(t *testing.T) {
		uc.EXPECT().GetDemandByID(gomock.Any(), ids.DemandID).Return(mocks.Demand, nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/demands/"+ids.DemandID.String(), nil))

		var response responseDemandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, response.Data.ID, mocks.Demand.ID)
	})

	t.Run("should return error when service_api error", func(t *testing.T) {
		uc.EXPECT().GetDemandByID(gomock.Any(), ids.DemandID).Return(nil, utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/demands/"+ids.DemandID.String(), nil))

		var response responseDemandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when id is invalid", func(t *testing.T) {
		uc.EXPECT().GetDemandByID(gomock.Any(), uuid.Nil).Return(nil, utils.NewBadRequestError("ID is invalid"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/demands/aa", nil))

		var response responseDemandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDemandHandler_GetDemandsByCommodityID(t *testing.T) {
	r, h, uc, ids, mocks := DemandHandlerSetUp(t)

	r.GET("/demands/commodity/:commodity_id", h.GetDemandsByCommodityID)

	t.Run("should get demands by commodity id successfully", func(t *testing.T) {
		uc.EXPECT().GetDemandsByCommodityID(gomock.Any(), ids.CommodityID).Return(mocks.Demands, nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/demands/commodity/"+ids.CommodityID.String(), nil))

		var response responseDemandsHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, len(mocks.Demands), len(response.Data))
		assert.Equal(t, response.Data[0].ID, (mocks.Demands)[0].ID)
	})

	t.Run("should return error when service_api error", func(t *testing.T) {
		uc.EXPECT().GetDemandsByCommodityID(gomock.Any(), ids.CommodityID).Return(nil, utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/demands/commodity/"+ids.CommodityID.String(), nil))

		var response responseDemandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when id is invalid", func(t *testing.T) {
		uc.EXPECT().GetDemandsByCommodityID(gomock.Any(), uuid.Nil).Return(nil, utils.NewBadRequestError("ID is invalid"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/demands/commodity/aa", nil))

		var response responseDemandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDemandHandler_GetDemandsByCityID(t *testing.T) {
	r, h, uc, ids, mocks := DemandHandlerSetUp(t)

	r.GET("/demands/city/:id", h.GetDemandsByCityID)

	t.Run("should get demands by city id successfully", func(t *testing.T) {
		uc.EXPECT().GetDemandsByCityID(gomock.Any(), ids.CityID).Return(mocks.Demands, nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, fmt.Sprintf("/demands/city/%d", ids.CityID), nil))

		var response responseDemandsHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, len(mocks.Demands), len(response.Data))
		assert.Equal(t, response.Data[0].ID, (mocks.Demands)[0].ID)
	})

	t.Run("should return error when service_api error", func(t *testing.T) {
		uc.EXPECT().GetDemandsByCityID(gomock.Any(), ids.CityID).Return(nil, utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/demands/city/"+ids.CityIDStr, nil))

		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when id is invalid", func(t *testing.T) {
		uc.EXPECT().GetDemandsByCityID(gomock.Any(), uuid.Nil).Return(nil, utils.NewBadRequestError("ID is invalid"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/demands/city/aa", nil))

		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDemandHandler_UpdateDemand(t *testing.T) {
	r, h, uc, ids, mocks := DemandHandlerSetUp(t)

	r.PATCH("/demands/:id", h.UpdateDemand)

	t.Run("should update demand successfully", func(t *testing.T) {
		uc.EXPECT().UpdateDemand(gomock.Any(), ids.DemandID, gomock.Any()).Return(mocks.UpdatedDemand, nil).Times(1)

		reqBody := `{"commodity_id":"` + ids.CommodityID.String() + `","city_id":` + ids.CityIDStr + `}`
		req, _ := http.NewRequest(http.MethodPatch, "/demands/"+ids.DemandID.String(), bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseDemandHandler
		assert.NoError(t, json.NewDecoder(w.Body).Decode(&response))

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, response.Data.ID, mocks.UpdatedDemand.ID)
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		uc.EXPECT().UpdateDemand(gomock.Any(), ids.DemandID, gomock.Any()).Times(0)

		req, _ := http.NewRequest(http.MethodPatch, "/demands/"+ids.DemandID.String(), bytes.NewReader([]byte("invalid-json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseDemandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
	})

	t.Run("should return error when service_api error", func(t *testing.T) {
		uc.EXPECT().UpdateDemand(gomock.Any(), ids.DemandID, gomock.Any()).Return(nil, utils.NewInternalError("Internal error"))

		reqBody := `{"commodity_id":"` + ids.CommodityID.String() + `","city_id":` + ids.CityIDStr + `}`
		req, _ := http.NewRequest(http.MethodPatch, "/demands/"+ids.DemandID.String(), bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseDemandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when id is invalid", func(t *testing.T) {

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodPatch, "/demands/aa", nil))

		var response responseDemandHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDemandHandler_DeleteDemand(t *testing.T) {
	r, h, uc, ids, _ := DemandHandlerSetUp(t)

	r.DELETE("/demands/:id", h.DeleteDemand)

	t.Run("should delete demand successfully", func(t *testing.T) {
		uc.EXPECT().DeleteDemand(gomock.Any(), ids.DemandID).Return(nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/demands/"+ids.DemandID.String(), nil))

		var response response.ResponseMessage
		assert.NoError(t, json.NewDecoder(w.Body).Decode(&response))

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, response.Data.Message, "Demand deleted successfully")
	})

	t.Run("should return error when service_api error", func(t *testing.T) {
		uc.EXPECT().DeleteDemand(gomock.Any(), ids.DemandID).Return(utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/demands/"+ids.DemandID.String(), nil))

		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when id is invalid", func(t *testing.T) {
		uc.EXPECT().DeleteDemand(gomock.Any(), uuid.Nil).Return(utils.NewBadRequestError("ID is invalid"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/demands/aa", nil))

		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDemandHandler_GetDemandHistoryByCommodityIDAndCityID(t *testing.T) {
	r, h, uc, ids, mocks := DemandHandlerSetUp(t)

	r.GET("/demands/commodity/:commodity_id/city/:city_id", h.GetDemandHistoryByCommodityIDAndCityID)

	t.Run("should get demand history by commodity id and city id successfully", func(t *testing.T) {
		uc.EXPECT().GetDemandHistoryByCommodityIDAndCityID(gomock.Any(), ids.CommodityID, ids.CityID).Return(mocks.DemandHistory, nil).Times(1)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/demands/commodity/"+ids.CommodityID.String()+"/city/"+ids.CityIDStr, nil))

		var response responseDemandHistoryHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, len(mocks.DemandHistory), len(response.Data))
		assert.Equal(t, response.Data[0].ID, (mocks.DemandHistory)[0].ID)
	})

	t.Run("should return error when service_api error", func(t *testing.T) {
		uc.EXPECT().GetDemandHistoryByCommodityIDAndCityID(gomock.Any(), ids.CommodityID, ids.CityID).Return(nil, utils.NewInternalError("Internal error"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/demands/commodity/"+ids.CommodityID.String()+"/city/"+ids.CityIDStr, nil))

		var response responseDemandHistoryHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when id is invalid", func(t *testing.T) {
		uc.EXPECT().GetDemandHistoryByCommodityIDAndCityID(gomock.Any(), uuid.Nil, uuid.Nil).Return(nil, utils.NewBadRequestError("ID is invalid"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/demands/commodity/aa/city/bb", nil))

		var response responseDemandHistoryHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when city id is invalid", func(t *testing.T) {
		uc.EXPECT().GetDemandHistoryByCommodityIDAndCityID(gomock.Any(), ids.CommodityID, uuid.Nil).Return(nil, utils.NewBadRequestError("ID is invalid"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/demands/commodity/"+ids.CommodityID.String()+"/city/aa", nil))

		var response responseDemandHistoryHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
