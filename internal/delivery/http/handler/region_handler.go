package handler

import "github.com/gin-gonic/gin"

type RegionHandler interface {
	CreateRegion(c *gin.Context)
	GetAllRegions(c *gin.Context)
	GetRegionByID(c *gin.Context)
}
