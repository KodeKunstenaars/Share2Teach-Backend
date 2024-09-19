package dbrepo

import (
	"backend/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBRepo struct {
	Client   *mongo.Client
	Database string
}

const dbTimeout = time.Second * 3

func (m *MongoDBRepo) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := m.Client.Database(m.Database).Collection("user_info")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *MongoDBRepo) GetUserByID(id primitive.ObjectID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := m.Client.Database(m.Database).Collection("user_info")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *MongoDBRepo) RegisterUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := m.Client.Database(m.Database).Collection("user_info")

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoDBRepo) UploadDocumentMetadata(document *models.Document) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := m.Client.Database(m.Database).Collection("metadata")

	_, err := collection.InsertOne(ctx, document)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoDBRepo) FindDocumentsByTitle(title string) ([]models.Document, error) {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel() //ensures that the context is canceled after the funtion returns

	collection := m.Client.Database(m.Database).Collection("metadata")

	//creates a filter for the query that only searches for the given title
	// search query for title that is case-insensitive
	filter := bson.M{"title": bson.M{"$regex": primitive.Regex{Pattern: title, Options: "i"}}}

	// cursor that loops through the DB to find the matching documents
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// slice that will hold all documents that match the filter
	var documents []models.Document

	// loops through the filter and adds the documents to the slice
	for cursor.Next(ctx) {
		var doc models.Document
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}

		documents = append(documents, doc)
	}

	if err := cursor.Err(); err != nil {

		return nil, err
	}

	return documents, nil
}
