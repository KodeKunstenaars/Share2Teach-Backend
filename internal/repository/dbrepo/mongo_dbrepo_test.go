package dbrepo_test

import (
	"backend/internal/models"
	"backend/internal/repository/dbrepo"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MockSingleResult simulates a single result from a MongoDB query
type MockSingleResult struct {
	mock.Mock
}

func (m *MockSingleResult) Decode(v interface{}) error {
	args := m.Called(v)
	return args.Error(0)
}

// MockClient simulates a MongoDB client
type MockClient struct {
	mock.Mock
}

func (m *MockClient) Database(name string, opts ...*options.DatabaseOptions) *mongo.Database {
	args := m.Called(name)
	return args.Get(0).(*mongo.Database) // This line should return a mock of mongo.Database
}

// MockCollection simulates the behavior of *mongo.Collection
type MockCollection struct {
	mock.Mock
}

func (m *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.SingleResult)
}

func (m *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

// MockDatabase simulates the behavior of *mongo.Database
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection {
	args := m.Called(name)
	return args.Get(0).(*mongo.Collection)
}

func setupMocks() (*dbrepo.MongoDBRepo, *MockClient, *MockDatabase, *MockCollection, *MockSingleResult) {
	mockClient := new(MockClient)
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockSingleResult := new(MockSingleResult)

	// Mock the expected return of the Database method
	mockClient.On("Database", "test_db").Return(mockDB)

	// Mock the expected return of the Collection method
	mockDB.On("Collection", "users").Return(mockCollection) // Ensure this matches your collection name

	// Use the interface type MongoClient
	repo := &dbrepo.MongoDBRepo{
		Client:   mockClient, // Pass the MockClient
		Database: "test_db",
	}

	return repo, mockClient, mockDB, mockCollection, mockSingleResult
}

func TestGetUserByEmail(t *testing.T) {
	// Setup mocks
	repo, mockClient, mockDB, mockCollection, mockSingleResult := setupMocks()

	email := "test@example.com"
	expectedUser := &models.User{
		ID:    primitive.NewObjectID(),
		Email: email,
	}

	// Mock the Decode method of MockSingleResult to populate the user
	mockSingleResult.On("Decode", mock.Anything).Run(func(args mock.Arguments) {
		user := args.Get(0).(*models.User)
		*user = *expectedUser // Populate the user with expected data
	}).Return(nil)

	// Mock FindOne to return the expected result
	mockCollection.On("FindOne", mock.Anything, bson.M{"email": email}).Return(mockSingleResult)

	// Call the method under test
	actualUser, err := repo.GetUserByEmail(email)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Email, actualUser.Email) // Compare emails

	// Verify the expectations
	mockClient.AssertExpectations(t)
	mockDB.AssertExpectations(t)
	mockCollection.AssertExpectations(t)
	mockSingleResult.AssertExpectations(t)
}

