package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/usecase/mock"
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
)

type ResponseLandhHandler struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    domain.Land `json:"data"`
	Errors  interface{} `json:"errors"`
}

func TestCreateLand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock LandUsecase
	usecase := mock.NewMockLandUsecase(ctrl)
	h := handler.NewLandHandler(usecase)

	// Set up Gin router
	r := gin.Default()

	// Mocking the GetAuthUserID function to return a valid userID
	originalGetAuthUserID := utils.GetAuthUserID
	defer func() { utils.GetAuthUserID = originalGetAuthUserID }() // Restore original after the test
	utils.GetAuthUserID = func(c *gin.Context) (uuid.UUID, error) {
		return uuid.New(), nil // Return a valid user ID
	}

	// Register handler
	r.POST("/lands", h.CreateLand)

	t.Run("Test CreateLand, successfully", func(t *testing.T) {
		// Prepare data
		userID := uuid.New()
		landID := uuid.New()
		mockLand := &domain.Land{ID: landID, UserID: userID, LandArea: 10, Certificate: "test"}

		// Mocking the usecase's CreateLand method
		usecase.EXPECT().CreateLand(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockLand, nil).Times(1)

		// Prepare request body
		reqBody := `{"land_area":10,"certificate":"test"}`
		req, _ := http.NewRequest(http.MethodPost, "/lands", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		// Prepare response recorder
		w := httptest.NewRecorder()

		// Execute the request
		r.ServeHTTP(w, req)

		// Assert the HTTP status code and response
		assert.Equal(t, http.StatusCreated, w.Code)

		// Optionally, assert that the response body is as expected
		// You can adjust this based on what the handler returns in the success response
		// assert.Contains(t, w.Body.String(), "created")
	})

	t.Run("Test CreateLand, when GetAuthUserID fails", func(t *testing.T) {
		// Mocking the GetAuthUserID function to return an error
		utils.GetAuthUserID = func(c *gin.Context) (uuid.UUID, error) {
			return uuid.UUID{}, utils.NewUnauthorizedError("unauthorized")
		}

		// Prepare data
		reqBody := `{"land_area":10,"certificate":"test"}`
		req, _ := http.NewRequest(http.MethodPost, "/lands", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		// Prepare response recorder
		w := httptest.NewRecorder()

		// Execute the request
		r.ServeHTTP(w, req)

		// Assert the HTTP status code for unauthorized error
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Test CreateLand, when ShouldBindJSON fails", func(t *testing.T) {
		// Prepare invalid request body
		reqBody := `{"land_area":"invalid","certificate":"test"}`
		req, _ := http.NewRequest(http.MethodPost, "/lands", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		// Prepare response recorder
		w := httptest.NewRecorder()

		// Execute the request
		r.ServeHTTP(w, req)

		// Assert the HTTP status code for bad request
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetOneLand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockLandUsecase(ctrl)
	h := handler.NewLandHandler(usecase)
	r := gin.Default()
	r.GET("/lands:id", h.GetLandByID)

	t.Run("Test GetLandByID, successfully", func(t *testing.T) {

	})

	t.Run("Test GetLandByID, database error", func(t *testing.T) {

	})

	t.Run("Test GetLandByID, invalid id", func(t *testing.T) {

	})
}

func TestGetOneLandByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockLandUsecase(ctrl)
	h := handler.NewLandHandler(usecase)
	r := gin.Default()
	r.GET("/lands/user/:id", h.GetLandByUserID)

	t.Run("Test GetLandByUserID, successfully", func(t *testing.T) {

	})

	t.Run("Test GetLandByUserID, database error", func(t *testing.T) {

	})

	t.Run("Test GetLandByUserID, invalid id", func(t *testing.T) {

	})
}

func TestGetAllLands(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockLandUsecase(ctrl)
	h := handler.NewLandHandler(usecase)
	r := gin.Default()
	r.GET("/lands", h.GetAllLands)

	t.Run("Test GetAllLands, successfully", func(t *testing.T) {

	})

	t.Run("Test GetAllLands, database error", func(t *testing.T) {

	})
}

func TestUpdateLand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockLandUsecase(ctrl)
	h := handler.NewLandHandler(usecase)
	r := gin.Default()
	r.PATCH("/lands/:id", h.UpdateLand)

	t.Run("Test UpdateLand, successfully", func(t *testing.T) {

	})

	t.Run("Test UpdateLand, database error", func(t *testing.T) {

	})

	t.Run("Test UpdateLand, invalid id", func(t *testing.T) {

	})

	t.Run("Test UpdateLand, bind error", func(t *testing.T) {

	})
}

func TestDeleteLand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockLandUsecase(ctrl)
	h := handler.NewLandHandler(usecase)
	r := gin.Default()
	r.DELETE("/lands/:id", h.DeleteLand)

	t.Run("Test DeleteLand, successfully", func(t *testing.T) {

	})

	t.Run("Test DeleteLand, database error", func(t *testing.T) {

	})

	t.Run("Test DeleteLand, invalid id", func(t *testing.T) {

	})
}

func TestRestoreLand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockLandUsecase(ctrl)
	h := handler.NewLandHandler(usecase)
	r := gin.Default()
	r.PATCH("/lands/:id/restore", h.RestoreLand)

	t.Run("Test RestoreLand, successfully", func(t *testing.T) {

	})

	t.Run("Test RestoreLand, database error", func(t *testing.T) {

	})

	t.Run("Test RestoreLand, invalid id", func(t *testing.T) {

	})
}
