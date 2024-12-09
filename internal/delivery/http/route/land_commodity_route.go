package route

import (
	"github.com/gin-gonic/gin"
	handler "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
)

func LandCommodityRoutes(public, protected *gin.RouterGroup, landCommodityHandler handler.LandCommodityHandler) {

	protected.POST("/land_commodities", landCommodityHandler.CreateLandCommodity)
	protected.GET("/land_commodities/:id", landCommodityHandler.GetLandCommodityByID)
	protected.GET("/land_commodities/land/:id", landCommodityHandler.GetLandCommodityByLandID)
	protected.GET("/land_commodities", landCommodityHandler.GetAllLandCommodity)
	protected.GET("/land_commodities/commodity/:id", landCommodityHandler.GetLandCommodityByCommodityID)
	protected.PATCH("/land_commodities/:id", landCommodityHandler.UpdateLandCommodity)
	protected.DELETE("/land_commodities/:id", landCommodityHandler.DeleteLandCommodity)
	protected.PATCH("/land_commodities/:id/restore", landCommodityHandler.RestoreLandCommodity)
}
