package handler_interface

import "github.com/gin-gonic/gin"

type ProvinceHandler interface {
	CreateProvince(c *gin.Context)
	GetAllProvinces(c *gin.Context)
	GetProvinceByID(c *gin.Context)
	UpdateProvince(c *gin.Context)
	DeleteProvince(c *gin.Context)
}
