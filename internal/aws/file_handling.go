package upload_file

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"time"
)

func Upload(filePath string) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Unable to open file %v", err)
	}
	defer file.Close()

	// Get the file name from the provided file path
	fileName := filepath.Base(filePath)

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Unable to load AWS SDK config %v", err)
	}

	// Create an S3 client
	client := s3.NewFromConfig(cfg)

	// Create an uploader with the S3 client
	uploader := manager.NewUploader(client)

	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	// Upload the file to S3
	result, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String("share2teach"),
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		log.Fatalf("Unable to upload file to S3 %v", err)
	}

	fmt.Printf("File uploaded to S3 successfully. Location: %s\n", result.Location)
}
