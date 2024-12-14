package handler

import "github.com/gin-gonic/gin"

type SupplyHandler interface {
	CreateSupply(c *gin.Context)
	GetAllSupply(c *gin.Context)
	GetSupplyByID(c *gin.Context)
	GetSupplyByCommodityID(c *gin.Context)
	GetSupplyByRegionID(c *gin.Context)
	UpdateSupply(c *gin.Context)
	DeleteSupply(c *gin.Context)
	GetSupplyHistoryByCommodityIDAndRegionID(c *gin.Context)
}
