package route

import (
	"github.com/gin-gonic/gin"
	handler "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
)

func UserRouter(r *gin.Engine, userHandler handler.UserHandler) {
	r.POST("/users", userHandler.RegisterUser)
	r.GET("/users", userHandler.GetAllUsers)
	r.GET("/users/:id", userHandler.GetOneUser)
	r.PATCH("/users/:id", userHandler.UpdateUser)
	r.DELETE("/users/:id", userHandler.DeleteUser)
	r.PATCH("/users/:id/restore", userHandler.RestoreUser)
}
