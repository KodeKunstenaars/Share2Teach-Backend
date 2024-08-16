package router

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/KodeKunstenaars/Share2Teach/internal/aws/config"
	"github.com/KodeKunstenaars/Share2Teach/internal/aws/s3"
	"github.com/KodeKunstenaars/Share2Teach/internal/db"
	"github.com/KodeKunstenaars/Share2Teach/internal/duplicatechecker"
	"github.com/KodeKunstenaars/Share2Teach/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
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

	// Create a new duplicate checker
	checker := duplicatechecker.NewChecker(collection)

	// Define the file upload endpoint
	router.POST("/upload", func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
			return
		}

		defer file.Close()

		// Create a new SHA-256 hash object
		hash := sha256.New()

		// Create a buffer to store the file content
		buf := new(bytes.Buffer)
		writer := io.MultiWriter(hash, buf)

		// Copy the file content to the writer (hash and buffer)
		if _, err := io.Copy(writer, file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
			return
		}

		// Get the hash as a hex string
		fileHash := hex.EncodeToString(hash.Sum(nil))

		// Check for duplicate
		existingMetadata, err := checker.GetMetadataByHash(context.Background(), fileHash)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for duplicates"})
			return
		}

		var s3Key string
		if existingMetadata != nil {
			// Duplicate found, use the existing S3 key
			s3Key = existingMetadata.Key
		} else {
			// No duplicate, generate a new ObjectID for the file metadata
			metadataID := primitive.NewObjectID()
			s3Key = metadataID.Hex()

			// Upload the file to S3 using the metadata ID as the key
			readSeeker := bytes.NewReader(buf.Bytes())
			err = s3.UploadFile(s3Client, "share2teach", s3Key, readSeeker)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
				return
			}
		}

		// Prepare the metadata with the ID and hash
		metadata := models.FileMetadata{
			ID:         primitive.NewObjectID().Hex(), // Convert ObjectID to a string for storage
			Filename:   header.Filename,
			Bucket:     "share2teach",
			Key:        s3Key, // Use the existing or new S3 key
			Size:       header.Size,
			UploadedAt: time.Now(),
			Hash:       fileHash, // Store the hash of the file
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
