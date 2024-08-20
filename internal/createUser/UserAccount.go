package userAccount

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	//try import mongo connection package
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID       string `bson:"_id,omitempty"`
	Username string `bson:"username"`
	Password string `bson:"password"`
}

var userCollection *mongo.Collection

// Initialize the MongoDB connection
func init() {
	client := InitMongo("mongodb://localhost:27017")                    // Call the InitMongo function
	userCollection = client.Database("Share2Teach").Collection("login") // .collection needs to be login
}

// Helper function to hash the password using SHA-256
func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

// Function to create a new user account
func createUser(username, password string) error {
	hashedPassword := hashPassword(password)

	user := User{
		Username: username,
		Password: hashedPassword,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if the user already exists
	var existingUser User
	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&existingUser)
	if err == nil {
		return errors.New("user already exists")
	} else if err != mongo.ErrNoDocuments {
		return err
	}

	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// Function to login a user
func loginUser(username, password string) (bool, error) {
	hashedPassword := hashPassword(password)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := userCollection.FindOne(ctx, bson.M{"username": username, "password": hashedPassword}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, errors.New("invalid username or password")
		}
		return false, err
	}

	return true, nil
}

func main() {
	// Example usage:
	username := "exampleUser"
	password := "examplePassword"

	// Create a user account
	err := createUser(username, password)
	if err != nil {
		fmt.Println("Error creating user:", err)
	} else {
		fmt.Println("User created successfully")
	}

	// Attempt to login
	isAuthenticated, err := loginUser(username, password)
	if err != nil {
		fmt.Println("Login failed:", err)
	} else if isAuthenticated {
		fmt.Println("User logged in successfully")
	} else {
		fmt.Println("Invalid login credentials")
	}
}
