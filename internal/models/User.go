package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName     string             `json:"first_name"`
	LastName      string             `json:"last_name"`
	Email         string             `json:"email"`
	Password      string             `json:"password"`
	Role          string             `json:"role"`
	Qualification string             `json:"qualification"`
}

func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
