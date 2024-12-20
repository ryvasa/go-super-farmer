package handler_interface

import "github.com/gin-gonic/gin"

type SupplyHandler interface {
	CreateSupply(c *gin.Context)
	GetAllSupply(c *gin.Context)
	GetSupplyByID(c *gin.Context)
	GetSupplyByCommodityID(c *gin.Context)
	GetSupplyByCityID(c *gin.Context)
	UpdateSupply(c *gin.Context)
	DeleteSupply(c *gin.Context)
	GetSupplyHistoryByCommodityIDAndCityID(c *gin.Context)
}
