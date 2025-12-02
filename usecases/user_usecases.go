package usecases

import (
	"errors"
	"task_manager/domain"
	"task_manager/infrastructure"
	"task_manager/repositories"
)

// RegisterUser performs validation, hashing and stores the user. If the DB is empty, the user is created as ADMIN.
func RegisterUser(user domain.User) error {
	// ensure username uniqueness
	count, err := repositories.CountByUsername(user.Username)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("username already exists")
	}

	total, err := repositories.CountAllUsers()
	if err != nil {
		return err
	}

	if total == 0 {
		user.UserType = "ADMIN"
	} else if user.UserType == "" {
		user.UserType = "USER"
	}

	hashed, err := infrastructure.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed

	return repositories.InsertUser(user)
}

// LoginUser authenticates and returns a token + user
func LoginUser(username, password string) (string, domain.User, error) {
	u, err := repositories.FindByUsername(username)
	if err != nil {
		return "", domain.User{}, errors.New("invalid username or password")
	}

	if err := infrastructure.ComparePassword(u.Password, password); err != nil {
		return "", domain.User{}, errors.New("invalid username or password")
	}

	token, err := infrastructure.GenerateToken(u)
	if err != nil {
		return "", domain.User{}, err
	}
	return token, u, nil
}

func PromoteUser(id string) error { return repositories.PromoteUserByID(id) }
