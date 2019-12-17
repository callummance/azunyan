package db

import (
	"context"
	"time"

	"github.com/callummance/azunyan/models"
	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AddSongToUpcomingQueue adds a song to the upcomingQueue collection to mark that a song
// has fulfilled the required number of singers and has entered upcoming queue.
func AddSongToUpcomingQueue(env databaseConfig, requestID primitive.ObjectID, queueItem models.QueueItem) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	upcomingQueueItem := models.UpcomingQueueItem{
		RequestID: requestID,
		QueueItem: queueItem,
	}

	_, err := getCompletedReqCollection(env).InsertOne(ctx, upcomingQueueItem)
	if err != nil {
		env.GetLog().Printf("Could not insert upcoming queue item; encountered error %q", err)
	}
}

// RemoveSongFromUpcomingQueue removes a song from the upcomingQueue collection. This should be called
// when the song has been played.
func RemoveSongFromUpcomingQueue(env databaseConfig, queueItemID primitive.ObjectID) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := getCompletedReqCollection(env).DeleteMany(ctx, bson.M{"queueitem.queueitemid": queueItemID})
	if err != nil {
		env.GetLog().Printf("Could not remove from upcoming queue item; encountered error %q", err)
	}
	if res.DeletedCount == 0 {
		env.GetLog().Printf("Could not remove queue id %s from upcoming queue item", queueItemID.Hex())
	}
}

// GetUpcomingSongs retrieves the documents in upcomingqueue
func GetUpcomingSongs(env databaseConfig) map[primitive.ObjectID]models.UpcomingQueueItem {
	col := getCompletedReqCollection(env)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mapRes := make(map[primitive.ObjectID]models.UpcomingQueueItem)

	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		env.GetLog().Printf("Failure whilst retrieving upcoming queue; error %q", err)
		return mapRes
	}
	for cursor.Next(ctx) {
		var upcomingQueueItem models.UpcomingQueueItem
		cursor.Decode(&upcomingQueueItem)
		mapRes[upcomingQueueItem.RequestID] = upcomingQueueItem
	}
	return mapRes
}

// ClearUpcomingSongs deletes all documents in upcomingqueue collection
func ClearUpcomingSongs(env databaseConfig) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := getCompletedReqCollection(env).DeleteMany(ctx, bson.M{})
	if err != nil {
		env.GetLog().Printf("Failure whilst trying to clear upcoming queue; error %q", err)
	}
}

func getCompletedReqCollection(env databaseConfig) *mongo.Collection {
	return env.GetClient().Database(env.GetConfig().DbConfig.DatabaseName).Collection("upcomingqueue")
}
