package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	handler_implementation "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/implementation"
	handler_interface "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/interface"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/test/response"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	mock_usecase "github.com/ryvasa/go-super-farmer/internal/usecase/mock"
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
)

type responseCityHandler struct {
	Status  int            `json:"status"`
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    domain.City    `json:"data"`
	Errors  response.Error `json:"errors"`
}

type CityHandlerMocks struct {
	City   *domain.City
	Cities []*domain.City
}
type CityHandlerIDs struct {
	CityID        int64
	ProvinceID    int64
	CityIDstr     string
	ProvinceIDstr string
}

func CityHandlerSetUp(t *testing.T) (*gin.Engine, handler_interface.CityHandler, *mock_usecase.MockCityUsecase, CityHandlerIDs, CityHandlerMocks) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc := mock_usecase.NewMockCityUsecase(ctrl)
	h := handler_implementation.NewCityHandler(uc)
	r := gin.Default()

	cityID := int64(1)
	provinceID := int64(1)

	cityIDstr := strconv.FormatInt(cityID, 10)
	provinceIDstr := strconv.FormatInt(provinceID, 10)
	ids := CityHandlerIDs{
		CityID:        cityID,
		ProvinceID:    provinceID,
		CityIDstr:     cityIDstr,
		ProvinceIDstr: provinceIDstr,
	}

	mocks := CityHandlerMocks{
		City: &domain.City{
			ID:         cityID,
			Name:       "city",
			ProvinceID: provinceID,
		},
		Cities: []*domain.City{
			{
				ID:         cityID,
				Name:       "city",
				ProvinceID: provinceID,
			},
		},
	}

	return r, h, uc, ids, mocks
}

func TestCityHandler_CreateCity(t *testing.T) {
	r, h, uc, _, mocks := CityHandlerSetUp(t)
	r.POST("/cities", h.CreateCity)

	t.Run("should create city successfully", func(t *testing.T) {
		uc.EXPECT().CreateCity(gomock.Any(), gomock.Any()).Return(mocks.City, nil).Times(1)

		reqBody := `{"name":"city","province_id":1}`

		req, _ := http.NewRequest(http.MethodPost, "/cities", bytes.NewReader([]byte(reqBody)))

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseCityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, response.Data.Name, "city")
		assert.Equal(t, response.Data.ProvinceID, int64(1))
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		uc.EXPECT().CreateCity(gomock.Any(), gomock.Any()).Times(0)

		req, _ := http.NewRequest(http.MethodPost, "/cities", bytes.NewReader([]byte(`invalid-json`)))

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

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().CreateCity(gomock.Any(), gomock.Any()).Return(nil, utils.NewInternalError("Internal error"))

		reqBody := `{"name":"city","province_id":1}`

		req, _ := http.NewRequest(http.MethodPost, "/cities", bytes.NewReader([]byte(reqBody)))

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseCityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestCityHandler_GetCityByID(t *testing.T) {
	r, h, uc, ids, mocks := CityHandlerSetUp(t)

	r.GET("/cities/:id", h.GetCityByID)

	t.Run("should return city by id successfully", func(t *testing.T) {
		uc.EXPECT().GetCityByID(gomock.Any(), ids.CityID).Return(mocks.City, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/cities/"+ids.CityIDstr, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseCityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, response.Data.Name, "city")
		assert.Equal(t, response.Data.ProvinceID, int64(1))
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().GetCityByID(gomock.Any(), ids.CityID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/cities/"+ids.CityIDstr, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseCityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/cities/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseCityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when not found", func(t *testing.T) {
		uc.EXPECT().GetCityByID(gomock.Any(), ids.CityID).Return(nil, utils.NewNotFoundError("city not found")).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/cities/"+ids.CityIDstr, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseCityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "NOT_FOUND")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestCityHandler_GetAllCities(t *testing.T) {
	r, h, uc, _, mocks := CityHandlerSetUp(t)
	r.GET("/cities", h.GetAllCities)
	t.Run("should return all cities successfully", func(t *testing.T) {
		uc.EXPECT().GetAllCities(gomock.Any()).Return(mocks.Cities, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/cities", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var response struct {
			Data []domain.City `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Data, len(mocks.Cities))
	})
	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().GetAllCities(gomock.Any()).Return(nil, utils.NewInternalError("Internal error")).Times(1)
		req, _ := http.NewRequest(http.MethodGet, "/cities", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var response responseCityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestCityHandler_UpdateCity(t *testing.T) {
	r, h, uc, ids, mocks := CityHandlerSetUp(t)
	r.PATCH("/cities/:id", h.UpdateCity)
	t.Run("should update city successfully", func(t *testing.T) {
		uc.EXPECT().UpdateCity(gomock.Any(), ids.CityID, gomock.Any()).Return(mocks.City, nil).Times(1)

		reqBody := `{"name":"updated"}`
		req, _ := http.NewRequest(http.MethodPatch, "/cities/"+ids.CityIDstr, bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response responseCityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, mocks.City.Name, response.Data.Name)
		assert.Equal(t, mocks.City.ProvinceID, response.Data.ProvinceID)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().UpdateCity(gomock.Any(), ids.CityID, gomock.Any()).Return(nil, utils.NewInternalError("internal error")).Times(1)

		reqBody := `{"name":"updated"}`
		req, _ := http.NewRequest(http.MethodPatch, "/cities/"+ids.CityIDstr, bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseCityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		uc.EXPECT().UpdateCity(gomock.Any(), ids.CityID, gomock.Any()).Times(0)

		req, _ := http.NewRequest(http.MethodPatch, "/cities/"+ids.CityIDstr, bytes.NewReader([]byte(`invalid-json`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseCityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, "/cities/abc", bytes.NewReader([]byte(`invalid-json`)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseCityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when not found", func(t *testing.T) {
		uc.EXPECT().UpdateCity(gomock.Any(), ids.CityID, gomock.Any()).Return(nil, utils.NewNotFoundError("city not found")).Times(1)

		reqBody := `{"name":"updated"}`
		req, _ := http.NewRequest(http.MethodPatch, "/cities/"+ids.CityIDstr, bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseCityHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "NOT_FOUND")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestCityHandler_DeleteCity(t *testing.T) {
	r, h, uc, ids, _ := CityHandlerSetUp(t)

	r.DELETE("/cities/:id", h.DeleteCity)

	t.Run("should delete city successfully", func(t *testing.T) {
		uc.EXPECT().DeleteCity(gomock.Any(), ids.CityID).Return(nil).Times(1)

		req, _ := http.NewRequest(http.MethodDelete, "/cities/"+ids.CityIDstr, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "City deleted successfully", response.Data.Message)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		uc.EXPECT().DeleteCity(gomock.Any(), ids.CityID).Return(utils.NewInternalError("internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodDelete, "/cities/"+ids.CityIDstr, nil)
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
		req, _ := http.NewRequest(http.MethodDelete, "/cities/abc", nil)
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
		uc.EXPECT().DeleteCity(gomock.Any(), ids.CityID).Return(utils.NewNotFoundError("city not found")).Times(1)

		req, _ := http.NewRequest(http.MethodDelete, "/cities/"+ids.CityIDstr, nil)
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
