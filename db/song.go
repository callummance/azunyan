package db

import (
	"io/ioutil"
	"github.com/callummance/azunyan/models"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"fmt"
)

func ImportJSONSongList(env databaseConfig, fileLoc string) {
	rawFile, err := ioutil.ReadFile(fileLoc)
	if err != nil {
		env.GetLog().Printf("Could not locate songlist file at %s; encountered error %s", fileLoc, err)
		return
	}

	var res [](map[string]string)
	json.Unmarshal(rawFile, &res)

	var count int

	for _, song := range res  {
		songObj := models.Song{Id: bson.NewObjectId(), Title: song["title"], Artist: song["artist"]}
		col := getCollection(env)
		if checkSongExists(env, songObj) {
			continue
		} else {
			err := col.Insert(songObj)
			if err != nil {
				env.GetLog().Printf("Could not insert song %q, encountered error %q", songObj, err)
			} else {
				count += 1
			}
		}
	}

	env.GetLog().Printf("Imported %d songs", count)
}

func GetSongs(env databaseConfig) []models.Song {
	var songs []models.Song
	err := getCollection(env).Find(bson.M{}).All(&songs)
	if err != nil {
		env.GetLog().Printf("Could not fetch songlist due to error %q", err)
	}
	return songs
}

func GetSongById(env databaseConfig, sid bson.ObjectId) (*models.Song, error) {
	var res models.Song
	err := getCollection(env).FindId(sid).One(&res)
	if err != nil && err.Error() == "not found"  {
		return nil, nil
	} else if err != nil {
		env.GetLog().Printf("Failed to check database for song id %q due to reason '%s'", sid, err)
		return nil, fmt.Errorf("database failure occurred: %q", err)
	} else {
		return &res, nil
	}
}


func getCollection(env databaseConfig) *mgo.Collection {
	return env.GetSession().DB(env.GetDbConfig().DatabaseName).C("song")
}

func checkSongExists(env databaseConfig, song models.Song) bool {
	col := getCollection(env)
	q := col.Find(bson.M{"title": song.Title, "artist": song.Artist})
	cnt, err := q.Count()
	if err != nil {
		env.GetLog().Printf("could not check if song %s exists, error was %q", song.Title, err)
		return false
	} else {
		return cnt > 0
	}
}
