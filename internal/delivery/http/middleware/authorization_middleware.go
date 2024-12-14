package middleware

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ryvasa/go-super-farmer/utils"
)

type AutzMiddleware struct {
	enforcer *casbin.Enforcer
}

func NewAutzMiddleware(enforcer *casbin.Enforcer) *AutzMiddleware {
	return &AutzMiddleware{enforcer: enforcer}
}

func (m *AutzMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {

		value, exists := c.Get("user")
		if !exists {
			utils.ErrorResponse(c, utils.NewUnauthorizedError("no user found in context"))
			c.Abort()
			return
		}
		claimsMap, ok := value.(jwt.MapClaims)
		if !ok {
			utils.ErrorResponse(c, utils.NewUnauthorizedError("invalid claims type"))
			c.Abort()
			return
		}

		role := claimsMap["role"].(string)
		path := c.Request.URL.Path
		method := c.Request.Method

		log.Println(role, path, method)
		allowed, err := m.enforcer.Enforce(role, path, method)
		if err != nil {
			utils.ErrorResponse(c, utils.NewUnauthorizedError("Authorization check failed"))
			c.Abort()
			return
		}

		if !allowed {
			utils.ErrorResponse(c, utils.NewForbiddenError("Insufficient permissions"))
			c.Abort()
			return
		}

		c.Next()
	}
}
