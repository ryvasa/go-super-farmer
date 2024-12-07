package route

import (
	"github.com/gin-gonic/gin"
	handler "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
)

func LandRoutes(public, protected *gin.RouterGroup, userHandler handler.LandHandler) {

	protected.GET("/lands", userHandler.GetAllLands)
	protected.POST("/lands", userHandler.CreateLand)
	protected.GET("/lands/:id", userHandler.GetLandByID)
	protected.PATCH("/lands/:id", userHandler.UpdateLand)
	protected.DELETE("/lands/:id", userHandler.DeleteLand)
	protected.PATCH("/lands/:id/restore", userHandler.RestoreLand)
	protected.GET("/lands/user/:id", userHandler.GetLandByUserID)
}
