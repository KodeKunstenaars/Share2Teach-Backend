package createUser

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	//"fmt"
	//"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	Client *mongo.Client
	// Replace with your actual Google OAuth credentials
	oauthConfig = &oauth2.Config{
		ClientID:     "<YOUR_GOOGLE_CLIENT_ID>",
		ClientSecret: "<YOUR_GOOGLE_CLIENT_SECRET>",
		RedirectURL:  "http://localhost:8080/auth/google/callback", //can only change when home screen front end is available
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}
	oauthStateString = "random" // Generate a secure state string in production
)

// User model
type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email    string             `json:"email" bson:"email"`
	Name     string             `json:"name" bson:"name"`
	Role     string             `json:"role" bson:"role"`
	Username string             `json:"username" bson:"username"`
	GoogleID string             `json:"googleId" bson:"googleId"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

func init() {
	// Load environment variables

	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	// Initialize MongoDB client
	mongoURI := os.Getenv("MONGO_URI")
	var err error
	Client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI)) //call connect
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	// Routes
	r.POST("/register", registerUser)
	r.POST("/login", loginUser)
	r.GET("/auth/google", handleGoogleLogin)
	r.GET("/auth/google/callback", handleGoogleCallback)
	r.GET("/profile", profile)

	r.Run(":8080")
}

func registerUser(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	collection := Client.Database("Share2Teach").Collection("user_info")
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func loginUser(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := Client.Database("Share2Teach").Collection("user_info")
	var user User
	err := collection.FindOne(context.Background(), bson.M{"email": credentials.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func handleGoogleLogin(c *gin.Context) {
	url := oauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func handleGoogleCallback(c *gin.Context) {
	code := c.DefaultQuery("code", "")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	client := oauthConfig.Client(context.Background(), token)
	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer userInfo.Body.Close()

	var user User
	if err := json.NewDecoder(userInfo.Body).Decode(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	collection := Client.Database("Share2Teach").Collection("user_info")
	var existingUser User
	err = collection.FindOne(context.Background(), bson.M{"googleId": user.GoogleID}).Decode(&existingUser)
	if err != nil {
		_, err = collection.InsertOne(context.Background(), user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, user)
}

func profile(c *gin.Context) {
	// For demonstration, just returning a static response.
	// Implement user authentication and retrieval logic as needed.
	c.JSON(http.StatusOK, gin.H{"message": "Profile data"})
}
