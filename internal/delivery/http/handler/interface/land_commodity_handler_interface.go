package handler_interface

import "github.com/gin-gonic/gin"

type LandCommodityHandler interface {
	CreateLandCommodity(c *gin.Context)
	GetLandCommodityByID(c *gin.Context)
	GetLandCommodityByLandID(c *gin.Context)
	GetAllLandCommodity(c *gin.Context)
	GetLandCommodityByCommodityID(c *gin.Context)
	UpdateLandCommodity(c *gin.Context)
	DeleteLandCommodity(c *gin.Context)
	RestoreLandCommodity(c *gin.Context)
}
