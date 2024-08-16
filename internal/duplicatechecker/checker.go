package duplicatechecker

import (
	"context"
	"fmt"
	"github.com/KodeKunstenaars/Share2Teach/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// Checker defines the structure for the duplicate checker
type Checker struct {
	collection *mongo.Collection
}

// NewChecker creates a new instance of Checker
func NewChecker(collection *mongo.Collection) *Checker {
	return &Checker{collection: collection}
}

// GetMetadataByHash checks if a file with the given hash already exists in the database
// and returns the metadata if it does.
func (dc *Checker) GetMetadataByHash(ctx context.Context, fileHash string) (*models.FileMetadata, error) {
	filter := map[string]interface{}{"hash": fileHash}

	var result models.FileMetadata
	err := dc.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No document found, so it's not a duplicate
		}
		return nil, fmt.Errorf("error checking for duplicate: %w", err)
	}

	return &result, nil
}
