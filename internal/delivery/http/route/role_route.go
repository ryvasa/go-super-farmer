package route

import (
	"github.com/gin-gonic/gin"
	handler "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
)

func RoleRouter(r *gin.Engine, roleHandler handler.RoleHandler) {
	r.POST("/roles", roleHandler.CreateRole)
	r.GET("/roles", roleHandler.GetAllRoles)
}
