package route

import (
	"github.com/gin-gonic/gin"
	handler "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
)

func ProvinceRoute(public, protected *gin.RouterGroup, handler handler.ProvinceHandler) {
	protected.POST("/provinces", handler.CreateProvince)
	protected.GET("/provinces", handler.GetAllProvinces)
	protected.GET("/provinces/:id", handler.GetProvinceById)
	protected.PATCH("/provinces/:id", handler.UpdateProvince)
	protected.DELETE("/provinces/:id", handler.DeleteProvince)
}
