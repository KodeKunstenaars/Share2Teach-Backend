package router

import (
	"fmt"
	"os"

	"github.com/KodeKunstenaars/Share2Teach/internal/aws/config"
	"github.com/KodeKunstenaars/Share2Teach/internal/aws/s3"
	"github.com/KodeKunstenaars/Share2Teach/internal/db"
	"github.com/KodeKunstenaars/Share2Teach/internal/duplicatechecker"
	"github.com/KodeKunstenaars/Share2Teach/internal/router/upload"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Initialize AWS session
	sess, err := config.InitSession()
	if err != nil {
		fmt.Println("Failed to initialize AWS session:", err)
		return nil
	}

	// Initialize S3 client
	s3Client := s3.NewS3Client(sess)

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

	// Create a new duplicate checker
	checker := duplicatechecker.NewChecker(collection)

	// Initialize upload routes
	upload.InitUploadRoutes(router, s3Client, checker, collection)

	return router
}
