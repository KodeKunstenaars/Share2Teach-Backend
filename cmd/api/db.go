// Everything here works

package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// func openMongoDB(dsn string) (*mongo.Client, error) {
// 	client, err := mongo.NewClient(options.Client().ApplyURI(dsn))
// 	if err != nil {
// 		return nil, err
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	err = client.Connect(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Ping the database to verify the connection
// 	err = client.Ping(ctx, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return client, nil
// }

func (app *application) connectToMongoDB() (*mongo.Client, error) {
	client, err := openMongo(app.DSN)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Mongo!")
	return client, nil
}

func openMongo(dsn string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(dsn)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
