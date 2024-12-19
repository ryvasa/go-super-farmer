package route

import (
	"github.com/gin-gonic/gin"
	handler_interface "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/interface"
)

type PriceRoute struct {
	handler handler_interface.PriceHandler
}

func NewPriceRoute(handler handler_interface.PriceHandler) *PriceRoute {
	return &PriceRoute{handler}
}

func (r *PriceRoute) Register(public, protected *gin.RouterGroup) {
	protected.POST("/prices", r.handler.CreatePrice)
	protected.GET("/prices", r.handler.GetAllPrices)
	public.GET("/prices/:id", r.handler.GetPriceByID)
	protected.GET("/prices/commodity/:id", r.handler.GetPricesByCommodityID)
	protected.GET("/prices/region/:id", r.handler.GetPricesByRegionID)
	protected.PATCH("/prices/:id", r.handler.UpdatePrice)
	protected.DELETE("/prices/:id", r.handler.DeletePrice)
	protected.PATCH("/prices/:id/restore", r.handler.RestorePrice)
	public.GET("/prices/current/commodity/:commodity_id/region/:region_id", r.handler.GetPriceByCommodityIDAndRegionID)
	public.GET("/prices/history/commodity/:commodity_id/region/:region_id", r.handler.GetPricesHistoryByCommodityIDAndRegionID)
	public.GET("/prices/history/commodity/:commodity_id/region/:region_id/download", r.handler.DownloadPricesHistoryByCommodityIDAndRegionID)
	public.GET("/prices/history/commodity/:commodity_id/region/:region_id/download/file", r.handler.GetPriceHistoryExcelFile)
}
