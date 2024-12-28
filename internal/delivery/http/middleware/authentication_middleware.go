package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ryvasa/go-super-farmer/pkg/auth/token"
	"github.com/ryvasa/go-super-farmer/utils"
)

type AuthMiddleware struct {
	token token.Token
	// enforcer  *casbin.Enforcer
}

func NewAuthMiddleware(token token.Token) *AuthMiddleware {
	return &AuthMiddleware{token}
}

func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authentication
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, utils.NewUnauthorizedError("no authorization header"))
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.ErrorResponse(c, utils.NewUnauthorizedError("invalid authorization format"))
			c.Abort()
			return
		}
		claims, err := m.token.ExtractClaims(tokenParts[1])
		if err != nil {
			utils.ErrorResponse(c, utils.NewUnauthorizedError("invalid token"))
			c.Abort()
			return
		}
		c.Set("user", claims)
		c.Next()
	}
}
