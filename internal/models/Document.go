package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

type Document struct {
    ID             primitive.ObjectID `json:"_id" bson:"_id"`
    Title          string             `json:"title" bson:"title"`
    CreatedAt      time.Time          `json:"-" bson:"created_at"`
    UserID         primitive.ObjectID `json:"user_id" bson:"user_id"`
    Moderated      bool               `json:"moderated" bson:"moderated"`
    Subject        string             `json:"subject" bson:"subject"`
    Grade          string             `json:"grade" bson:"grade"`
    Reported       bool               `json:"reported" bson:"reported"`
    RatingID       primitive.ObjectID `json:"rating_id" bson:"rating_id"`
    ApprovalStatus string             `json:"approvalStatus" bson:"approvalStatus"` 
}
