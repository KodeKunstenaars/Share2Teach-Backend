package storagerepo

import "github.com/aws/aws-sdk-go-v2/service/s3"

type StorageRepo struct {
	S3Client      *s3.Client
	PresignClient *s3.PresignClient
}
