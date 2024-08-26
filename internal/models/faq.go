package models

// FileMetadata represents the metadata for a file
type FAQ struct {
	ID       string `bson:"_id,omitempty"`
	Question string `bson:"question"`
	Answer   string `bson:"answer"`
}
