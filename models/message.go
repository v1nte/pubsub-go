package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	PostedAt time.Time          `json:"postedAt" bson:"posted_at"`
	Author   string             `json:"author" bson:"author"`
	Topic    string             `json:"topic" bson:"topic"`
	Message  string             `json:"message" bson:"message"`
}
