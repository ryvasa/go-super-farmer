package handler

import "github.com/gin-gonic/gin"

type ProvinceHandler interface {
	CreateProvince(c *gin.Context)
	GetAllProvinces(c *gin.Context)
	GetProvinceById(c *gin.Context)
	UpdateProvince(c *gin.Context)
	DeleteProvince(c *gin.Context)
}
