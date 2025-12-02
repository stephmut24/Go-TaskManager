package infrastructure

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates the Authorization: Bearer <token> header and injects claims into context
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
		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		if v, ok := claims["user_id"]; ok {
			c.Set("user_id", v)
		}
		if v, ok := claims["username"]; ok {
			c.Set("username", v)
		}
		if v, ok := claims["user_type"]; ok {
			c.Set("user_type", v)
		}

		c.Next()
	}
}

// AdminOnly ensures the current user role is ADMIN
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if v, ok := c.Get("user_type"); ok {
			if s, ok2 := v.(string); ok2 && s == "ADMIN" {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin role required"})
	}
}
