package route

import (
	"github.com/gin-gonic/gin"
	handler "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
)

func PriceRoutes(public, protected *gin.RouterGroup, priceHandler handler.PriceHandler) {
	protected.POST("/prices", priceHandler.CreatePrice)
	protected.GET("/prices", priceHandler.GetAllPrices)
	public.GET("/prices/:id", priceHandler.GetPriceByID)
	protected.GET("/prices/commodity/:id", priceHandler.GetPricesByCommodityID)
	protected.GET("/prices/region/:id", priceHandler.GetPricesByRegionID)
	protected.PATCH("/prices/:id", priceHandler.UpdatePrice)
	protected.DELETE("/prices/:id", priceHandler.DeletePrice)
	protected.PATCH("/prices/:id/restore", priceHandler.RestorePrice)
	public.GET("/prices/current/commodity/:commodity_id/region/:region_id", priceHandler.GetPriceByCommodityIDAndRegionID)
	public.GET("/prices/history/commodity/:commodity_id/region/:region_id", priceHandler.GetPricesHistoryByCommodityIDAndRegionID)
}
