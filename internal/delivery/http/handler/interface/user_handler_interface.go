package handler_interface

import "github.com/gin-gonic/gin"

type UserHandler interface {
	RegisterUser(c *gin.Context)
	GetOneUser(c *gin.Context)
	GetAllUsers(c *gin.Context)
	DeleteUser(c *gin.Context)
	RestoreUser(c *gin.Context)
	UpdateUser(c *gin.Context)
}
