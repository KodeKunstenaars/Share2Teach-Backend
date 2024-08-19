package createUser

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID       string `bson:"_id,omitempty"`
	Username string `bson:"username"`
	Password string `bson:"password"`
}

var userCollection *mongo.Collection

// Initialize the MongoDB connection and collection (assuming you have a connection already set up)
func init() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	userCollection = client.Database("Share2Teach").Collection("users") //get the name of the log in
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

// Function to get user input from the console
func getUserInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func main() {
	// Get user input for account creation
	username := getUserInput("Enter username: ")
	password := getUserInput("Enter password: ")

	// Create user account
	err := createUser(username, password)
	if err != nil {
		fmt.Println("Error creating user:", err)
	} else {
		fmt.Println("User created successfully")
	}

	// Get user input for login
	loginUsername := getUserInput("Enter username for login: ")
	loginPassword := getUserInput("Enter password for login: ")

	// Attempt login
	isAuthenticated, err := loginUser(loginUsername, loginPassword)
	if err != nil {
		fmt.Println("Login failed:", err)
	} else if isAuthenticated {
		fmt.Println("User logged in successfully")
	} else {
		fmt.Println("Invalid login credentials")
	}
}
