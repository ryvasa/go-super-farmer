package handler_interface

import "github.com/gin-gonic/gin"

type HarvestHandler interface {
	CreateHarvest(c *gin.Context)
	GetAllHarvest(c *gin.Context)
	GetHarvestByID(c *gin.Context)
	GetHarvestByCommodityID(c *gin.Context)
	GetHarvestByLandID(c *gin.Context)
	GetHarvestByLandCommodityID(c *gin.Context)
	GetHarvestByRegionID(c *gin.Context)
	UpdateHarvest(c *gin.Context)
	DeleteHarvest(c *gin.Context)
	RestoreHarvest(c *gin.Context)
	GetAllDeletedHarvest(c *gin.Context)
	GetHarvestDeletedByID(c *gin.Context)
	DownloadHarvestByLandCommodityID(c *gin.Context)
	GetHarvestExcelFile(c *gin.Context)
}
