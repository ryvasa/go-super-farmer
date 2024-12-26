package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthUtil interface {
	GetAuthUserID(c *gin.Context) (uuid.UUID, error)
	GetAuthRole(c *gin.Context) (string, error)
}

type AuthUtilImpl struct{}

func NewAuthUtil() AuthUtil { // Mengembalikan interface, bukan pointer langsung
	return &AuthUtilImpl{}
}

func (a *AuthUtilImpl) GetAuthUserID(c *gin.Context) (uuid.UUID, error) {
	value, exists := c.Get("user")
	if !exists {
		return uuid.UUID{}, NewUnauthorizedError("unauthorized")
	}
	claims, ok := value.(jwt.MapClaims)
	if !ok {
		return uuid.UUID{}, NewUnauthorizedError("invalid claims type")
	}

	userId, ok := claims["sub"].(string)
	if !ok {
		return uuid.UUID{}, NewUnauthorizedError("invalid user id")
	}
	id, err := uuid.Parse(userId)
	if err != nil {
		return uuid.UUID{}, NewInternalError(err.Error())
	}

	return id, nil
}

func (a *AuthUtilImpl) GetAuthRole(c *gin.Context) (string, error) {
	value, exists := c.Get("user")
	if !exists {
		return "", NewUnauthorizedError("unauthorized")
	}
	claims, ok := value.(jwt.MapClaims)
	if !ok {
		return "", NewUnauthorizedError("invalid claims type")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", NewUnauthorizedError("invalid role")
	}

	return role, nil
}
