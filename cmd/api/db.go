package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
