package route

import (
	"github.com/gin-gonic/gin"
	handler_interface "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/interface"
)

type SupplyRoute struct {
	handler handler_interface.SupplyHandler
}

func NewSupplyRoute(handler handler_interface.SupplyHandler) *SupplyRoute {
	return &SupplyRoute{handler}
}

func (r *SupplyRoute) Register(public, protected *gin.RouterGroup) {
	protected.POST("/supplies", r.handler.CreateSupply)
	protected.GET("/supplies", r.handler.GetAllSupply)
	protected.GET("/supplies/:id", r.handler.GetSupplyByID)
	protected.GET("/supplies/commodity/:commodity_id", r.handler.GetSupplyByCommodityID)
	protected.GET("/supplies/region/:id", r.handler.GetSupplyByRegionID)
	protected.PATCH("/supplies/:id", r.handler.UpdateSupply)
	protected.DELETE("/supplies/:id", r.handler.DeleteSupply)
	protected.GET("/supplies/commodity/:commodity_id/region/:region_id", r.handler.GetSupplyHistoryByCommodityIDAndRegionID)
}
