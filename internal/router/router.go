package router

import (
	"fmt"
	"github.com/KodeKunstenaars/Share2Teach/internal/aws"
	"os"

	"github.com/KodeKunstenaars/Share2Teach/internal/db"
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

	return router
}
