package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type NewsPost struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Slug        string             `json:"slug" bson:"slug"`
	Content     string             `json:"content" bson:"content"`
	Summary     string             `json:"summary" bson:"summary"`
	Image       string             `json:"image" bson:"image"`
	Author      string             `json:"author" bson:"author"`
	PublishedAt primitive.DateTime `json:"publishedAt" bson:"publishedAt"`
	CreatedAt   primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt   primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}
