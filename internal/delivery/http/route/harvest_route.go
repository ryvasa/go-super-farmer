package route

import (
	"github.com/gin-gonic/gin"
	handler_interface "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/interface"
)

type HarvestRoute struct {
	handler handler_interface.HarvestHandler
}

func NewHarvestRoute(handler handler_interface.HarvestHandler) *HarvestRoute {
	return &HarvestRoute{handler}
}

func (r *HarvestRoute) Register(public, protected *gin.RouterGroup) {
	protected.POST("/harvests", r.handler.CreateHarvest)
	protected.GET("/harvests", r.handler.GetAllHarvest)
	protected.GET("/harvests/:id", r.handler.GetHarvestByID)
	protected.GET("/harvests/commodity/:id", r.handler.GetHarvestByCommodityID)
	protected.GET("/harvests/land/:id", r.handler.GetHarvestByLandID)
	protected.GET("/harvests/land_commodity/:id", r.handler.GetHarvestByLandCommodityID)
	protected.GET("/harvests/city/:id", r.handler.GetHarvestByCityID)
	protected.PATCH("/harvests/:id", r.handler.UpdateHarvest)
	protected.DELETE("/harvests/:id", r.handler.DeleteHarvest)
	protected.PATCH("/harvests/:id/restore", r.handler.RestoreHarvest)
	protected.GET("/harvests/deleted", r.handler.GetAllDeletedHarvest)
	protected.GET("/harvests/deleted/:id", r.handler.GetHarvestDeletedByID)
	public.GET("/harvests/land_commodity/:id/download", r.handler.DownloadHarvestByLandCommodityID)
	public.GET("/harvests/land_commodity/:id/download/file", r.handler.GetHarvestExcelFile)
}
