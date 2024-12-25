package handler_interface

import "github.com/gin-gonic/gin"

type SaleHandler interface {
	CreateSale(c *gin.Context)
	GetAllSales(c *gin.Context)
	GetSaleByID(c *gin.Context)
	GetSalesByCommodityID(c *gin.Context)
	GetSalesByCityID(c *gin.Context)
	UpdateSale(c *gin.Context)
	DeleteSale(c *gin.Context)
	RestoreSale(c *gin.Context)
	GetAllDeletedSales(c *gin.Context)
	GetDeletedSaleByID(c *gin.Context)
}
