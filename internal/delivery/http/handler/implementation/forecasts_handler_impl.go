package handler_implementation

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	handler_interface "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/interface"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	"github.com/ryvasa/go-super-farmer/utils"
)

type ForecastsHandlerImpl struct {
	usecase usecase_interface.ForecastsUsecase
}

func NewForecastsHandler(usecase usecase_interface.ForecastsUsecase) handler_interface.ForecastsHandler {
	return &ForecastsHandlerImpl{usecase}
}

func (h *ForecastsHandlerImpl) GetForecastsByCommodityIDAndCityID(c *gin.Context) {
	cityIDStr := c.Param("city_id")
	cityID, err := strconv.ParseInt(cityIDStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	commodityIDStr := c.Param("land_commodity_id")
	commodityID, err := uuid.Parse(commodityIDStr)
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	logrus.Log.Info("forecasts handle request")

	forecasts, err := h.usecase.GetForecastsByCommodityIDAndCityID(c, commodityID, cityID)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, forecasts)
}
