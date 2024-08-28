package UserAccount

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/KodeKunstenaars/Share2Teach/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// User structure
type User struct {
	Email    string `bson:"email" json:"email"`
	Name     string `bson:"name" json:"name"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
	Role     string `bson:"role" json:"role"`
}

var Client *mongo.Client

// OAuth configuration (Google)
var googleOAuthConfig *oauth2.Config

// Function to hash the password using SHA-256
func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

// Function to create a user account
func createUserAccount(client *mongo.Client, email, name, username, password, role string) error {
	collection := client.Database("Share2Teach").Collection("user_info")

	// Check if the user already exists
	var existingUser User
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&existingUser)
	if err == nil {
		return fmt.Errorf("user with email %s already exists", email)
	}

	// Hash the password
	hashedPassword := hashPassword(password)

	// Create a new user instance
	newUser := User{
		Email:    email,
		Name:     name,
		Username: username,
		Password: hashedPassword,
		Role:     role,
	}

	// Insert the new user into the collection
	_, err = collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

// HTTP handler for creating a user account
func createUserAccountHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
			return
		}

		err = createUserAccount(client, user.Email, user.Name, user.Username, user.Password, user.Role)
		if err != nil {
			http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "User account created successfully!")
	}
}

// Function to check user credentials
func checkUserCredentials(client *mongo.Client, email, password string) (bool, error) {
	collection := client.Database("Share2Teach").Collection("user_info")

	// Find the user by email
	var user User
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return false, fmt.Errorf("user not found")
	}

	// Compare hashed passwords
	hashedPassword := hashPassword(password)
	if hashedPassword != user.Password {
		return false, fmt.Errorf("incorrect password")
	}

	return true, nil
}

// HTTP handler for standard login
func loginHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var credentials struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Check user credentials
		isValid, err := checkUserCredentials(client, credentials.Email, credentials.Password)
		if err != nil || !isValid {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Login successful!")
	}
}

// OAuth login redirect handler
func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := googleOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// OAuth callback handler
func googleCallbackHandler(Client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		token, err := googleOAuthConfig.Exchange(context.TODO(), code)
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		client := googleOAuthConfig.Client(context.TODO(), token)
		response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer response.Body.Close()

		var userInfo struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		}
		if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
			http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if the user exists in the database
		var user User
		collection := Client.Database("Share2Teach").Collection("user_info")
		err = collection.FindOne(context.TODO(), bson.M{"email": userInfo.Email}).Decode(&user)
		if err != nil {
			// User does not exist, create a new account
			user = User{
				Email:    userInfo.Email,
				Name:     userInfo.Name,
				Username: userInfo.Email, // You can modify this logic for the username
				Role:     "user",
			}
			_, err := collection.InsertOne(context.TODO(), user)
			if err != nil {
				http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Google login successful! Welcome, %s", userInfo.Name)
	}
}

func Create() {
	uri := "your_mongodb_uri"

	// Call InitMongo and handle both return values
	client, err := db.InitMongo(uri)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB client: %v", err)
	}

	// Ensure MongoDB connection is established
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	log.Println("MongoDB connection established successfully!")

	// Now you can use 'client' to interact with your MongoDB
	// OAuth configuration for Google
	googleOAuthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	// HTTP Handlers
	http.HandleFunc("/createUser", createUserAccountHandler(client))
	http.HandleFunc("/login", loginHandler(client))
	http.HandleFunc("/googleLogin", googleLoginHandler)
	http.HandleFunc("/callback", googleCallbackHandler(client))

	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
