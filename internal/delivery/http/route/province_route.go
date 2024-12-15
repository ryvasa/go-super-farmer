package route

import (
	"github.com/gin-gonic/gin"
	handler_interface "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/interface"
)

type ProvinceRoute struct {
	handler handler_interface.ProvinceHandler
}

func NewProvinceRoute(handler handler_interface.ProvinceHandler) *ProvinceRoute {
	return &ProvinceRoute{handler}
}

func (r *ProvinceRoute) Register(public, protected *gin.RouterGroup) {
	protected.POST("/provinces", r.handler.CreateProvince)
	protected.GET("/provinces", r.handler.GetAllProvinces)
	protected.GET("/provinces/:id", r.handler.GetProvinceByID)
	protected.PATCH("/provinces/:id", r.handler.UpdateProvince)
	protected.DELETE("/provinces/:id", r.handler.DeleteProvince)
}
