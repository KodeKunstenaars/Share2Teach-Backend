package repository

import (
	"backend/internal/models"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DatabaseRepo interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id primitive.ObjectID) (*models.User, error)
	RegisterUser(user *models.User) error
}

type StorageRepo interface {
	ListBuckets() ([]types.Bucket, error)
	BucketExists(bucketName string) (bool, error)
	CreateBucket(name string, region string) error
}
