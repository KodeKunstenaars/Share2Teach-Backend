package models

//import "go.mongodb.org/mongo-driver/bson/primitive"

type FAQs struct {
	Question string `bson:"question" json:"question"`
	Answer   string `bson:"answer" json:"answer"`
}
