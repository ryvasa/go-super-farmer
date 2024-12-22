package route

import (
	"github.com/gin-gonic/gin"
	handler_interface "github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler/interface"
)

type CityRoute struct {
	handler handler_interface.CityHandler
}

func NewCityRoute(handler handler_interface.CityHandler) *CityRoute {
	return &CityRoute{handler}
}

func (r *CityRoute) Register(public, protected *gin.RouterGroup) {
	protected.POST("/cities", r.handler.CreateCity)
	protected.GET("/cities", r.handler.GetAllCities)
	protected.GET("/cities/:id", r.handler.GetCityByID)
	protected.PATCH("/cities/:id", r.handler.UpdateCity)
	protected.DELETE("/cities/:id", r.handler.DeleteCity)
}
