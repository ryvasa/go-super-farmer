package route

import (
	"github.com/gin-gonic/gin"
	handler "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
)

func RoleRoutes(public, protected *gin.RouterGroup, roleHandler handler.RoleHandler) {
	protected.POST("/roles", roleHandler.CreateRole)
	protected.GET("/roles", roleHandler.GetAllRoles)
}
