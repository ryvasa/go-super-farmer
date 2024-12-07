package route

import (
	"github.com/gin-gonic/gin"
	handler "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
)

func RoleRoutes(public *gin.RouterGroup, roleHandler handler.RoleHandler) {
	public.POST("/roles", roleHandler.CreateRole)
	public.GET("/roles", roleHandler.GetAllRoles)
}