/*

func TestGetUserByID(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := MongoDBRepo{
		Client:   mockClient,
		Database: "test-db",
	}

	id := primitive.NewObjectID()
	expectedUser := &models.User{
		ID: id,
	}

	// Mock MongoDB response (user found)
	mockCollection.On("FindOne", mock.Anything, bson.M{"_id": id}).Return(nil).Run(func(args mock.Arguments) {
		result := args.Get(0).(*mongo.SingleResult)
		// Simulate decoding the user
		result.Decode(expectedUser)
	})

	user, err := repo.GetUserByID(id)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)

	// Test user not found
	mockCollection.On("FindOne", mock.Anything, bson.M{"_id": primitive.NewObjectID()}).Return(errors.New("user not found"))

	user, err = repo.GetUserByID(primitive.NewObjectID())
	assert.Error(t, err)
	assert.Nil(t, user)
}


func TestRegisterUser(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := MongoDBRepo{
		Client:   mockClient,
		Database: "test-db",
	}

	newUser := &models.User{
		Email: "newuser@example.com",
	}

	// Mock MongoDB response for InsertOne
	mockCollection.On("InsertOne", mock.Anything, newUser).Return(&mongo.InsertOneResult{}, nil)

	err := repo.RegisterUser(newUser)

	// Assert
	assert.NoError(t, err)

	// Test error scenario
	mockCollection.On("InsertOne", mock.Anything, newUser).Return(nil, errors.New("insert error"))

	err = repo.RegisterUser(newUser)
	assert.Error(t, err)
}


func TestUploadDocumentMetadata(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := MongoDBRepo{
		Client:   mockClient,
		Database: "test-db",
	}

	newDocument := &models.Document{
		Title:   "Test Document",
		UserID:  primitive.NewObjectID(),
		Subject: "Math",
		Grade:   "A",
	}

	// Mock MongoDB response for InsertOne
	mockCollection.On("InsertOne", mock.Anything, newDocument).Return(&mongo.InsertOneResult{}, nil)

	err := repo.UploadDocumentMetadata(newDocument)

	// Assert
	assert.NoError(t, err)

	// Test error scenario
	mockCollection.On("InsertOne", mock.Anything, newDocument).Return(nil, errors.New("insert error"))

	err = repo.UploadDocumentMetadata(newDocument)
	assert.Error(t, err)
}


func TestFindDocuments(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := MongoDBRepo{
		Client:   mockClient,
		Database: "test-db",
	}

	title := "Math"
	subject := "Math"
	grade := "A"
	correctRole := true

	expectedDocuments := []models.Document{
		{Title: "Document 1", Subject: subject, Grade: grade},
		{Title: "Document 2", Subject: subject, Grade: grade},
	}

	// Mock MongoDB response for Find
	mockCollection.On("Find", mock.Anything, bson.M{
		"title":   title,
		"subject": subject,
		"grade":   grade,
	}).Return(expectedDocuments, nil)

	documents, err := repo.FindDocuments(title, subject, grade, correctRole)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, documents, 2)
	assert.Equal(t, expectedDocuments[0].Title, documents[0].Title)

	// Test error scenario
	mockCollection.On("Find", mock.Anything, bson.M{
		"title":   title,
		"subject": subject,
		"grade":   grade,
	}).Return(nil, errors.New("find error"))

	documents, err = repo.FindDocuments(title, subject, grade, correctRole)
	assert.Error(t, err)
	assert.Nil(t, documents)
}


func TestGetFAQs(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := MongoDBRepo{
		Client:   mockClient,
		Database: "test-db",
	}

	expectedFAQs := []models.FAQs{
		{Question: "What is Share2Teach?", Answer: "It's a knowledge-sharing platform."},
		{Question: "How to use Share2Teach?", Answer: "Create an account and start sharing."},
	}

	// Mock MongoDB response for Find
	mockCollection.On("Find", mock.Anything, bson.M{}).Return(expectedFAQs, nil)

	faqs, err := repo.GetFAQs()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, faqs, 2)
	assert.Equal(t, expectedFAQs[0].Question, faqs[0].Question)

	// Test error scenario
	mockCollection.On("Find", mock.Anything, bson.M{}).Return(nil, errors.New("find error"))

	faqs, err = repo.GetFAQs()
	assert.Error(t, err)
	assert.Nil(t, faqs)
}


func TestGetDocumentByID(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := MongoDBRepo{
		Client:   mockClient,
		Database: "test-db",
	}

	id := primitive.NewObjectID()
	expectedDocument := &models.Document{
		ID:    id,
		Title: "Test Document",
	}

	// Mock MongoDB response for FindOne
	mockCollection.On("FindOne", mock.Anything, bson.M{"_id": id}).Return(nil).Run(func(args mock.Arguments) {
		result := args.Get(0).(*mongo.SingleResult)
		// Simulate decoding the document
		result.Decode(expectedDocument)
	})

	document, err := repo.GetDocumentByID(id)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedDocument.ID, document.ID)

	// Test document not found
	mockCollection.On("FindOne", mock.Anything, bson.M{"_id": primitive.NewObjectID()}).Return(errors.New("document not found"))

	document, err = repo.GetDocumentByID(primitive.NewObjectID())
	assert.Error(t, err)
	assert.Nil(t, document)
}


func TestGetDocumentRating(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := MongoDBRepo{
		Client:   mockClient,
		Database: "test-db",
	}

	id := primitive.NewObjectID()
	expectedRating := &models.Rating{
		DocumentID: id,
		Rating:     4.5,
	}

	// Mock MongoDB response for FindOne
	mockCollection.On("FindOne", mock.Anything, bson.M{"documentid": id}).Return(nil).Run(func(args mock.Arguments) {
		result := args.Get(0).(*mongo.SingleResult)
		// Simulate decoding the rating
		result.Decode(expectedRating)
	})

	rating, err := repo.GetDocumentRating(id)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedRating.Rating, rating.Rating)

	// Test rating not found
	mockCollection.On("FindOne", mock.Anything, bson.M{"documentid": primitive.NewObjectID()}).Return(errors.New("rating not found"))

	rating, err = repo.GetDocumentRating(primitive.NewObjectID())
	assert.Error(t, err)
	assert.Nil(t, rating)
}


func TestSetDocumentRating(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := MongoDBRepo{
		Client:   mockClient,
		Database: "test-db",
	}

	id := primitive.NewObjectID()
	newRating := &models.Rating{
		DocumentID: id,
		Rating:     4.5,
	}

	// Mock MongoDB response for UpdateOne
	mockCollection.On("UpdateOne", mock.Anything, bson.M{"documentid": id}, bson.M{"$set": bson.M{"rating": newRating.Rating}}).Return(&mongo.UpdateResult{}, nil)

	err := repo.SetDocumentRating(id, newRating)

	// Assert
	assert.NoError(t, err)

	// Test error scenario
	mockCollection.On("UpdateOne", mock.Anything, bson.M{"documentid": id}, bson.M{"$set": bson.M{"rating": newRating.Rating}}).Return(nil, errors.New("update error"))

	err = repo.SetDocumentRating(id, newRating)
	assert.Error(t, err)
}


func TestCreateDocumentRating(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := MongoDBRepo{
		Client:   mockClient,
		Database: "test-db",
	}

	newRating := &models.Rating{
		DocumentID: primitive.NewObjectID(),
		Rating:     5.0,
	}

	// Mock MongoDB response for InsertOne
	mockCollection.On("InsertOne", mock.Anything, newRating).Return(&mongo.InsertOneResult{}, nil)

	err := repo.CreateDocumentRating(newRating)

	// Assert
	assert.NoError(t, err)

	// Test error scenario
	mockCollection.On("InsertOne", mock.Anything, newRating).Return(nil, errors.New("insert error"))

	err = repo.CreateDocumentRating(newRating)
	assert.Error(t, err)
}


func TestStoreResetToken(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := MongoDBRepo{
		Client:   mockClient,
		Database: "test-db",
	}

	resetToken := &models.PasswordReset{
		UserID:    primitive.NewObjectID(),
		ResetKey:  "reset-token-123",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	// Mock MongoDB response for InsertOne
	mockCollection.On("InsertOne", mock.Anything, resetToken).Return(&mongo.InsertOneResult{}, nil)

	err := repo.StoreResetToken(resetToken)

	// Assert
	assert.NoError(t, err)

	// Test error scenario
	mockCollection.On("InsertOne", mock.Anything, resetToken).Return(nil, errors.New("insert error"))

	err = repo.StoreResetToken(resetToken)
	assert.Error(t, err)
}


func TestVerifyResetToken(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := MongoDBRepo{
		Client:   mockClient,
		Database: "test-db",
	}

	id := primitive.NewObjectID()
	resetToken := "reset-token-123"
	passwordReset := &models.PasswordReset{
		UserID:   id,
		ResetKey: resetToken,
	}

	// Mock MongoDB response for FindOne
	mockCollection.On("FindOne", mock.Anything, bson.M{"userid": id, "resetkey": resetToken}).Return(nil).Run(func(args mock.Arguments) {
		result := args.Get(0).(*mongo.SingleResult)
		// Simulate decoding the PasswordReset document
		result.Decode(passwordReset)
	})

	valid, err := repo.VerifyResetToken(id, resetToken)

	// Assert
	assert.NoError(t, err)
	assert.True(t, valid)

	// Test invalid token
	mockCollection.On("FindOne", mock.Anything, bson.M{"userid": id, "resetkey": "invalid-token"}).Return(nil, errors.New("token not found"))

	valid, err = repo.VerifyResetToken(id, "invalid-token")
	assert.Error(t, err)
	assert.False(t, valid)
}

func TestChangeUserPassword(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := MongoDBRepo{
		Client:   mockClient,
		Database: "test-db",
	}

	id := primitive.NewObjectID()
	newPassword := "newSecurePassword123"

	// Mock MongoDB response for UpdateOne
	mockCollection.On("UpdateOne", mock.Anything, bson.M{"_id": id}, bson.M{"$set": bson.M{"password": newPassword}}).Return(&mongo.UpdateResult{}, nil)

	err := repo.ChangeUserPassword(id, newPassword)

	// Assert
	assert.NoError(t, err)

	// Test error scenario
	mockCollection.On("UpdateOne", mock.Anything, bson.M{"_id": id}, bson.M{"$set": bson.M{"password": newPassword}}).Return(nil, errors.New("update error"))

	err = repo.ChangeUserPassword(id, newPassword)
	assert.Error(t, err)
}


func TestUpdateDocumentsByID(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := MongoDBRepo{
		Client:   mockClient,
		Database: "test-db",
	}

	documentID := primitive.NewObjectID()
	updateData := bson.M{"title": "Updated Title"}

	// Mock MongoDB response for UpdateOne
	mockCollection.On("UpdateOne", mock.Anything, bson.M{"_id": documentID}, bson.M{"$set": updateData}).Return(&mongo.UpdateResult{}, nil)

	err := repo.UpdateDocumentsByID(documentID, updateData)

	// Assert
	assert.NoError(t, err)

	// Test error scenario
	mockCollection.On("UpdateOne", mock.Anything, bson.M{"_id": documentID}, bson.M{"$set": updateData}).Return(nil, errors.New("update error"))

	err = repo.UpdateDocumentsByID(documentID, updateData)
	assert.Error(t, err)
}


func TestInsertModerationData(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := MongoDBRepo{
		Client:   mockClient,
		Database: "test-db",
	}

	userID := primitive.NewObjectID()
	documentID := primitive.NewObjectID()
	approvalStatus := "approved"
	comments := "Looks good"

	moderationData := bson.M{
		"userid":         userID,
		"documentid":     documentID,
		"approvalstatus": approvalStatus,
		"comments":       comments,
	}

	// Mock MongoDB response for InsertOne
	mockCollection.On("InsertOne", mock.Anything, moderationData).Return(&mongo.InsertOneResult{}, nil)

	err := repo.InsertModerationData(userID, documentID, approvalStatus, comments)

	// Assert
	assert.NoError(t, err)

	// Test error scenario
	mockCollection.On("InsertOne", mock.Anything, moderationData).Return(nil, errors.New("insert error"))

	err = repo.InsertModerationData(userID, documentID, approvalStatus, comments)
	assert.Error(t, err)
}


func TestInsertReport(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := MongoDBRepo{
		Client:   mockClient,
		Database: "test-db",
	}

	report := bson.M{
		"documentid": primitive.NewObjectID(),
		"userid":     primitive.NewObjectID(),
		"reason":     "Inappropriate content",
	}

	// Mock MongoDB response for InsertOne
	mockCollection.On("InsertOne", mock.Anything, report).Return(&mongo.InsertOneResult{}, nil)

	insertResult, err := repo.InsertReport(report)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, insertResult)

	// Test error scenario
	mockCollection.On("InsertOne", mock.Anything, report).Return(nil, errors.New("insert error"))

	insertResult, err = repo.InsertReport(report)
	assert.Error(t, err)
	assert.Nil(t, insertResult)
}
*/
