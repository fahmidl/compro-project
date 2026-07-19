package models

type NewsPost struct {
	ID          string `json:"id" dynamodbav:"id"`
	Title       string `json:"title" dynamodbav:"title"`
	Slug        string `json:"slug" dynamodbav:"slug"`
	Content     string `json:"content" dynamodbav:"content"`
	Summary     string `json:"summary" dynamodbav:"summary"`
	Image       string `json:"image" dynamodbav:"image"`
	Author      string `json:"author" dynamodbav:"author"`
	PublishedAt int64  `json:"publishedAt" dynamodbav:"publishedAt"`
	CreatedAt   int64  `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt" dynamodbav:"updatedAt"`
}
