package repositories

import (
	"context"
	"task_manager/config"
	"task_manager/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection

// Init initializes the users collection. Call after DB connection established.
func InitUserCollection() {
	userCollection = config.GetCollection("task_manager_db", "users")
}

func CountByUsername(username string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return userCollection.CountDocuments(ctx, bson.M{"username": username})
}

func CountAllUsers() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return userCollection.CountDocuments(ctx, bson.M{})
}

func InsertUser(u domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if u.ID.IsZero() {
		u.ID = primitive.NewObjectID()
	}
	_, err := userCollection.InsertOne(ctx, u)
	return err
}

func FindByUsername(username string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var u domain.User
	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&u)
	return u, err
}

func PromoteUserByID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res, err := userCollection.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"user_type": "ADMIN"}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
