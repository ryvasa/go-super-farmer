package route

import (
	"github.com/gin-gonic/gin"
	handler_interface "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/interface"
)

type AuthRoute struct {
	handler handler_interface.AuthHandler
}

func NewAuthRoute(handler handler_interface.AuthHandler) *AuthRoute {
	return &AuthRoute{handler}
}

func (r *AuthRoute) Register(public, protected *gin.RouterGroup) {
	public.POST("/auth/login", r.handler.Login)
}
