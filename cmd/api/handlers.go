package main

import (
	"backend/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (app *application) uploadDocumentMetadata(w http.ResponseWriter, r *http.Request) {
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
		DocumentID primitive.ObjectID `json:"document_id"`
		Title      string             `json:"title"`
		Subject    string             `json:"subject"`
		Grade      string             `json:"grade"`
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
		ID:        payload.DocumentID,
		Title:     payload.Title,
		CreatedAt: time.Now(),
		UserID:    userID,
		Moderated: false,
		Subject:   payload.Subject,
		Grade:     payload.Grade,
		Reported:  false,
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

func (app *application) generatePresignedURLForUpload(w http.ResponseWriter, r *http.Request) {
	documentID := primitive.NewObjectID()
	objectKey := fmt.Sprintf( /*"documents/%s",*/ documentID.Hex()) // Object key for S3

	// Generate the presigned URL for the client to upload the document
	presignedRequest, err := app.Storage.PutObject("share2teach", objectKey, 3600)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("error generating presigned URL: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the presigned URL to the client along with the document ID
	response := struct {
		DocumentID   primitive.ObjectID `json:"document_id"`
		PresignedURL string             `json:"presigned_url"`
	}{
		DocumentID:   documentID,
		PresignedURL: presignedRequest.URL,
	}

	err = app.writeJSON(w, http.StatusOK, response)
	if err != nil {
		return
	}
}

func (app *application) searchDocuments(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	subject := r.URL.Query().Get("subject")
	grade := r.URL.Query().Get("grade")
	correctRole := false

	// finds the documents that match the given title
	documents, err := app.DB.FindDocuments(title, subject, grade, correctRole)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("error finding documents: %v", err), http.StatusInternalServerError)
		log.Println("error finding documents:", err)
		return
	}

	if len(documents) == 0 {
		app.errorJSON(w, fmt.Errorf("no documents found"), http.StatusNotFound)
		return
	}

	err = app.writeJSON(w, http.StatusOK, documents)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (app *application) searchDocumentsAdminOrModerator(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	subject := r.URL.Query().Get("subject")
	grade := r.URL.Query().Get("grade")
	correctRole := true

	// finds the documents that match the given title
	documents, err := app.DB.FindDocuments(title, subject, grade, correctRole)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("error finding documents: %v", err), http.StatusInternalServerError)
		log.Println("error finding documents:", err)
		return
	}

	if len(documents) == 0 {
		app.errorJSON(w, fmt.Errorf("no documents found"), http.StatusNotFound)
		return
	}

	err = app.writeJSON(w, http.StatusOK, documents)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (app *application) generatePresignedURLForDownload(w http.ResponseWriter, r *http.Request) {
	documentIDStr := chi.URLParam(r, "id")

	// Ensure documentID is provided
	if documentIDStr == "" {
		app.errorJSON(w, fmt.Errorf("document ID is missing"), http.StatusBadRequest)
		return
	}

	// Convert the documentID from string to primitive.ObjectID
	documentID, err := primitive.ObjectIDFromHex(documentIDStr)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("invalid document ID: %v", err), http.StatusBadRequest)
		return
	}

	objectKey := fmt.Sprintf(documentID.Hex())

	// Generate the presigned URL for the S3 object using the document ID
	presignedRequest, err := app.Storage.GetObject("share2teach", objectKey, 3600)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("error generating presigned URL: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the presigned URL and document ID to the client
	response := struct {
		PresignedURL string `json:"presigned_url"`
	}{
		PresignedURL: presignedRequest.URL,
	}

	// Respond with the JSON data
	err = app.writeJSON(w, http.StatusOK, response)
	if err != nil {
		return
	}
}

func (app *application) FAQs(w http.ResponseWriter, r *http.Request) {

	faqs, err := app.DB.GetFAQs()
	if err != nil {
		http.Error(w, "Failed to fetch FAQs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(faqs)
}
