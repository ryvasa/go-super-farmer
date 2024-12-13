package handler

import "github.com/gin-gonic/gin"

type DemandHandler interface {

	CreateDemand(c *gin.Context)
	GetAllDemands(c *gin.Context)
	GetDemandByID(c *gin.Context)
	GetDemandsByCommodityID(c *gin.Context)
	GetDemandsByRegionID(c *gin.Context)
	UpdateDemand(c *gin.Context)
	DeleteDemand(c *gin.Context)

}
