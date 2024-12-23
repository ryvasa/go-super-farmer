package handler_interface

import "github.com/gin-gonic/gin"

type DemandHandler interface {
	CreateDemand(c *gin.Context)
	GetAllDemands(c *gin.Context)
	GetDemandByID(c *gin.Context)
	GetDemandsByCommodityID(c *gin.Context)
	GetDemandsByCityID(c *gin.Context)
	UpdateDemand(c *gin.Context)
	DeleteDemand(c *gin.Context)
	GetDemandHistoryByCommodityIDAndCityID(c *gin.Context)
}
