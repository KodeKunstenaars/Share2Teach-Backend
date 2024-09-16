package fileapproval

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/KodeKunstenaars/Share2Teach/internal/db" // Import the db package

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var fileCollection *mongo.Collection

func init() {
	// Initialize the MongoDB client and get the collection
	client, err := db.InitMongo("your_mongo_uri")
	if err != nil {
		panic(err)
	}
	fileCollection = client.Database("your_database_name").Collection("your_collection_name")
}

func ModerateFile(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to extract file ID, status, and comments
	type RequestBody struct {
		ID             string `json:"id"`
		ApprovalStatus string `json:"approvalStatus"`
		Comments       string `json:"comments,omitempty"`
		Moderator      string `json:"moderator"`
	}
	var reqBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	fileID, err := primitive.ObjectIDFromHex(reqBody.ID)
	if err != nil {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	// Create the update document
	update := bson.M{
		"$set": bson.M{
			"approvalStatus": reqBody.ApprovalStatus,
			"comments":       reqBody.Comments,
			"moderator":      reqBody.Moderator,
			"updatedAt":      time.Now(),
		},
	}

	// Update the document in MongoDB
	filter := bson.M{"_id": fileID}
	opts := options.Update().SetUpsert(false)
	_, err = fileCollection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		http.Error(w, "Failed to update file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File moderated successfully"))
}
