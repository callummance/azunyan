package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Song contains metadata on a single song but without image and other large or
//binary data
type Song struct {
	ID       primitive.ObjectID `json:"ID" bson:"_id"`
	Title    string             `json:"title" bson:"title"`
	Artist   string             `json:"artist" bson:"artist"`
	IsDuet   bool               `json:"isDuet" bson:"is_duet"`
	Language string             `json:"language" bson:"language"`
	BPM      float64            `json:"bpm" bson:"bpm"`
	Genre    string             `json:"genre" bson:"genre"`
	Source   string             `json:"source" bson:"source"`
	Year     int                `json:"year" bson:"year"`
}

//SongSearchData contains just the Title, Artist, Source and ID of each track,
//and is intended to only be used when searching the database
type SongSearchData struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	Title  string             `json:"title" bson:"title"`
	Artist string             `json:"artist" bson:"artist"`
	Source string             `json:"source" bson:"source"`
}
