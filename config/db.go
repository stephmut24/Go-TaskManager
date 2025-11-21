package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client 

func ConnectDB() {
	dbURI := GetEnv("DB_URL")
	if dbURI == "" {
		log.Fatal("DB_URL missing in .env")
	}
	
	clientOptions := options.Client().ApplyURI(dbURI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB: ",err)
	}
	//check connection
	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal("Unable to ping MongoDB : ",err)
	}
	fmt.Println("MongoDB connected successfully")
	DB = client
}

func GetCollection(dbName, collName string) *mongo.Collection {
	if DB == nil {
		log.Fatal("MongoDB client is not initialized. Call ConnectDB() first.")
	}
	return DB.Database(dbName).Collection(collName)
}