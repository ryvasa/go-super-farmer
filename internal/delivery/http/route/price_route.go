package route

import (
	"github.com/gin-gonic/gin"
	handler "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
)

func PriceRoutes(public, protected *gin.RouterGroup, priceHandler handler.PriceHandler) {
	public.POST("/prices", priceHandler.CreatePrice)
	protected.GET("/prices", priceHandler.GetAllPrices)
	public.GET("/prices/:id", priceHandler.GetPriceById)
	protected.GET("/prices/commodity/:id", priceHandler.GetPricesByCommodityID)
	protected.GET("/prices/region/:id", priceHandler.GetPricesByRegionID)
	protected.PATCH("/prices/:id", priceHandler.UpdatePrice)
	protected.DELETE("/prices/:id", priceHandler.DeletePrice)
	protected.PATCH("/prices/:id/restore", priceHandler.RestorePrice)
}
