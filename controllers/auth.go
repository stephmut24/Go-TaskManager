package controllers

import (
	"task_manager/config"
	"task_manager/models"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// GenerateJWT creates a signed JWT containing basic user info and expiry
func GenerateJWT(user models.User) (string, error) {
	secret := config.GetEnv("JWT_SECRET")
	if secret == "" {
		// fallback secret for local testing (not secure) — set JWT_SECRET in production
		secret = "dev_secret"
	}

	claims := jwt.MapClaims{
		"user_id":   user.ID.Hex(),
		"username":  user.UserName,
		"user_type": user.User_type,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}
