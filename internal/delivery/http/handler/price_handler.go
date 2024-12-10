package handler

import (
	"github.com/gin-gonic/gin"
)

type PriceHandler interface {
	CreatePrice(c *gin.Context)
	GetAllPrices(c *gin.Context)
	GetPriceById(c *gin.Context)
	GetPricesByCommodityID(c *gin.Context)
	GetPricesByRegionID(c *gin.Context)
	UpdatePrice(c *gin.Context)
	DeletePrice(c *gin.Context)
	RestorePrice(c *gin.Context)
}
