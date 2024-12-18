package handler_interface

import (
	"github.com/gin-gonic/gin"
)

type PriceHandler interface {
	CreatePrice(c *gin.Context)
	GetAllPrices(c *gin.Context)
	GetPriceByID(c *gin.Context)
	GetPricesByCommodityID(c *gin.Context)
	GetPricesByRegionID(c *gin.Context)
	UpdatePrice(c *gin.Context)
	DeletePrice(c *gin.Context)
	RestorePrice(c *gin.Context)
	GetPriceByCommodityIDAndRegionID(c *gin.Context)
	GetPricesHistoryByCommodityIDAndRegionID(c *gin.Context)
	DownloadPricesHistoryByCommodityIDAndRegionID(c *gin.Context)
	GetPriceHistoryExcelFile(c *gin.Context)
}
