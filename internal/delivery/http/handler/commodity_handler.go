package handler

import "github.com/gin-gonic/gin"

type CommodityHandler interface {
	CreateCommodity(c *gin.Context)
	GetAllCommodities(c *gin.Context)
	GetCommodityById(c *gin.Context)
	UpdateCommodity(c *gin.Context)
	DeleteCommodity(c *gin.Context)
	RestoreCommodity(c *gin.Context)
}
