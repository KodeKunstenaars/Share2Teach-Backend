package models

import "time"

// FileMetadata represents the metadata for a file
type FileMetadata struct {
	ID         string    `bson:"_id,omitempty"`
	Filename   string    `bson:"filename"`
	FileType   string    `bson:"filetype"`
	Subject    string    `bson:"subject"`
	Grade      int       `bson:"grade"`
	Bucket     string    `bson:"bucket"`
	Key        string    `bson:"key"`
	Size       int64     `bson:"size"`
	UploadedAt time.Time `bson:"uploaded_at"`
	Hash       string    `bson:"hash"`
	User       string    `bson:"user"` // To be added later
}
