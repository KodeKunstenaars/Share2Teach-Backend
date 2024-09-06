package main

import (
	"backend/internal/models"
	"errors"
	"net/http"

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

	// validate payload against database
	existingUser, err := app.DB.GetUserByEmail(payload.Email)
	if err != nil || existingUser != nil {
		app.errorJSON(w, errors.New("payload already exists"), http.StatusBadRequest)
		return
	}

	err = app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// hash password
	// hashedPassword, err := hashPassword(payload.Password)
	// if err != nil {
	// 	app.errorJSON(w, err, http.StatusInternalServerError)
	// 	return
	// }

	// create a new payload
	newUser := &models.User{
		FirstName:     payload.FirstName,
		LastName:      payload.LastName,
		Email:         payload.Email,
		Password:      payload.Password,
		Role:          payload.Role,
		Qualification: payload.Qualification,
	}

	// save payload to database
	err = app.DB.RegisterUser(newUser)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// return success
	app.writeJSON(w, http.StatusCreated, newUser)
}

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	// read json payload
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate user against database
	user, err := app.DB.GetUserByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// check password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// create a jwt user
	u := jwtUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	// generate tokens
	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	app.writeJSON(w, http.StatusAccepted, tokens)
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
				app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}

			// Convert the Subject (userID) to primitive.ObjectID
			userID, err := primitive.ObjectIDFromHex(claims.Subject)
			if err != nil {
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			user, err := app.DB.GetUserByID(userID)
			if err != nil {
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			u := jwtUser{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			}

			tokenPairs, err := app.auth.GenerateTokenPair(&u)
			if err != nil {
				app.errorJSON(w, errors.New("error generating token"), http.StatusUnauthorized)
				return
			}

			http.SetCookie(w, app.auth.GetRefreshCookie(tokenPairs.RefreshToken))
			app.writeJSON(w, http.StatusOK, tokenPairs)
		}
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, app.auth.GetRefreshCookie(""))
	w.WriteHeader(http.StatusAccepted)
}

// ListBuckets lists the buckets in the current account.
func (app *application) listBuckets(w http.ResponseWriter, r *http.Request) {
	// Call the ListBuckets method on your app's Storage
	result, err := app.Storage.ListBuckets()
	if err != nil {
		// Handle error and send a failure response
		var payload = struct {
			Status  string `json:"status"`
			Message string `json:"message"`
			Version string `json:"version"`
		}{
			Status:  "error",
			Message: "Could not list buckets",
			Version: "1.0.0",
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
		Version string   `json:"version"`
	}{
		Status:  "active",
		Message: bucketNames,
		Version: "1.0.0",
	}

	// Write the JSON response
	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		http.Error(w, "Unable to send response", http.StatusInternalServerError)
	}
}
