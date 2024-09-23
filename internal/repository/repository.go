package repository

import (
	"backend/internal/models"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DatabaseRepo interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id primitive.ObjectID) (*models.User, error)
	RegisterUser(user *models.User) error
	UploadDocumentMetadata(document *models.Document) error
	FindDocuments(title, subject, grade string, correctRole bool) ([]models.Document, error)
	GetFAQs() ([]models.FAQs, error)
	GetDocumentByID(id primitive.ObjectID) (*models.Document, error)
	GetDocumentRating(id primitive.ObjectID) (*models.Rating, error)
	SetDocumentRating(id primitive.ObjectID, rating *models.Rating) error
	CreateDocumentRating(rating *models.Rating) error
	StoreResetToken(resetKey *models.PasswordReset) error
	VerifyResetToken(id primitive.ObjectID, resetKey string) (bool, error)
	ChangeUserPassword(id primitive.ObjectID, newPassword string) error
	UpdateDocumentsByID(documentID primitive.ObjectID, updateData bson.M) error
	InsertModerationData(userID, documentID primitive.ObjectID, approvalStatus, comments string) error
}

type StorageRepo interface {
	ListBuckets() ([]types.Bucket, error)
	BucketExists(bucketName string) (bool, error)
	CreateBucket(name string, region string) error
	PutObject(bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error)
	GetObject(bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error)
}

type MailRepo interface {
	SendPasswordResetRequest(email, token string) error
	SendWelcomeEmail(email string, firstName string, lastName string) error
}
