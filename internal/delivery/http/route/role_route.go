package route

import (
	"github.com/gin-gonic/gin"
	handler_interface "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/interface"
)

type RoleRoute struct {
	handler handler_interface.RoleHandler
}

func NewRoleRoute(handler handler_interface.RoleHandler) *RoleRoute {
	return &RoleRoute{handler}
}

func (r *RoleRoute) Register(public, protected *gin.RouterGroup) {
	protected.POST("/roles", r.handler.CreateRole)
	protected.GET("/roles", r.handler.GetAllRoles)
}
