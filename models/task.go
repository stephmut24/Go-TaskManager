package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



type Task struct {
	ID  			primitive.ObjectID 	`bson:"id"`
	Title  			string 				`bson:"title"`
	Description  	string 				`bson:"description"`
	DueDate 		time.Time  			`bson:"due_date"`
	Status 			string 				`bson:"status"`
}

