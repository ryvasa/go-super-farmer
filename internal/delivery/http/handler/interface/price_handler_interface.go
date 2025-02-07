package handler_interface

import (
	"github.com/gin-gonic/gin"
)

type PriceHandler interface {
	CreatePrice(c *gin.Context)
	GetAllPrices(c *gin.Context)
	GetPriceByID(c *gin.Context)
	GetPricesByCommodityID(c *gin.Context)
	GetPricesByCityID(c *gin.Context)
	UpdatePrice(c *gin.Context)
	DeletePrice(c *gin.Context)
	RestorePrice(c *gin.Context)
	GetPriceByCommodityIDAndCityID(c *gin.Context)
	GetPricesHistoryByCommodityIDAndCityID(c *gin.Context)
	GetReportPricesHistoryByCommodityIDAndCityID(c *gin.Context)
	DownloadFileReport(c *gin.Context)
}
