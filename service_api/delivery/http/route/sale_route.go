package route

import (
	"github.com/gin-gonic/gin"
	handler_interface "github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler/interface"
)

type SaleRoute struct {
	handler handler_interface.SaleHandler
}

func NewSaleRoute(handler handler_interface.SaleHandler) *SaleRoute {
	return &SaleRoute{handler}
}

func (r *SaleRoute) Register(public, protected *gin.RouterGroup) {
	protected.POST("/sales", r.handler.CreateSale)
	protected.GET("/sales", r.handler.GetAllSales)
	protected.GET("/sales/:id", r.handler.GetSaleByID)
	protected.GET("/sales/commodity/:id", r.handler.GetSalesByCommodityID)
	protected.GET("/sales/city/:id", r.handler.GetSalesByCityID)
	protected.PUT("/sales/:id", r.handler.UpdateSale)
	protected.DELETE("/sales/:id", r.handler.DeleteSale)
	protected.POST("/sales/:id/restore", r.handler.RestoreSale)
	protected.GET("/sales/deleted", r.handler.GetAllDeletedSales)
	protected.GET("/sales/deleted/:id", r.handler.GetDeletedSaleByID)
}
