package route

import (
	"github.com/gin-gonic/gin"
	handler_interface "github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler/interface"
)

type UserRoute struct {
	handler handler_interface.UserHandler
}

func NewUserRoute(handler handler_interface.UserHandler) *UserRoute {
	return &UserRoute{handler}
}

func (r *UserRoute) Register(public, protected *gin.RouterGroup) {
	public.POST("/users", r.handler.RegisterUser)
	protected.GET("/users", r.handler.GetAllUsers)
	protected.GET("/users/:id", r.handler.GetOneUser)
	protected.PATCH("/users/:id", r.handler.UpdateUser)
	protected.DELETE("/users/:id", r.handler.DeleteUser)
	protected.PATCH("/users/:id/restore", r.handler.RestoreUser)
}
