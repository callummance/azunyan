package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime/debug"

	"github.com/callummance/azunyan/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

	for _, song := range res {
		songObj := models.Song{ID: bson.NewObjectId(), Title: song["title"], Artist: song["artist"]}
		col := getCollection(env)
		if checkSongExists(env, songObj) {
			continue
		} else {
			err := col.Insert(songObj)
			if err != nil {
				env.GetLog().Printf("Could not insert song %+v, encountered error %q", songObj, err)
			} else {
				count++
			}
		}
	}

	env.GetLog().Printf("Imported %d songs", count)
}

func GetSongs(env databaseConfig) []models.Song {
	var songs []models.Song
	err := getCollection(env).Find(bson.M{}).Select(bson.M{
		"creator":  1,
		"title":    1,
		"artist":   1,
		"source":   1,
		"language": 1,
		"bpm":      1,
		"is_duet":  1,
		"genre":    1,
		"year":     1}).All(&songs)
	if err != nil {
		env.GetLog().Printf("Could not fetch songlist due to error %q", err)
	}
	return songs
}

func GetSongTAS(env databaseConfig) []models.SongSearchData {
	var songs []models.SongSearchData
	err := getCollection(env).Find(bson.M{}).Select(bson.M{
		"title":  1,
		"artist": 1,
		"source": 1}).All(&songs)
	if err != nil {
		env.GetLog().Printf("Could not fetch songlist due to error %q", err)
	}
	return songs
}

func GetSongByID(env databaseConfig, sid bson.ObjectId) (*models.Song, error) {
	var res models.Song
	err := getCollection(env).FindId(sid).Select(bson.M{
		"creator":  1,
		"title":    1,
		"artist":   1,
		"language": 1,
		"bpm":      1,
		"is_duet":  1,
		"genre":    1,
		"year":     1}).One(&res)
	if err != nil && err.Error() == "not found" {
		return nil, nil
	} else if err != nil {
		debug.PrintStack()
		fmt.Print("\n\n\n")
		env.GetLog().Printf("Failed to check database for song id %q due to reason '%s'", sid, err)
		return nil, fmt.Errorf("database failure occurred: %q", err)
	} else {
		return &res, nil
	}
}

func GetSongCoverByID(env databaseConfig, sid bson.ObjectId) ([]byte, error) {
	var res struct {
		Cover bson.Binary `bson:"cover"`
	}
	err := getCollection(env).FindId(sid).Select(bson.M{
		"_id":   0,
		"cover": 1}).One(&res)
	if err != nil && err.Error() == "not found" {
		return nil, nil
	} else if err != nil {
		env.GetLog().Printf("Failed to check database for song id %q due to reason '%s'", sid, err)
		return nil, fmt.Errorf("database failure occurred: %q", err)
	} else {
		return res.Cover.Data, nil
	}
}

func getCollection(env databaseConfig) *mgo.Collection {
	return env.GetSession().DB(env.GetConfig().DbConfig.DatabaseName).C("song")
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
