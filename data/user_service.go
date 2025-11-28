package data

import (
	"context"
	"task_manager/config"
	"task_manager/models"
	"time"

	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection

// InitUserCollection should be called after the database connection is established
// (e.g. from main) to initialize the users collection handle.
func InitUserCollection() {
	userCollection = config.GetCollection("task_manager_db", "users")
}

var ErrUserAlreadyExists = errors.New("username already exists")

func Signup(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// check if the username exists
	count, err := userCollection.CountDocuments(ctx, bson.M{"username": user.UserName})
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrUserAlreadyExists
	}

	// If there are no users in the database yet, make the first created user an ADMIN.
	totalUsers, err := userCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if totalUsers == 0 {
		// Ensure the first user is admin
		user.User_type = "ADMIN"
	} else if user.User_type == "" {
		// default role for subsequent users
		user.User_type = "USER"
	}

	//hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.ID = primitive.NewObjectID()

	// save data in mongodb

	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func Login(username, password string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User

	//found user by username
	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, errors.New("invalid username or password")
		}
		return models.User{}, err
	}

	//compare bcrypt password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return models.User{}, errors.New("invalid username or password")
	}

	//return user
	return user, nil
}

// PromoteUser sets the user's role to ADMIN (used by admin users)
func PromoteUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": bson.M{"user_type": "ADMIN"}}

	res, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}
