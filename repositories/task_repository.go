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

var taskCollection *mongo.Collection

func InitTaskCollection() {
	taskCollection = config.GetCollection("task_manager_db", "tasks")
}

func GetAllTasks() ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := taskCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []domain.Task
	for cursor.Next(ctx) {
		var t domain.Task
		if err := cursor.Decode(&t); err != nil {
			continue
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func GetTaskByID(id string) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var t domain.Task
	if err := taskCollection.FindOne(ctx, bson.M{"_id": oid}).Decode(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

func AddTask(t domain.Task) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if t.ID.IsZero() {
		t.ID = primitive.NewObjectID()
	}
	_, err := taskCollection.InsertOne(ctx, t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func UpdateTask(id string, updated domain.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{}
	if updated.Title != "" {
		update["title"] = updated.Title
	}
	if updated.Description != "" {
		update["description"] = updated.Description
	}
	if !updated.DueDate.IsZero() {
		update["due_date"] = updated.DueDate
	}
	if updated.Status != "" {
		update["status"] = updated.Status
	}

	_, err = taskCollection.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": update})
	return err
}

func DeleteTask(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = taskCollection.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
