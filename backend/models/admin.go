package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Admin struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username     string             `json:"username" bson:"username"`
	PasswordHash string             `json:"-" bson:"passwordHash"`
	Role         string             `json:"role" bson:"role"`
}
