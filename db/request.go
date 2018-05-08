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
