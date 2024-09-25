package dbrepo

import (
	"backend/internal/models"
	"backend/pkg/db"
	"context"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBRepo struct {
	//database           db.Database
	userInfoCollection      db.Collection
	passwordResetCollection db.Collection
	metadataCollection      db.Collection
	ratingsCollection       db.Collection
	faqsCollection          db.Collection
	moderateCollection      db.Collection
	reportsCollection       db.Collection
}

func NewMongoDBRepo(client *mongo.Client, databaseName string) *MongoDBRepo {
	database := client.Database(databaseName)
	return &MongoDBRepo{
		userInfoCollection:      database.Collection("user_info"),
		passwordResetCollection: database.Collection("password_reset"),
		metadataCollection:      database.Collection("metadata"),
		ratingsCollection:       database.Collection("ratings"),
		faqsCollection:          database.Collection("faqs"),
		moderateCollection:      database.Collection("moderate"),
		reportsCollection:       database.Collection("reports"),
	}
}

const dbTimeout = time.Second * 3

func (m *MongoDBRepo) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := m.userInfoCollection

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

	collection := m.userInfoCollection

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

	collection := m.userInfoCollection

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoDBRepo) UploadDocumentMetadata(document *models.Document) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := m.metadataCollection

	_, err := collection.InsertOne(ctx, document)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoDBRepo) CreateDocumentRating(initialRating *models.Rating) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := m.ratingsCollection

	_, err := collection.InsertOne(ctx, initialRating)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoDBRepo) FindDocuments(title, subject, grade string, correctRole bool) ([]models.Document, error) {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel() // ensures that the context is canceled after the function returns

	collection := m.metadataCollection

	// Creates a filter for the query that only searches for the given parameters
	// Search query for the parameters that is case-insensitive
	// Creates a filter for the query that only searches for the given parameters
	// Search query for the parameters that is case-insensitive
	filter := bson.M{}

	if title != "" {
		filter["title"] = bson.M{"$regex": primitive.Regex{Pattern: title, Options: "i"}}
	}
	if subject != "" {
		filter["subject"] = bson.M{"$regex": primitive.Regex{Pattern: subject, Options: "i"}}
	}
	if grade != "" {

		normalizedGrade := strings.ToLower(strings.TrimSpace(grade))
		normalizedGrade = strings.ReplaceAll(normalizedGrade, " ", "")
		normalizedGrade = strings.TrimPrefix(normalizedGrade, "grade")

		filter["$or"] = []bson.M{
			{"grade": normalizedGrade},
			{"grade": bson.M{"$regex": primitive.Regex{Pattern: normalizedGrade, Options: "i"}}},
		}

	}
	if correctRole == true {
		filter["moderated"] = bson.M{"$in": []bool{true, false}}
		//filter["reported"] = bson.M{"$in": []bool{true, false}}
	} else {
		filter["moderated"] = true
		filter["reported"] = false
	}

	// Cursor that loops through the DB to find the matching documents
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Slice that will hold all documents that match the filter
	var documents []models.Document

	// Loops through the filter and adds the documents to the slice
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

func (m *MongoDBRepo) GetFAQs() ([]models.FAQs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := m.faqsCollection

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var faqs []models.FAQs

	for cursor.Next(ctx) {
		var faq models.FAQs
		if err := cursor.Decode(&faq); err != nil {
			return nil, err
		}
		faqs = append(faqs, faq)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return faqs, nil
}

// UpdateDocumentsByID updates the document data by ID.
func (m *MongoDBRepo) UpdateDocumentsByID(documentID primitive.ObjectID, updateData bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := m.metadataCollection

	filter := bson.M{"_id": documentID}
	// No need for another $set here, we assume updateData has the correct update format
	update := updateData

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating document with ID %s: %v", documentID.Hex(), err)
		return err
	}

	return nil
}

// InsertModerationData inserts the moderation data into the "moderate" collection.
func (m *MongoDBRepo) InsertModerationData(userID, documentID primitive.ObjectID, approvalStatus, comments string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := m.moderateCollection

	moderationData := bson.M{
		"moderatedBy":    userID,
		"documentID":     documentID,
		"approvalStatus": approvalStatus,
		"comments":       comments,
		"moderatedAt":    time.Now(),
	}

	_, err := collection.InsertOne(ctx, moderationData)
	if err != nil {
		return err
	}

	return nil
}

// GetDocumentByID retrieves a document by its ID.
func (m *MongoDBRepo) GetDocumentByID(id primitive.ObjectID) (*models.Document, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := m.metadataCollection

	var document models.Document
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&document)
	if err != nil {
		return nil, err
	}

	return &document, nil
}

// SetDocumentRating inserts or updates the rating for a given document by its ID.
func (m *MongoDBRepo) SetDocumentRating(docID primitive.ObjectID, rating *models.Rating) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collectionForDocument := m.metadataCollection
	var document models.Document

	err := collectionForDocument.FindOne(ctx, bson.M{"_id": docID}).Decode(&document)
	if err != nil {
		return err
	}

	collectionForRating := m.ratingsCollection
	filter := bson.M{"_id": document.RatingID}
	update := bson.M{
		"$inc": bson.M{"times_rated": 1, "total_rating": rating.TotalRating},
		"$set": bson.M{"average_rating": bson.M{"$divide": []interface{}{"$total_rating", "$times_rated"}}},
	}

	_, err = collectionForRating.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

// GetDocumentRating retrieves the rating for a given document by its ID.
func (m *MongoDBRepo) GetDocumentRating(docID primitive.ObjectID) (*models.Rating, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := m.ratingsCollection

	var rating models.Rating
	err := collection.FindOne(ctx, bson.M{"doc_id": docID}).Decode(&rating)
	if err != nil {
		return nil, err
	}

	return &rating, nil
}

func (m *MongoDBRepo) StoreResetToken(resetEntry *models.PasswordReset) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := m.passwordResetCollection

	_, err := collection.InsertOne(ctx, resetEntry)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoDBRepo) VerifyResetToken(userID primitive.ObjectID, token string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Query the password_reset collection using "user_id"
	collection := m.passwordResetCollection
	filter := bson.M{"user_id": userID, "reset_token": token, "spent": false}

	var resetEntry models.PasswordReset
	err := collection.FindOne(ctx, filter).Decode(&resetEntry)
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if time.Now().After(resetEntry.ExpiresAt) {
		return false, nil
	}

	return true, nil
}

// ChangeUserPassword updates the password for a given user by their ID.
func (m *MongoDBRepo) ChangeUserPassword(userID primitive.ObjectID, newPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := m.userInfoCollection

	filter := bson.M{"_id": userID}
	update := bson.M{"$set": bson.M{"password": newPassword}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	resetCollection := m.passwordResetCollection

	resetFilter := bson.M{"user_id": userID}
	resetUpdate := bson.M{"$set": bson.M{"spent": true}}

	_, err = resetCollection.UpdateOne(ctx, resetFilter, resetUpdate)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoDBRepo) InsertReport(report bson.M) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	collection := m.reportsCollection
	result, err := collection.InsertOne(ctx, report)
	if err != nil {
		return nil, err
	}

	return result, nil
}
