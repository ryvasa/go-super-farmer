package route

import (
	"github.com/gin-gonic/gin"
	handler "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
)

func CityRoute(public, protected *gin.RouterGroup, handler handler.CityHandler) {
	protected.POST("/cities", handler.CreateCity)
	protected.GET("/cities", handler.GetAllCities)
	protected.GET("/cities/:id", handler.GetCityByID)
	protected.PATCH("/cities/:id", handler.UpdateCity)
	protected.DELETE("/cities/:id", handler.DeleteCity)
}
