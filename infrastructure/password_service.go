package infrastructure

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// ComparePassword compares a bcrypt hashed password with its possible plaintext equivalent.
// Returns nil on success, or an error on failure.
func ComparePassword(hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
