package route

import (
	"github.com/gin-gonic/gin"
	handler_interface "github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler/interface"
)

type LandRoute struct {
	handler handler_interface.LandHandler
}

func NewLandRoute(handler handler_interface.LandHandler) *LandRoute {
	return &LandRoute{handler}
}

func (r *LandRoute) Register(public, protected *gin.RouterGroup) {
	protected.GET("/lands", r.handler.GetAllLands)
	protected.POST("/lands", r.handler.CreateLand)
	protected.GET("/lands/:id", r.handler.GetLandByID)
	protected.PATCH("/lands/:id", r.handler.UpdateLand)
	protected.DELETE("/lands/:id", r.handler.DeleteLand)
	protected.PATCH("/lands/:id/restore", r.handler.RestoreLand)
	protected.GET("/lands/user/:id", r.handler.GetLandByUserID)
}
