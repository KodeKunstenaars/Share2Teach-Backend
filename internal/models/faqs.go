package models

//import "go.mongodb.org/mongo-driver/bson/primitive"

type FAQs struct {
	ID       string `bson:"_id,omitempty" json:"id,omitempty"`
	Question string `bson:"question" json:"question"`
	Answer   string `bson:"answer" json:"answer"`
}
