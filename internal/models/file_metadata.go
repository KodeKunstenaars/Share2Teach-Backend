package models

import "time"

// FileMetadata represents the metadata for a file
type FileMetadata struct {
	ID         string    `bson:"_id,omitempty"` // MongoDB document ID
	Filename   string    `bson:"filename"`      // Original file name
	Bucket     string    `bson:"bucket"`        // S3 bucket name
	Key        string    `bson:"key"`           // S3 object key (same as MongoDB ID)
	Size       int64     `bson:"size"`          // File size in bytes
	UploadedAt time.Time `bson:"uploaded_at"`   // Upload timestamp
}
