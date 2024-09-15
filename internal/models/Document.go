package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Document struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Title     string             `json:"title"`
	CreatedAt time.Time          `json:"-"`
	UserID    primitive.ObjectID `json:"user_id"`
	Moderated bool               `json:"moderated"`
	Subject   string             `json:"subject"`
	Grade     string             `json:"grade"`
}
