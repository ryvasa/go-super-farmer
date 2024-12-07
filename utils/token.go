package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/pkg/env"
)

type TokenUtil interface {
	GenerateToken(id uuid.UUID, role string) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	ExtractClaims(tokenString string) (jwt.MapClaims, error)
}

type TokenUtilImpl struct {
	env *env.Env
}

func NewTokenUtil(cfg *env.Env) TokenUtil {
	return &TokenUtilImpl{
		env: cfg,
	}
}

func (t *TokenUtilImpl) GenerateToken(id uuid.UUID, role string) (string, error) {
	claims := jwt.MapClaims{
		"iss":  "go-restaurant-api",
		"sub":  id,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
		"role": role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.env.Secret.JwtSecretKey))
}

func (t *TokenUtilImpl) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.env.Secret.JwtSecretKey), nil
	})
}

func (t *TokenUtilImpl) ExtractClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := t.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, NewInternalError("Invalid token claims")
}
