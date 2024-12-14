package route

import (
	"github.com/gin-gonic/gin"
	handler "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
)

func SupplyRoutes(public, protected *gin.RouterGroup, supplyHandler handler.SupplyHandler) {
	protected.POST("/supplies", supplyHandler.CreateSupply)
	protected.GET("/supplies", supplyHandler.GetAllSupply)
	protected.GET("/supplies/:id", supplyHandler.GetSupplyByID)
	protected.GET("/supplies/commodity/:commodity_id", supplyHandler.GetSupplyByCommodityID)
	protected.GET("/supplies/region/:id", supplyHandler.GetSupplyByRegionID)
	protected.PATCH("/supplies/:id", supplyHandler.UpdateSupply)
	protected.DELETE("/supplies/:id", supplyHandler.DeleteSupply)
	protected.GET("/supplies/commodity/:commodity_id/region/:region_id", supplyHandler.GetSupplyHistoryByCommodityIDAndRegionID)
}
