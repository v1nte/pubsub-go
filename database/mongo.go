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
	LogsDB   *mongo.Collection
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
	LogsDB = client.Database("logs").Collection("appLogs")

	if err = client.Database("db").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		return err
	}

	if err = client.Database("logs").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		return err
	}

	return nil
}

func Close() error {
	return client.Disconnect(context.Background())
}
