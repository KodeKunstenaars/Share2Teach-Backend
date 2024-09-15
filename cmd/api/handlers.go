package main

import (
	"backend/internal/models"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Go Share2Teach up and running",
		Version: "1.0.0",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		FirstName     string `json:"first_name"`
		LastName      string `json:"last_name"`
		Email         string `json:"email"`
		Password      string `json:"password"`
		Role          string `json:"role"`
		Qualification string `json:"qualification"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		err := app.errorJSON(w, err, http.StatusBadRequest)
		if err != nil {
			return
		}
		return
	}

	hashedPassword, err := models.HashPassword(payload.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		err := app.errorJSON(w, err, http.StatusInternalServerError)
		if err != nil {
			return
		}
		return
	}

	newUser := &models.User{
		ID:            primitive.NewObjectID(),
		FirstName:     payload.FirstName,
		LastName:      payload.LastName,
		Email:         payload.Email,
		Password:      hashedPassword,
		Role:          payload.Role, //"educator",
		Qualification: payload.Qualification,
	}

	err = app.DB.RegisterUser(newUser)
	if err != nil {
		log.Printf("Error inserting user into MongoDB: %v", err)
		err := app.errorJSON(w, err, http.StatusInternalServerError)
		if err != nil {
			return
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, newUser)
	if err != nil {
		return
	}
}

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	// read json payload
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		err := app.errorJSON(w, err, http.StatusBadRequest)
		if err != nil {
			return
		}
		return
	}

	// validate user against database
	user, err := app.DB.GetUserByEmail(requestPayload.Email)
	if err != nil {
		err := app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		if err != nil {
			return
		}
		return
	}

	// check password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		err := app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		if err != nil {
			return
		}
		return
	}

	// create a jwt user
	u := jwtUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	}

	// generate tokens
	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
		return
	}

	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	err = app.writeJSON(w, http.StatusAccepted, tokens)
	if err != nil {
		return
	}
}

func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == app.auth.CookieName {
			claims := &Claims{}
			refreshToken := cookie.Value

			// parse the token to get the claims
			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(app.JWTSecret), nil
			})
			if err != nil {
				err := app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
				if err != nil {
					return
				}
				return
			}

			// Convert the Subject (userID) to primitive.ObjectID
			userID, err := primitive.ObjectIDFromHex(claims.Subject)
			if err != nil {
				err := app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				if err != nil {
					return
				}
				return
			}

			user, err := app.DB.GetUserByID(userID)
			if err != nil {
				err := app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				if err != nil {
					return
				}
				return
			}

			u := jwtUser{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Role:      user.Role,
			}

			tokenPairs, err := app.auth.GenerateTokenPair(&u)
			if err != nil {
				err := app.errorJSON(w, errors.New("error generating token"), http.StatusUnauthorized)
				if err != nil {
					return
				}
				return
			}

			http.SetCookie(w, app.auth.GetRefreshCookie(tokenPairs.RefreshToken))
			err = app.writeJSON(w, http.StatusOK, tokenPairs)
			if err != nil {
				return
			}
		}
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, app.auth.GetRefreshCookie(""))
	w.WriteHeader(http.StatusAccepted)
}

func (app *application) listBuckets(w http.ResponseWriter, r *http.Request) {
	result, err := app.Storage.ListBuckets()
	if err != nil {
		var payload = struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{
			Status:  "error",
			Message: "Could not list buckets",
		}

		_ = app.writeJSON(w, http.StatusInternalServerError, payload)
		return
	}

	// Prepare the response payload for successful listing
	var bucketNames []string
	for _, bucket := range result {
		bucketNames = append(bucketNames, *bucket.Name)
	}

	var payload = struct {
		Status  string   `json:"status"`
		Message []string `json:"message"`
	}{
		Status:  "active",
		Message: bucketNames,
	}

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		http.Error(w, "Unable to send response", http.StatusInternalServerError)
	}
}

func (app *application) uploadDocument(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the token
	userIDStr, err := app.auth.GetUserIDFromHeader(w, r)
	if err != nil {
		err := app.errorJSON(w, fmt.Errorf("error extracting user ID from token: %v", err), http.StatusUnauthorized)
		if err != nil {
			return
		}
		return
	}

	// Convert the UserID string to MongoDB ObjectID
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		err := app.errorJSON(w, fmt.Errorf("invalid UserID: %v", err), http.StatusBadRequest)
		if err != nil {
			return
		}
		return
	}

	var payload struct {
		Title string `json:"title"`
	}

	err = app.readJSON(w, r, &payload)
	if err != nil {
		err := app.errorJSON(w, err, http.StatusBadRequest)
		if err != nil {
			return
		}
		return
	}

	newDocument := &models.Document{
		ID:        primitive.NewObjectID(),
		Title:     payload.Title,
		DocHash:   "docHash",
		CreatedAt: time.Now(),
		UserID:    userID,
		Moderated: false,
		Subject:   "Subject",
		Grade:     "Grade",
		AWSKey:    "awsKey",
	}

	err = app.DB.UploadDocumentMetadata(newDocument)
	if err != nil {
		log.Printf("Error inserting document into MongoDB: %v", err)
		err := app.errorJSON(w, err, http.StatusInternalServerError)
		if err != nil {
			return
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, newDocument)
	if err != nil {
		return
	}
}
