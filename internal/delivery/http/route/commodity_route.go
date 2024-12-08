package route

import (
	"github.com/gin-gonic/gin"
	handler "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
)

func CommodityRoutes(public, protected *gin.RouterGroup, commodityHandler handler.CommodityHandler) {
	protected.POST("/commodities", commodityHandler.CreateCommodity)
	protected.GET("/commodities", commodityHandler.GetAllCommodities)
	protected.GET("/commodities/:id", commodityHandler.GetCommodityById)
	protected.PATCH("/commodities/:id", commodityHandler.UpdateCommodity)
	protected.DELETE("/commodities/:id", commodityHandler.DeleteCommodity)
	protected.PATCH("/commodities/:id/restore", commodityHandler.RestoreCommodity)
}
