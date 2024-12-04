package route

import (
	"github.com/gin-gonic/gin"
	handler "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/user"
)

func NewRouter(userHandler handler.UserHandler) *gin.Engine {
	r := gin.Default()
	r.POST("/users", userHandler.RegisterUser)
	r.GET("/users/:id", userHandler.GetOneUser)
	return r
}
