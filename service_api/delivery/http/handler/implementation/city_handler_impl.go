package handler_implementation

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	handler_interface "github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler/interface"
	"github.com/ryvasa/go-super-farmer/service_api/model/dto"
	usecase_interface "github.com/ryvasa/go-super-farmer/service_api/usecase/interface"
	"github.com/ryvasa/go-super-farmer/utils"
)

type CityHandlerImpl struct {
	uc usecase_interface.CityUsecase
}

func NewCityHandler(uc usecase_interface.CityUsecase) handler_interface.CityHandler {
	return &CityHandlerImpl{uc}
}

func (h *CityHandlerImpl) CreateCity(c *gin.Context) {
	req := new(dto.CityCreateDTO)
	if err := c.ShouldBindJSON(req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	city, err := h.uc.CreateCity(c, req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	logrus.Log.Info("API service started successfully")

	utils.SuccessResponse(c, http.StatusCreated, city)
}

func (h *CityHandlerImpl) GetAllCities(c *gin.Context) {
	cities, err := h.uc.GetAllCities(c)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, cities)
}

func (h *CityHandlerImpl) GetCityByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	city, err := h.uc.GetCityByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, city)
}

func (h *CityHandlerImpl) UpdateCity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	req := new(dto.CityUpdateDTO)
	if err := c.ShouldBindJSON(req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	city, err := h.uc.UpdateCity(c, id, req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, city)
}

func (h *CityHandlerImpl) DeleteCity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	err = h.uc.DeleteCity(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "City deleted successfully"})
}
