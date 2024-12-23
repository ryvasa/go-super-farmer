package route

import (
	"github.com/gin-gonic/gin"
	handler_interface "github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler/interface"
)

type LandCommodityRoute struct {
	handler handler_interface.LandCommodityHandler
}

func NewLandCommodityRoute(handler handler_interface.LandCommodityHandler) *LandCommodityRoute {
	return &LandCommodityRoute{handler}
}

func (r *LandCommodityRoute) Register(public, protected *gin.RouterGroup) {
	protected.POST("/land_commodities", r.handler.CreateLandCommodity)
	protected.GET("/land_commodities", r.handler.GetAllLandCommodity)
	protected.GET("/land_commodities/:id", r.handler.GetLandCommodityByID)
	protected.GET("/land_commodities/land/:id", r.handler.GetLandCommodityByLandID)
	protected.GET("/land_commodities/commodity/:id", r.handler.GetLandCommodityByCommodityID)
	protected.PATCH("/land_commodities/:id", r.handler.UpdateLandCommodity)
	protected.DELETE("/land_commodities/:id", r.handler.DeleteLandCommodity)
	protected.PATCH("/land_commodities/:id/restore", r.handler.RestoreLandCommodity)
}
