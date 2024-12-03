package handler

import "github.com/gin-gonic/gin"

type UserHandler interface {
	RegisterUser(c *gin.Context)
	GetUser(c *gin.Context)
}
