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

func (m *MongoDBRepo) Connection() *mongo.Client {
	return m.Client
}

// func (m *MongoDBRepo) AllMovies() ([]*models.Movie, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
// 	defer cancel()

// 	// get collection
// 	collection := m.Client.Database(m.Database).Collection("movies")

// 	cursor, err := collection.Find(ctx, bson.M{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	var movies []*models.Movie
// 	for cursor.Next(ctx) {
// 		var movie models.Movie
// 		err := cursor.Decode(&movie)
// 		if err != nil {
// 			return nil, err
// 		}
// 		movies = append(movies, &movie)
// 	}

// 	if err := cursor.Err(); err != nil {
// 		return nil, err
// 	}

// 	return movies, nil
// }

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
