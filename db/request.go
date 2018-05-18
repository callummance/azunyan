package db

import (
	"fmt"
	"time"

	"github.com/callummance/azunyan/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
		err := getReqCollection(env).Insert(request)
		if err != nil {
			env.GetLog().Printf("Could not insert request %+v; encountered error %q", request, err)
			return fmt.Errorf("request for song %s could not be inserted due to error %s", request.Song, err)
		}
		return nil
	}
}

//CheckDupeRequest returns true iff the given request is a request for the same
//song by the same person as another *active* request (that is, one that has
//not already been played).
func CheckDupeRequest(env databaseConfig, request models.Request) (bool, error) {
	q := getReqCollection(env).Find(bson.M{
		"singer":     request.Singer,
		"song":       request.Song,
		"playedtime": bson.M{"$exists": false},
	})
	cnt, err := q.Count()
	if err != nil {
		env.GetLog().Printf("Failure when checking if request is a duplicate: %v", err)
		return false, err
	}
	return cnt > 0, nil
}

//GetLiveRequestsForSong returns a list of requests for a given song out of
//those that have not yet been played
func GetLiveRequestsForSong(env databaseConfig, sid bson.ObjectId) ([]models.Request, error) {
	var res []models.Request
	err := getReqCollection(env).Find(bson.M{
		"song":       sid,
		"playedtime": bson.M{"$exists": false},
	}).All(res)
	if err != nil {
		env.GetLog().Printf("Failure whilst finding existing requests for a song: %v", err)
		return nil, err
	}
	return res, nil
}

type songRequests struct {
	ID       bson.ObjectId    `bson:"_ID"`
	Requests []models.Request `bson:"requests"`
}

//GetLiveAggregatedSongRequests returns a map mapping song IDs to a list of live
//requests for that song.
func GetLiveAggregatedSongRequests(env databaseConfig) (map[bson.ObjectId][]models.Request, error) {
	var res []songRequests
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

	err := getReqCollection(env).Pipe(pipeline).All(&res)
	if err != nil {
		env.GetLog().Printf("Failed to fetch song requests list due to error %q", err)
		return nil, err
	}

	mapRes := make(map[bson.ObjectId][]models.Request)
	for _, song := range res {
		mapRes[song.ID] = song.Requests
	}
	return mapRes, nil
}

func GetPreviousRequestsBySong(env databaseConfig, songId bson.ObjectId, submitted time.Time) []models.Request {
	var res []models.Request
	col := getReqCollection(env)
	q := col.Find(bson.M{
		"songid": songId,
		"time":   bson.M{"$lt": submitted},
	})
	err := q.All(&res)
	if err != nil {
		env.GetLog().Printf("Failure when getting previous requests: %q", err)
	}
	return res
}

func GetPreviousRequestsBySinger(env databaseConfig, singer string, submitted time.Time) []models.Request {
	var res []models.Request
	col := getReqCollection(env)
	q := col.Find(bson.M{
		"singers": singer,
		"time":    bson.M{"$lt": submitted},
	})
	err := q.All(&res)
	if err != nil {
		env.GetLog().Printf("Failure when getting previous requests: %q", err)
	}
	return res
}

func GetQueued(env databaseConfig) []models.Request {
	var res []models.Request
	col := getReqCollection(env)
	q := col.Find(bson.M{
		"playedtime": bson.M{"$exists": false},
	}).Sort("-priority")
	err := q.All(&res)
	if err != nil {
		env.GetLog().Printf("Failure when getting queued songs: %q", err)
	}
	return res
}

func SetRequestPlayed(env databaseConfig, reqId bson.ObjectId, playedTime time.Time) error {
	col := getReqCollection(env)
	err := col.UpdateId(reqId, bson.M{"$set": bson.M{"playedtime": playedTime}})
	if err != nil {
		env.GetLog().Printf("Couldn't update priority due to error %q", err)
	}
	return err
}

func UpdateReqPrio(env databaseConfig, reqId bson.ObjectId, newPrio float64) error {
	col := getReqCollection(env)
	err := col.UpdateId(reqId, bson.M{"$set": bson.M{"priority": newPrio}})
	if err != nil {
		env.GetLog().Printf("Couldn't update priority due to error %q", err)
	}
	return err
}

func getReqCollection(env databaseConfig) *mgo.Collection {
	return env.GetSession().DB(env.GetConfig().DbConfig.DatabaseName).C("request")
}
