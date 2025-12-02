package infrastructure

import (
	"task_manager/config"
	"task_manager/domain"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// GenerateToken builds and signs a JWT for the given user
func GenerateToken(user domain.User) (string, error) {
	secret := config.GetEnv("JWT_SECRET")
	if secret == "" {
		secret = "dev_secret"
	}

	claims := jwt.MapClaims{
		"user_id":   user.ID.Hex(),
		"username":  user.Username,
		"user_type": user.UserType,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken parses and validates a JWT token string and returns claims.
func ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	secret := config.GetEnv("JWT_SECRET")
	if secret == "" {
		secret = "dev_secret"
	}

	token, err := jwt.Parse(tokenStr, func(tok *jwt.Token) (interface{}, error) {
		if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenUnverifiable
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenMalformed
}
