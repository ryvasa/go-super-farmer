package route

import (
	"github.com/gin-gonic/gin"
	handler_interface "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/interface"
)

type DemandRoute struct {
	handler handler_interface.DemandHandler
}

func NewDemandRoute(handler handler_interface.DemandHandler) *DemandRoute {
	return &DemandRoute{handler}
}

func (r *DemandRoute) Register(public, protected *gin.RouterGroup) {
	protected.POST("/demands", r.handler.CreateDemand)
	protected.GET("/demands", r.handler.GetAllDemands)
	protected.GET("/demands/:id", r.handler.GetDemandByID)
	protected.GET("/demands/commodity/:commodity_id", r.handler.GetDemandsByCommodityID)
	protected.GET("/demands/region/:id", r.handler.GetDemandsByRegionID)
	protected.PATCH("/demands/:id", r.handler.UpdateDemand)
	protected.DELETE("/demands/:id", r.handler.DeleteDemand)
	protected.GET("/demands/commodity/:commodity_id/region/:region_id", r.handler.GetDemandHistoryByCommodityIDAndRegionID)
}
