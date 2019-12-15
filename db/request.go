package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/callummance/azunyan/models"
	"github.com/globalsign/mgo/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

func InsertRequest(env databaseConfig, request models.Request) error {
	//Check the song is valid
	song, err := GetSongByID(env, request.Song)
	if err != nil {
		return err
	} else if song == nil {
		//Song is not in database
		return fmt.Errorf("song %q does not exist in the song database", request.Song)
	} else {
		//Song exists \o/
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := getReqCollection(env).InsertOne(ctx, request)
		if err != nil {
			env.GetLog().Printf("Could not insert request %+v; encountered error %q", request, err)
			return fmt.Errorf("request for song %s could not be inserted due to error %s", request.Song, err)
		}
		return nil
	}
}

//ResetRequests removes all requests stored in the current azunyan session
func ResetRequests(env databaseConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := getReqCollection(env).DeleteMany(ctx, bson.M{})
	return err
}

//RemoveRequests removes all the requests by a given singer and returns the number of items deleted
func RemoveRequests(env databaseConfig, singer string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	inf, err := getReqCollection(env).DeleteMany(ctx, bson.M{
		"singer": singer,
	})

	if err != nil {
		env.GetLog().Printf("Failure when removing requests: %v", err)
		return -1, err
	}
	return inf.DeletedCount, nil
}

//CheckDupeRequest returns true iff the given request is a request for the same
//song by the same person as another *active* request (that is, one that has
//not already been played).
func CheckDupeRequest(env databaseConfig, request models.Request) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cnt, err := getReqCollection(env).CountDocuments(ctx, bson.M{
		"singer":     request.Singer,
		"song":       request.Song,
		"playedtime": bson.M{"$exists": false},
	})
	if err != nil {
		env.GetLog().Printf("Failure when checking if request is a duplicate: %v", err)
		return false, err
	}
	return cnt > 0, nil
}

//GetLiveRequestsForSong returns a list of requests for a given song out of
//those that have not yet been played
func GetLiveRequestsForSong(env databaseConfig, sid primitive.ObjectID) ([]models.Request, error) {
	var res []models.Request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := getReqCollection(env).Find(ctx, bson.M{
		"song":       sid,
		"playedtime": bson.M{"$exists": false},
	})
	defer cursor.Close(context.Background())
	if err != nil {
		env.GetLog().Printf("Failure whilst finding existing requests for a song: %v", err)
		return nil, err
	}
	res, err = resultsToRequestsArray(cursor)
	if err != nil {
		env.GetLog().Printf("Failure whilst converting results to an array: %v", err)
		return res, err
	}
	return res, nil
}

type songRequests struct {
	ID       primitive.ObjectID `bson:"_id"`
	Requests []models.Request   `bson:"requests"`
}

//GetLiveAggregatedSongRequests returns a map mapping song IDs to a list of live
//requests for that song.
func GetLiveAggregatedSongRequests(env databaseConfig) (map[primitive.ObjectID][]models.Request, error) {
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"playedtime": bson.M{"$exists": false},
			},
		},
		bson.M{
			"$sort": bson.M{"time": 1},
		},
		bson.M{
			"$group": bson.M{
				"_id":      "$songid",
				"requests": bson.M{"$push": "$$ROOT"},
			},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := getReqCollection(env).Aggregate(ctx, pipeline)
	if err != nil {
		env.GetLog().Printf("Failed to fetch song requests list due to error %q", err)
		return nil, err
	}
	mapRes := make(map[primitive.ObjectID][]models.Request)
	for cursor.Next(ctx) {
		var song songRequests
		cursor.Decode(&song)
		mapRes[song.ID] = song.Requests
	}
	return mapRes, nil
}

func GetPreviousRequestsBySong(env databaseConfig, songId primitive.ObjectID, submitted time.Time) ([]models.Request, error) {
	var res []models.Request
	col := getReqCollection(env)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := col.Find(ctx, bson.M{
		"songid": songId,
		"time":   bson.M{"$lt": submitted},
	})
	if err != nil {
		env.GetLog().Printf("Failure when getting previous requests: %q", err)
		return nil, err
	}
	res, err = resultsToRequestsArray(cursor)
	if err != nil {
		env.GetLog().Printf("Failure whilst converting results to an array: %v", err)
		return res, err
	}
	return res, nil
}

func GetPreviousRequestsBySinger(env databaseConfig, singer string, submitted time.Time) []models.Request {
	var res []models.Request
	col := getReqCollection(env)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := col.Find(ctx, bson.M{
		"singers": singer,
		"time":    bson.M{"$lt": submitted},
	})
	if err != nil {
		env.GetLog().Printf("Failure when getting previous requests by singer: %q", err)
	}
	res, err = resultsToRequestsArray(cursor)
	if err != nil {
		env.GetLog().Printf("Failure whilst converting results to an array: %v", err)
		return res
	}
	return res
}

func GetLiveRequests(env databaseConfig) []models.Request {
	var res []models.Request
	col := getReqCollection(env)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"time": 1})
	cursor, err := col.Find(ctx, bson.M{
		"playedtime": bson.M{"$exists": false},
	}, findOptions)
	if err != nil {
		env.GetLog().Printf("Failure when getting queued songs: %q", err)
	}
	res, err = resultsToRequestsArray(cursor)
	if err != nil {
		env.GetLog().Printf("Failure whilst converting results to an array: %v", err)
		return res
	}
	return res
}

func SetRequestPlayed(env databaseConfig, reqId primitive.ObjectID, playedTime time.Time) error {
	col := getReqCollection(env)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := col.UpdateOne(ctx, bson.M{"_id": reqId}, bson.M{"$set": bson.M{"playedtime": playedTime}})
	if err != nil {
		env.GetLog().Printf("Couldn't update playedtime due to error %q", err)
	}
	return err
}

func UpdateReqPrio(env databaseConfig, reqId primitive.ObjectID, newPrio float64) error {
	col := getReqCollection(env)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := col.UpdateOne(ctx, bson.M{"_id": reqId}, bson.M{"$set": bson.M{"priority": newPrio}})
	if err != nil {
		env.GetLog().Printf("Couldn't update priority due to error %q", err)
	}
	return err
}

func getReqCollection(env databaseConfig) *mongo.Collection {
	return env.GetClient().Database(env.GetConfig().DbConfig.DatabaseName).Collection("request")
}

func resultsToRequestsArray(cursor *mongo.Cursor) ([]models.Request, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var res []models.Request
	for cursor.Next(ctx) {
		var elem models.Request
		err := cursor.Decode(&elem)
		if err != nil {
			return res, err
		}
		res = append(res, elem)
	}
	return res, nil
}
