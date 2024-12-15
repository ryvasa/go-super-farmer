package route

import (
	"github.com/gin-gonic/gin"
	handler_interface "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/interface"
)

type RegionRoute struct {
	handler handler_interface.RegionHandler
}

func NewRegionRoute(handler handler_interface.RegionHandler) *RegionRoute {
	return &RegionRoute{handler}
}

func (r *RegionRoute) Register(public, protected *gin.RouterGroup) {
	public.GET("/regions", r.handler.GetAllRegions)
	public.GET("/regions/:id", r.handler.GetRegionByID)
	protected.POST("/regions", r.handler.CreateRegion)
}
