package db

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime/debug"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/callummance/azunyan/models"
	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo"
	// "github.com/globalsign/mgo"
	// "github.com/globalsign/mgo/bson"
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
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
		songObj := models.Song{ID: primitive.NewObjectID(), Title: song["title"], Artist: song["artist"]}
		col := getCollection(env)
		if checkSongExists(env, songObj) {
			continue
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err := col.InsertOne(ctx, songObj)
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	findOptions := options.Find()
	findOptions.SetProjection(bson.M{
		"creator":  1,
		"title":    1,
		"artist":   1,
		"source":   1,
		"language": 1,
		"bpm":      1,
		"is_duet":  1,
		"genre":    1,
		"year":     1})
	cursor, err := getCollection(env).Find(ctx, bson.M{}, findOptions)
	if err != nil {
		env.GetLog().Printf("Could not fetch songlist due to error %q", err)
	}
	songs, err = resultsToSongsArray(cursor)
	if err != nil {
		env.GetLog().Printf("Failure whilst converting results to an array: %v", err)
		return songs
	}
	return songs
}

func GetSongTAS(env databaseConfig) []models.SongSearchData {
	var songs []models.SongSearchData
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	findOptions := options.Find()
	findOptions.SetProjection(bson.M{
		"title":  1,
		"artist": 1,
		"source": 1})
	cursor, err := getCollection(env).Find(ctx, bson.M{}, findOptions)
	// env.GetLog().Printf("%+v\n", songs)
	if err != nil {
		env.GetLog().Printf("Could not fetch songlist due to error %q", err)
	}
	songs, err = resultsToSongSearchDataArray(cursor)
	if err != nil {
		env.GetLog().Printf("Failure whilst converting results to an array: %v", err)
		return songs
	}
	return songs
}

func GetSongsByIDs(env databaseConfig, songIDs []primitive.ObjectID) ([]models.Song, error) {
	findOptions := options.Find()
	findOptions.SetProjection(bson.M{
		"_id":      1,
		"creator":  1,
		"title":    1,
		"artist":   1,
		"language": 1,
		"bpm":      1,
		"is_duet":  1,
		"genre":    1,
		"year":     1})

	var cursor *mongo.Cursor
	var err error
	if len(songIDs) == 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		cursor, err = getCollection(env).Find(ctx, bson.M{}, findOptions)
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancel()
		cursor, err = getCollection(env).Find(ctx, bson.M{"_id": bson.M{"$in": songIDs}}, findOptions)
	}
	if err != nil {
		env.GetLog().Printf("Failed to check database for song ids due to reason '%s'", err)
	}
	songs, err := resultsToSongsArray(cursor)
	if err != nil {
		env.GetLog().Printf("Failure whilst converting results to an array: %v", err)
		return nil, err
	}
	return songs, err
}

func GetSongByID(env databaseConfig, sid primitive.ObjectID) (*models.Song, error) {
	var res models.Song
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	findOptions := options.FindOne()
	findOptions.SetProjection(bson.M{
		"_id":      1,
		"creator":  1,
		"title":    1,
		"artist":   1,
		"language": 1,
		"bpm":      1,
		"is_duet":  1,
		"genre":    1,
		"year":     1})
	singleResult := getCollection(env).FindOne(ctx, bson.M{"_id": sid}, findOptions)
	err := singleResult.Err()
	if err != nil && err.Error() == "not found" {
		env.GetLog().Printf("Couldn't find song with id %q", sid)
		return nil, nil
	} else if err != nil {
		debug.PrintStack()
		fmt.Print("\n\n\n")
		env.GetLog().Printf("Failed to check database for song id %q due to reason '%s'", sid, err)
		return nil, fmt.Errorf("database failure occurred: %q", err)
	} else {
		singleResult.Decode(&res)
		return &res, nil
	}
}

func GetSongCoverByID(env databaseConfig, sid primitive.ObjectID) ([]byte, error) {
	var res struct {
		ID    primitive.ObjectID `json:"id" bson:"_id"`
		Cover []byte             `json:"cover" bson:"cover"`
	}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	findOptions := options.FindOne()
	findOptions.SetProjection(bson.M{
		"_id":   0,
		"cover": 1})
	singleResult := getCollection(env).FindOne(ctx, bson.M{"_id": sid}, findOptions)
	err := singleResult.Err()
	if err != nil && err.Error() == "not found" {
		return nil, nil
	} else if err != nil {
		env.GetLog().Printf("Failed to check database for song id %q due to reason '%s'", sid, err)
		return nil, fmt.Errorf("database failure occurred: %q", err)
	} else {
		singleResult.Decode(&res)
		return res.Cover, nil
	}
}

func getCollection(env databaseConfig) *mongo.Collection {
	return env.GetClient().Database(env.GetConfig().DbConfig.DatabaseName).Collection(env.GetConfig().DbConfig.DatabaseCollectionName)
}

func checkSongExists(env databaseConfig, song models.Song) bool {
	col := getCollection(env)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cnt, err := col.CountDocuments(ctx, bson.M{"title": song.Title, "artist": song.Artist})
	if err != nil {
		env.GetLog().Printf("could not check if song %s exists, error was %q", song.Title, err)
		return false
	} else {
		return cnt > 0
	}
}

func resultsToSongsArray(cursor *mongo.Cursor) ([]models.Song, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var res []models.Song
	for cursor.Next(ctx) {
		var elem models.Song
		err := cursor.Decode(&elem)
		if err != nil {
			return res, err
		}
		res = append(res, elem)
	}
	return res, nil
}

func resultsToSongSearchDataArray(cursor *mongo.Cursor) ([]models.SongSearchData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var res []models.SongSearchData
	for cursor.Next(ctx) {
		var elem models.SongSearchData
		err := cursor.Decode(&elem)
		if err != nil {
			return res, err
		}
		res = append(res, elem)
	}
	return res, nil
}
