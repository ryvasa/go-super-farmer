package handler_interface

import "github.com/gin-gonic/gin"

type RoleHandler interface {
	CreateRole(c *gin.Context)
	GetAllRoles(c *gin.Context)
}
