package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Rating struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	DocID         primitive.ObjectID `json:"doc_id" bson:"doc_id"`
	TimesRated    int                `json:"times_rated" bson:"times_rated"`
	TotalRating   int                `json:"total_rating" bson:"total_rating"`
	AverageRating float64            `json:"average_rating" bson:"average_rating"`
}
