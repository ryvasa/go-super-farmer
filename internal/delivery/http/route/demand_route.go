package route

import (
	"github.com/gin-gonic/gin"
	handler "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
)

func DemandRoutes(public, protected *gin.RouterGroup, demandHandler handler.DemandHandler) {
	protected.POST("/demands", demandHandler.CreateDemand)
	protected.GET("/demands", demandHandler.GetAllDemands)
	protected.GET("/demands/:id", demandHandler.GetDemandByID)
	protected.GET("/demands/commodity/:commodity_id", demandHandler.GetDemandsByCommodityID)
	protected.GET("/demands/region/:id", demandHandler.GetDemandsByRegionID)
	protected.PATCH("/demands/:id", demandHandler.UpdateDemand)
	protected.DELETE("/demands/:id", demandHandler.DeleteDemand)
	protected.GET("/demands/commodity/:commodity_id/region/:region_id", demandHandler.GetDemandHistoryByCommodityIDAndRegionID)
}
