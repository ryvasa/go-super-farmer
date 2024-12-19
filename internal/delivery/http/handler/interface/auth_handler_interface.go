package handler_interface

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	Login(c *gin.Context)
	SendOTP(c *gin.Context)
	VerifyOTP(c *gin.Context)
}
