package route

import (
	"github.com/gin-gonic/gin"
	handler_interface "github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler/interface"
)

type CommodityRoute struct {
	handler handler_interface.CommodityHandler
}

func NewCommodityRoute(handler handler_interface.CommodityHandler) *CommodityRoute {
	return &CommodityRoute{handler}
}

func (r *CommodityRoute) Register(public, protected *gin.RouterGroup) {
	protected.POST("/commodities", r.handler.CreateCommodity)
	protected.GET("/commodities", r.handler.GetAllCommodities)
	protected.GET("/commodities/:id", r.handler.GetCommodityById)
	protected.PATCH("/commodities/:id", r.handler.UpdateCommodity)
	protected.DELETE("/commodities/:id", r.handler.DeleteCommodity)
	protected.PATCH("/commodities/:id/restore", r.handler.RestoreCommodity)
}
