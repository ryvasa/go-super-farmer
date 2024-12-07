package route

import (
	"github.com/gin-gonic/gin"
	handler "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
)

func UserRoutes(public, protected *gin.RouterGroup, userHandler handler.UserHandler) {
	public.POST("/users", userHandler.RegisterUser)
	protected.GET("/users", userHandler.GetAllUsers)
	protected.GET("/users/:id", userHandler.GetOneUser)
	protected.PATCH("/users/:id", userHandler.UpdateUser)
	protected.DELETE("/users/:id", userHandler.DeleteUser)
	protected.PATCH("/users/:id/restore", userHandler.RestoreUser)
}
