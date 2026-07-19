package models

type Admin struct {
	ID           string `json:"id" dynamodbav:"id"`
	Username     string `json:"username" dynamodbav:"username"`
	PasswordHash string `json:"-" dynamodbav:"passwordHash"`
	Role         string `json:"role" dynamodbav:"role"`
}
