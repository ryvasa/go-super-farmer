package route

import (
	"github.com/gin-gonic/gin"
	handler "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
)

func NewRouter(handler *handler.Handlers) *gin.Engine {
	r := gin.Default()

	RoleRouter(r, handler.RoleHandler)
	UserRouter(r, handler.UserHandler)
	return r

}
