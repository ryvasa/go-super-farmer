package route

import (
	"github.com/gin-gonic/gin"
	handler "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
)

func HarvestRoutes(public, protected *gin.RouterGroup, harvestHandler handler.HarvestHandler) {
	protected.POST("/harvests", harvestHandler.CreateHarvest)
	protected.GET("/harvests", harvestHandler.GetAllHarvest)
	protected.GET("/harvests/:id", harvestHandler.GetHarvestByID)
	protected.GET("/harvests/commodity/:id", harvestHandler.GetHarvestByCommodityID)
	protected.GET("/harvests/land/:id", harvestHandler.GetHarvestByLandID)
	protected.GET("/harvests/land_commodity/:id", harvestHandler.GetHarvestByLandCommodityID)
	protected.GET("/harvests/region/:id", harvestHandler.GetHarvestByRegionID)
	protected.PATCH("/harvests/:id", harvestHandler.UpdateHarvest)
	protected.DELETE("/harvests/:id", harvestHandler.DeleteHarvest)
	protected.PATCH("/harvests/:id/restore", harvestHandler.RestoreHarvest)
	protected.GET("/harvests/deleted", harvestHandler.GetAllDeletedHarvest)
	protected.GET("/harvests/deleted/:id", harvestHandler.GetHarvestDeletedByID)
}
