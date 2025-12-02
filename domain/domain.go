package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is the core user entity in the domain layer
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username  string             `bson:"username" json:"username" validate:"required,min=4,max=100"`
	Password  string             `bson:"password" json:"password" validate:"required,min=6"`
	UserType  string             `bson:"user_type" json:"user_type" validate:"required,oneof=ADMIN USER"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// Task is the core task entity in the domain layer
type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	DueDate     time.Time          `bson:"due_date" json:"due_date"`
	Status      string             `bson:"status" json:"status"`
}
