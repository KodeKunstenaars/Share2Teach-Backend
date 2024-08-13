package router

import (
	"bytes"
	"fmt"
	"github.com/KodeKunstenaars/Share2Teach/internal/aws/config"
	"github.com/KodeKunstenaars/Share2Teach/internal/aws/s3"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
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

	// Define the file upload endpoint
	router.POST("/upload", func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
			return
		}

		defer file.Close()

		// Read the file content into a buffer
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
			return
		}

		// Create a ReadSeeker from the buffer
		readSeeker := bytes.NewReader(buf.Bytes())

		// Generate the key name for the uploaded file in S3
		key := header.Filename

		// Upload the file to S3
		err = s3.UploadFile(s3Client, "share2teach", key, readSeeker)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("File %s uploaded successfully", key)})
	})

	return router
}
