package router

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/KodeKunstenaars/Share2Teach/internal/aws/config"
	"github.com/KodeKunstenaars/Share2Teach/internal/aws/s3"
	"github.com/KodeKunstenaars/Share2Teach/internal/db"
	"github.com/KodeKunstenaars/Share2Teach/internal/duplicatechecker"
	"github.com/KodeKunstenaars/Share2Teach/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

		defer func(file multipart.File) {
			if err := file.Close(); err != nil {
				fmt.Println("Failed to close file:", err)
			}
		}(file)

		// Read the file content into a buffer
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file content"})
			return
		}

		// Detect the file type based on the content
		kind, _ := filetype.Match(fileBytes)
		var fileType string
		if kind == filetype.Unknown {
			// Fallback to checking the file extension if type is unknown
			ext := filepath.Ext(header.Filename)
			fileType = mime.TypeByExtension(ext)
			if fileType == "" {
				fileType = "application/octet-stream" // Final fallback if extension is also unrecognized
			}
		} else {
			fileType = kind.MIME.Value
		}

		// Calculate the hash
		hash := sha256.New()
		hash.Write(fileBytes)
		fileHash := hex.EncodeToString(hash.Sum(nil))

		// Generate an ObjectID for the first-time upload
		objectID := primitive.NewObjectID().Hex()

		// Check for duplicate
		existingMetadata, err := checker.GetMetadataByHash(context.Background(), fileHash)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for duplicates"})
			return
		}

		var s3Key string
		var mongoID string

		if existingMetadata != nil {
			// Duplicate found, reuse the existing S3 key
			s3Key = existingMetadata.Key
			mongoID = primitive.NewObjectID().Hex() // Generate a new MongoDB ID
		} else {
			// No duplicate, use the same ObjectID for both MongoDB `_id` and S3 key
			s3Key = objectID
			mongoID = objectID

			// Upload the file to S3 using the ObjectID as the key
			readSeeker := bytes.NewReader(fileBytes)
			err = s3.UploadFile(s3Client, "share2teach", s3Key, readSeeker)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
				return
			}
		}

		// Prepare the metadata with the ID, hash, and file type
		metadata := models.FileMetadata{
			ID:         mongoID, // Use the same ID for both MongoDB and S3 key on first upload
			Filename:   header.Filename,
			FileType:   fileType,
			Bucket:     "share2teach",
			Key:        s3Key, // Reuse the existing S3 key if duplicate, or use the new one
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
