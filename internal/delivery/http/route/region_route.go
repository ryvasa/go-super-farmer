package route

import (
	"github.com/gin-gonic/gin"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
)

func RegionRoute(public, protected *gin.RouterGroup, handler handler.RegionHandler) {

	public.GET("/regions", handler.GetAllRegions)
	public.GET("/regions/:id", handler.GetRegionByID)
	protected.POST("/regions", handler.CreateRegion)

}
