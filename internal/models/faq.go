package models

// FAQs for faq page
type FAQ struct {
	ID       string `bson:"_id,omitempty"`
	Question string `bson:"question"`
	Answer   string `bson:"answer"`
}
