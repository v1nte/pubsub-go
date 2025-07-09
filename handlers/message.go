package handlers

import (
	"context"
	"time"

	"github.com/v1nte/pubsub-go/database"
	"github.com/v1nte/pubsub-go/logger"
	"github.com/v1nte/pubsub-go/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
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
		logger.Log.Error("Failed to insert message into DB", zap.Error(err))
		return
	}

	logger.Log.Info("Message inserted", zap.Any("insertedID", res.InsertedID))
}
