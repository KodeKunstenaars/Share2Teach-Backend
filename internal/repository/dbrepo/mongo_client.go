package dbrepo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient interface includes only the Database method used by the repo.
type MongoClient interface {
	Database(name string, opts ...*options.DatabaseOptions) *mongo.Database
}

// Ensure that *mongo.Client satisfies the MongoClient interface
var _ MongoClient = (*mongo.Client)(nil)
