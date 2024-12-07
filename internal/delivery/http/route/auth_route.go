package route

import (
	"github.com/gin-gonic/gin"
	handler "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
)

func AuthRoutes(public *gin.RouterGroup, authHandler handler.AuthHandler) {
	public.POST("/auth/login", authHandler.Login)
}
