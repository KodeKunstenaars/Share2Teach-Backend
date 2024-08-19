package s3

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
)

// NewS3Client initializes a new S3 client using the provided session
func NewS3Client(sess *session.Session) *s3.S3 {
	return s3.New(sess)
}

// UploadFile uploads a file to the specified S3 bucket
func UploadFile(client *s3.S3, bucketName, key string, body io.ReadSeeker) error {
	_, err := client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   body,
	})
	return err
}
