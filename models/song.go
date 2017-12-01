package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Song struct {
	Id 		bson.ObjectId	`json:"id" bson:"_id"`
	Title	string			`json:"title" bson:"title"`
	Artist 	string 			`json:"artist" bson:"artist"`
}
