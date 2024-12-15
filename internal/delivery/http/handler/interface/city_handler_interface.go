package handler_interface

import (
	"github.com/gin-gonic/gin"
)

type CityHandler interface {
	CreateCity(c *gin.Context)
	GetAllCities(c *gin.Context)
	GetCityByID(c *gin.Context)
	UpdateCity(c *gin.Context)
	DeleteCity(c *gin.Context)
}
