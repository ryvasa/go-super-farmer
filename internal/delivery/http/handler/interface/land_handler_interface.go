package handler_interface

import "github.com/gin-gonic/gin"

type LandHandler interface {
	CreateLand(c *gin.Context)
	GetLandByID(c *gin.Context)
	GetLandByUserID(c *gin.Context)
	GetAllLands(c *gin.Context)
	UpdateLand(c *gin.Context)
	DeleteLand(c *gin.Context)
	RestoreLand(c *gin.Context)
}
