package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	//add import to call db connectToMongo
)

func FAQ() {
	client, err := db.connectToMongoDB("uri")
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	db := client.Database("")
	coll := db.Collection("")

	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}

	for cursor.Next(context.TODO()) {
		var results bson.M
		if err := cursor.Decode(&results); err != nil {
			panic(err)
		}
	}
}
