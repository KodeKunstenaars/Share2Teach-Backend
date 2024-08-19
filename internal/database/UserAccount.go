go get go.mongodb.org/mongo-driver/mongo
go get go.mongodb.org/mongo-driver/mongo/options
go get go.mongodb.org/mongo-driver/mongo/bson
go get golang.org/x/crypto/bcrypt




package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
    "golang.org/x/crypto/bcrypt"
)

// User to represent a user
type User struct {
    Username string `json:"username" bson:"username"`
    Email    string `json:"email" bson:"email"`
    Password string `json:"-" bson:"password"`
}

// MongoDB connection string and database
const (
    mongoURI        = "mongodb://localhost:27017"
    dbName          = "testdb"
    usersCollection = "users"
)

// Initialize the MongoDB client
func initDB() *mongo.Client {
    client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }

    // Check the connection
    err = client.Ping(ctx, readpref.Primary())
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Connected to MongoDB!")
    return client
}

// HashPassword hashes the given password
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

// CreateAccountHandler handles the user registration
func CreateAccountHandler(client *mongo.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var user User
        if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        // Hash the user's password
        hashedPassword, err := HashPassword(user.Password)
        if err != nil {
            http.Error(w, "Failed to hash password", http.StatusInternalServerError)
            return
        }
        user.Password = hashedPassword

        collection := client.Database(dbName).Collection(usersCollection)

        // Check if the user already exists
        var existingUser User
        err = collection.FindOne(context.Background(), bson.M{"$or": []bson.M{{"username": user.Username}, {"email": user.Email}}}).Decode(&existingUser)
        if err == nil {
            http.Error(w, "Username or Email already exists", http.StatusBadRequest)
            return
        }

        // Save the user to MongoDB
        _, err = collection.InsertOne(context.Background(), user)
        if err != nil {
            http.Error(w, "Failed to create account", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{"message": "Account created successfully!"})
    }
}

func main() {
    client := initDB()
    defer func() {
        if err := client.Disconnect(context.Background()); err != nil {
            log.Fatal(err)
        }
    }()

    http.HandleFunc("/register", CreateAccountHandler(client))
    log.Println("Server started at :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
