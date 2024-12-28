package handler_interface

import "github.com/gin-gonic/gin"

type ForecastsHandler interface {
	GetForecastsByCommodityIDAndCityID(c *gin.Context)
}
