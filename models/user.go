package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserName  string             `json:"username" validate:"required,min=4,max=100"`
	Password  string             `json:"password" validate:"required,min=6"`
	User_type string             `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	CreatedAt time.Time          `json:"created_at"`
}
