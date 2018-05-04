package models

import (
	"gopkg.in/mgo.v2/bson"
)

//Song contains metadata on a single song but without image and other large or
//binary data
type Song struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	Title    string        `json:"title" bson:"title"`
	Artist   string        `json:"artist" bson:"artist"`
	IsDuet   bool          `json:"isDuet" bson:"is_duet"`
	Language string        `json:"language" bson:"language"`
	BPM      float64       `json:"bpm" bson:"bpm"`
	Genre    string        `json:"genre" bson:"genre"`
	Source   string        `json:"source" bson:"source"`
	Year     int           `json:"year" bson:"year"`
}
