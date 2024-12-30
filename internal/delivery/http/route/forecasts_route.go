package route

import (
	"github.com/gin-gonic/gin"
	handler_interface "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/interface"
)

type ForecastsRoute struct {
	handler handler_interface.ForecastsHandler
}

func NewForecastsRoute(handler handler_interface.ForecastsHandler) *ForecastsRoute {
	return &ForecastsRoute{handler}
}

func (r *ForecastsRoute) Register(public, protected *gin.RouterGroup) {
	protected.GET("/forecasts/city/:city_id/land_commodity/:land_commodity_id", r.handler.GetForecastsByCommodityIDAndCityID)
}
