package upload

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
	"path/filepath"
	"time"

	customS3 "github.com/KodeKunstenaars/Share2Teach/internal/aws/s3"
	"github.com/KodeKunstenaars/Share2Teach/internal/duplicatechecker"
	"github.com/KodeKunstenaars/Share2Teach/internal/models"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// InitUploadRoutes initializes the routes for file uploads
func InitUploadRoutes(router *gin.Engine, s3Client *s3.S3, checker *duplicatechecker.Checker, collection *mongo.Collection) {
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
			err = customS3.UploadFile(s3Client, "share2teach", s3Key, readSeeker)
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
			UploadedAt: time.Now().UTC().Add(2 * time.Hour),
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
}
