package router

import (
	"encoding/gob"
	"fmt"
	"github.com/KodeKunstenaars/Share2Teach/internal/authenticator"
	"github.com/KodeKunstenaars/Share2Teach/internal/aws"
	"github.com/KodeKunstenaars/Share2Teach/internal/middleware"
	"github.com/KodeKunstenaars/Share2Teach/web/app/callback"
	"github.com/KodeKunstenaars/Share2Teach/web/app/login"
	"github.com/KodeKunstenaars/Share2Teach/web/app/user"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"github.com/KodeKunstenaars/Share2Teach/internal/db"
	"github.com/KodeKunstenaars/Share2Teach/web/app/logout"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Initialize AWS session
	sess, err := aws.InitSession()
	if err != nil {
		fmt.Println("Failed to initialize AWS session:", err)
		return nil
	}

	// Initialize S3 client
	s3Client := aws.NewS3Client(sess)

	// Retrieve MongoDB URI from environment variable
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		fmt.Println("MONGODB_URI environment variable not set")
		return nil
	}

	// Initialize MongoDB client with Atlas connection string
	mongoClient, err := db.InitMongo(mongoURI)
	if err != nil {
		fmt.Println("Failed to initialize MongoDB client:", err)
		return nil
	}

	// Select the database and collection
	collection := mongoClient.Database("Share2Teach").Collection("metadata")

	// Initialize upload route
	InitUploadRoutes(router, s3Client, collection)

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

	router.Static("/public", "web/static")
	router.LoadHTMLGlob("web/template/*")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "home.html", nil)
	})
	router.GET("/login", login.Handler(auth))
	router.GET("/callback", callback.Handler(auth))
	//router.GET("/user", user.Handler)
	router.GET("/logout", logout.Handler)
	router.GET("/user", middleware.IsAuthenticated, user.Handler)

	return router
}
