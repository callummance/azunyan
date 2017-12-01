package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Request struct {
	ReqId		bson.ObjectId	`json:"_" bson:"_id"`
	ReqTime 	time.Time		`json:"time" bson:"time"`
	Singers 	[]string		`json:"singers" bson:"singers"`
	Song		bson.ObjectId	`json:"songid" bson:"songid"`
	PriorityMod int				`json:"prioritymod" bson:"prioritymod"`
	Priority	float64			`json:"priority" bson:"priority"`
	PlayedTime  *time.Time		`json:"playedtime" bson:"playedtime,omitempty"`
}

type AbbreviatedQueueItem struct {
	ReqId		string			`json:"reqid"`
	SongTitle	string			`json:"songtitle"`
	SongArtist	string			`json:"songartist"`
	Singers 	[]string		`json:"singers"`
}

func (r Request) Abbreviate(song Song) AbbreviatedQueueItem {
	return AbbreviatedQueueItem{
		ReqId:r.ReqId.Hex(),
		SongTitle:song.Title,
		SongArtist:song.Artist,
		Singers:r.Singers,
	}
}