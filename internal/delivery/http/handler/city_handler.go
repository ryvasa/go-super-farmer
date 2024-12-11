package handler

import (
	"github.com/gin-gonic/gin"
)

type CityHandler interface {
	CreateCity(c *gin.Context)
	GetAllCities(c *gin.Context)
	GetCityById(c *gin.Context)
	UpdateCity(c *gin.Context)
	DeleteCity(c *gin.Context)
}
