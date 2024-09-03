package models

import "time"

type Document struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	DocHash   string    `json:"dochash"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	UserID    int       `json:"user_id"`
	Moderated bool      `json:"moderated"`
}
