package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Request contains details on a single request that someone made for a song.
type Request struct {
	ReqID       primitive.ObjectID `json:"_" bson:"_id"`
	ReqTime     time.Time          `json:"time" bson:"time"`
	Singer      string             `json:"singer" bson:"singer"`
	Song        primitive.ObjectID `json:"songid" bson:"songid"`
	PriorityMod int                `json:"prioritymod" bson:"prioritymod"`
	PlayedTime  *time.Time         `json:"playedtime" bson:"playedtime,omitempty"`
}

//QueueItem contains details on an enqueued song, aggregated from its
//requests
type QueueItem struct {
	QueueItemID  primitive.ObjectID   `json:"queueitemid" bson:"queueitemid"`
	RequestIDs   []primitive.ObjectID `json:"ids" bson:"reqids"`
	SongID       primitive.ObjectID   `json:"sid" bson:"songid"`
	SongTitle    string               `json:"title" bson:"title"`
	SongArtist   string               `json:"artist" bson:"artist"`
	Singers      []string             `json:"singers" bson:"singers"`
	RequestTimes []time.Time          `json:"times" bson:"times"`
}

type AbbreviatedQueueItem struct {
	ReqId      string `json:"reqid"`
	SongTitle  string `json:"songtitle"`
	SongArtist string `json:"songartist"`
	Singer     string `json:"singers"`
}

func (r Request) Abbreviate(song Song) AbbreviatedQueueItem {
	return AbbreviatedQueueItem{
		ReqId:      r.ReqID.Hex(),
		SongTitle:  song.Title,
		SongArtist: song.Artist,
		Singer:     r.Singer,
	}
}
