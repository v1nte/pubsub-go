package handlers

import (
	"context"
	"log"
	"time"

	"github.com/v1nte/pubsub-go/database"
	"github.com/v1nte/pubsub-go/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SaveMessageToDB(author, topic, message string) {
	msg := models.Message{
		ID:       primitive.NewObjectID(),
		PostedAt: time.Now(),
		Author:   author,
		Topic:    topic,
		Message:  message,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := database.Messages.InsertOne(ctx, msg)
	if err != nil {
		log.Println("Failed to insert message into DB", err)
		return
	}

	log.Println("Message inserted", res.InsertedID)
}
