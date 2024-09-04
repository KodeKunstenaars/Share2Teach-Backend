package models

import "time"

type Document struct {
	ID        int       `json:"_id" bson:"_id"`
	Title     string    `json:"title"`
	DocHash   string    `json:"doc_hash"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	UserID    int       `json:"user_id"`
	Moderated bool      `json:"moderated"`
}
