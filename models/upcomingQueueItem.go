package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpcomingQueueItem struct {
	RequestID primitive.ObjectID `json:"ID" bson:"_id"`
	QueueItem QueueItem          `json:"queueitem" bson:"queueitem"`
}
