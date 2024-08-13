package router

import (
	"context"
	"fmt"
	"github.com/KodeKunstenaars/Share2Teach/internal/aws/config"
	"github.com/KodeKunstenaars/Share2Teach/internal/aws/s3"
	"github.com/KodeKunstenaars/Share2Teach/internal/db"
	"github.com/KodeKunstenaars/Share2Teach/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
	"time"
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

	// Define the file upload endpoint
	router.POST("/upload", func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
			return
		}

		defer file.Close()

		// Generate a new ObjectID for the file metadata
		metadataID := primitive.NewObjectID()

		// Prepare the metadata with the ID
		metadata := models.FileMetadata{
			ID:         metadataID.Hex(), // Convert ObjectID to a string for storage
			Filename:   header.Filename,
			Bucket:     "share2teach",
			Key:        metadataID.Hex(), // Use the same ID as the S3 key
			Size:       header.Size,
			UploadedAt: time.Now(),
		}

		// Upload the file to S3 using the metadata ID as the key
		err = s3.UploadFile(s3Client, "share2teach", metadataID.Hex(), file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
			return
		}

		// Store metadata in MongoDB
		_, err = collection.InsertOne(context.Background(), metadata)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store metadata"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("File %s uploaded and metadata stored successfully", header.Filename)})
	})

	return router
}
