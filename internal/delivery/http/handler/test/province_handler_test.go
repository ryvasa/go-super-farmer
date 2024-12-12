package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/test/response"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/usecase/mock"
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
)

type responseProvinceHandler struct {
	Status  int             `json:"status"`
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    domain.Province `json:"data"`
	Errors  response.Error  `json:"errors"`
}

type responseProvincesHandler struct {
	Status  int               `json:"status"`
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    []domain.Province `json:"data"`
	Errors  response.Error    `json:"errors"`
}

type ProvinceHandlerMocks struct {
	Province        *domain.Province
	Provinces       *[]domain.Province
	UpdatedProvince *domain.Province
}
type ProvinceHandlerIDs struct {
	ProvinceID    int64
	ProvinceIDstr string
}

func ProvinceHandlerSetUp(t *testing.T) (*gin.Engine, handler.ProvinceHandler, *mock.MockProvinceUsecase, ProvinceHandlerIDs, ProvinceHandlerMocks) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc := mock.NewMockProvinceUsecase(ctrl)
	h := handler.NewProvinceHandler(uc)
	r := gin.Default()

	ids := ProvinceHandlerIDs{
		ProvinceID:    1,
		ProvinceIDstr: "1",
	}

	mocks := ProvinceHandlerMocks{
		Province: &domain.Province{
			ID:   ids.ProvinceID,
			Name: "Province",
		},
		Provinces: &[]domain.Province{
			{
				ID:   ids.ProvinceID,
				Name: "Province",
			},
		},
		UpdatedProvince: &domain.Province{
			ID:   ids.ProvinceID,
			Name: "updated",
		},
	}

	return r, h, uc, ids, mocks
}

func TestProvinceHandler_CreateProvince(t *testing.T) {
	r, h, uc, _, mocks := ProvinceHandlerSetUp(t)
	r.POST("/provinces", h.CreateProvince)

	t.Run("should create province successfully", func(t *testing.T) {
		uc.EXPECT().CreateProvince(gomock.Any(), gomock.Any()).Return(mocks.Province, nil).Times(1)

		reqBody := `{"name":"Province"}`
		req, _ := http.NewRequest(http.MethodPost, "/provinces", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseProvinceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, response.Data.Name, "Province")
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		uc.EXPECT().CreateProvince(gomock.Any(), gomock.Any()).Times(0)

		req, _ := http.NewRequest(http.MethodPost, "/provinces", bytes.NewReader([]byte(`invalid-json`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseProvinceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.EXPECT().CreateProvince(gomock.Any(), gomock.Any()).Return(nil, utils.NewInternalError("internal server error")).Times(1)

		reqBody := `{"name":"Province"}`
		req, _ := http.NewRequest(http.MethodPost, "/provinces", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseProvinceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestProvinceHandler_GetAllProvinces(t *testing.T) {
	r, h, uc, _, mocks := ProvinceHandlerSetUp(t)
	r.GET("/provinces", h.GetAllProvinces)

	t.Run("should return all provinces successfully", func(t *testing.T) {
		uc.EXPECT().GetAllProvinces(gomock.Any()).Return(mocks.Provinces, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/provinces", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseProvincesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Data, len(*mocks.Provinces))
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.EXPECT().GetAllProvinces(gomock.Any()).Return(nil, utils.NewInternalError("internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/provinces", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseProvincesHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestProvinceHandler_GetProvinceByID(t *testing.T) {
	r, h, uc, ids, mocks := ProvinceHandlerSetUp(t)
	r.GET("/provinces/:id", h.GetProvinceByID)

	t.Run("should return province by id successfully", func(t *testing.T) {
		uc.EXPECT().GetProvinceByID(gomock.Any(), ids.ProvinceID).Return(mocks.Province, nil).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/provinces/"+ids.ProvinceIDstr, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseProvinceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, response.Data.Name, "Province")
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.EXPECT().GetProvinceByID(gomock.Any(), ids.ProvinceID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		req, _ := http.NewRequest(http.MethodGet, "/provinces/"+ids.ProvinceIDstr, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseProvinceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/provinces/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseProvinceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestProvinceHandler_UpdateProvince(t *testing.T) {
	r, h, uc, ids, mocks := ProvinceHandlerSetUp(t)
	r.PATCH("/provinces/:id", h.UpdateProvince)

	t.Run("should update province successfully", func(t *testing.T) {
		uc.EXPECT().UpdateProvince(gomock.Any(), ids.ProvinceID, gomock.Any()).Return(mocks.UpdatedProvince, nil).Times(1)

		reqBody := `{"name":"updated"}`
		req, _ := http.NewRequest(http.MethodPatch, "/provinces/"+ids.ProvinceIDstr, bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseProvinceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, response.Data.Name, "updated")
	})

	t.Run("should return error when bind error", func(t *testing.T) {
		uc.EXPECT().UpdateProvince(gomock.Any(), ids.ProvinceID, gomock.Any()).Times(0)

		req, _ := http.NewRequest(http.MethodPatch, "/provinces/"+ids.ProvinceIDstr, bytes.NewReader([]byte(`invalid-json`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseProvinceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.EXPECT().UpdateProvince(gomock.Any(), ids.ProvinceID, gomock.Any()).Return(nil, utils.NewInternalError("internal server error")).Times(1)

		reqBody := `{"name":"updated"}`
		req, _ := http.NewRequest(http.MethodPatch, "/provinces/"+ids.ProvinceIDstr, bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseProvinceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, "/provinces/abc", bytes.NewReader([]byte(`invalid-json`)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseProvinceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "BAD_REQUEST")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when not found", func(t *testing.T) {
		uc.EXPECT().UpdateProvince(gomock.Any(), ids.ProvinceID, gomock.Any()).Return(nil, utils.NewNotFoundError("province not found")).Times(1)

		reqBody := `{"name":"updated"}`
		req, _ := http.NewRequest(http.MethodPatch, "/provinces/"+ids.ProvinceIDstr, bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response responseProvinceHandler
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Errors)
		assert.Equal(t, response.Errors.Code, "NOT_FOUND")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestProvinceHandler_DeleteProvince(t *testing.T) {
	r, h, uc, ids, _ := ProvinceHandlerSetUp(t)
	r.DELETE("/provinces/:id", h.DeleteProvince)

	t.Run("should delete province successfully", func(t *testing.T) {
		uc.EXPECT().DeleteProvince(gomock.Any(), ids.ProvinceID).Times(1)

		req, _ := http.NewRequest(http.MethodDelete, "/provinces/"+ids.ProvinceIDstr, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Province deleted successfully", response.Data.Message)
	})

	t.Run("should return error when usecase error", func(t *testing.T) {
		uc.EXPECT().DeleteProvince(gomock.Any(), ids.ProvinceID).Return(utils.NewInternalError("internal server error")).Times(1)

		req, _ := http.NewRequest(http.MethodDelete, "/provinces/"+ids.ProvinceIDstr, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response response.ResponseMessage
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.Errors.Code, "INTERNAL_ERROR")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/provinces/abc", nil)
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
		uc.EXPECT().DeleteProvince(gomock.Any(), ids.ProvinceID).Return(utils.NewNotFoundError("province not found")).Times(1)

		req, _ := http.NewRequest(http.MethodDelete, "/provinces/"+ids.ProvinceIDstr, nil)
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
