package database

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client   *mongo.Client
	Messages *mongo.Collection
)

const (
	defaultMongoURI = "mongodb://root:root@localhost:27017"
	mongoURIVarName = "MONGO_URI"
)

func getMongoURI() string {
	if uri := os.Getenv(mongoURIVarName); uri != "" {
		return uri
	}
	return defaultMongoURI
}

func Init() error {
	opts := options.Client().ApplyURI(getMongoURI())

	localClient, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return err
	}

	client = localClient

	Messages = client.Database("db").Collection("messages")

	err = client.Database("db").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err()
	return err
}

func Close() error {
	return client.Disconnect(context.Background())
}
