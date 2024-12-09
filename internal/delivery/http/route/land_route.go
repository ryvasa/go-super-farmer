package route

import (
	"github.com/gin-gonic/gin"
	handler "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
)

func LandRoutes(public, protected *gin.RouterGroup, landHandler handler.LandHandler) {

	protected.GET("/lands", landHandler.GetAllLands)
	protected.POST("/lands", landHandler.CreateLand)
	protected.GET("/lands/:id", landHandler.GetLandByID)
	protected.PATCH("/lands/:id", landHandler.UpdateLand)
	protected.DELETE("/lands/:id", landHandler.DeleteLand)
	protected.PATCH("/lands/:id/restore", landHandler.RestoreLand)
	protected.GET("/lands/user/:id", landHandler.GetLandByUserID)
}
