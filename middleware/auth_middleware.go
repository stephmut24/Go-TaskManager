package middleware

import (
	"errors"
	"net/http"
	"strings"

	"task_manager/config"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates JWT in the Authorization header and sets claims in the context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		tokenString := parts[1]
		secret := config.GetEnv("JWT_SECRET")
		if secret == "" {
			secret = "dev_secret"
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// ensure token method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// copy relevant claim fields into context for handlers
			if v, exists := claims["user_id"]; exists {
				c.Set("user_id", v)
			}
			if v, exists := claims["username"]; exists {
				c.Set("username", v)
			}
			if v, exists := claims["user_type"]; exists {
				c.Set("user_type", v)
			}
		}

		c.Next()
	}
}

// AdminOnly middleware enforces the user_type claim equals "ADMIN". Must run after AuthMiddleware.
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if v, ok := c.Get("user_type"); ok {
			if vs, ok2 := v.(string); ok2 && vs == "ADMIN" {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin role required"})
	}
}
