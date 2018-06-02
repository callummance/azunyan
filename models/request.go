package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//Request contains details on a single request that someone made for a song.
type Request struct {
	ReqID       bson.ObjectId `json:"_" bson:"_id"`
	ReqTime     time.Time     `json:"time" bson:"time"`
	Singer      string        `json:"singer" bson:"singer"`
	Song        bson.ObjectId `json:"songid" bson:"songid"`
	PriorityMod int           `json:"prioritymod" bson:"prioritymod"`
	PlayedTime  *time.Time    `json:"playedtime" bson:"playedtime,omitempty"`
}

//QueueItem contains details on an enqueued song, aggregated from its
//requests
type QueueItem struct {
	RequestIDs   []bson.ObjectId `json:"ids"`
	SongID       bson.ObjectId   `json:"sid"`
	SongTitle    string          `json:"title"`
	SongArtist   string          `json:"artist"`
	Singers      []string        `json:"singers"`
	RequestTimes []time.Time     `json:"times"`
}

//CompileQueueItems takes a list of requests for a song, as well as the song's
//struct itself and the number of singers which can be contained in a single
//queued item and returns a slice of complete queue items along with an optional
//incomplete queue item.
func CompileQueueItems(reqs []Request, song *Song, maxSingers int) ([]QueueItem, *QueueItem) {
	if len(reqs) == 0 {
		return nil, nil
	}
	noRequests := len(reqs)
	incompleteQueueLen := noRequests % maxSingers
	fullQueueItemsCount := noRequests / maxSingers
	completeRes := make([]QueueItem, fullQueueItemsCount)

	for i := 0; i < fullQueueItemsCount; i++ {
		item := QueueItem{
			SongID:       song.ID,
			SongTitle:    song.Title,
			SongArtist:   song.Artist,
			RequestIDs:   []bson.ObjectId{},
			Singers:      []string{},
			RequestTimes: []time.Time{},
		}
		for j := 0; j < maxSingers; j++ {
			offset := i*maxSingers + j
			item.RequestIDs = append(item.RequestIDs, reqs[offset].ReqID)
			item.Singers = append(item.Singers, reqs[offset].Singer)
			item.RequestTimes = append(item.RequestTimes, reqs[offset].ReqTime)
		}
		completeRes[i] = item
	}

	if incompleteQueueLen != 0 {
		item := QueueItem{
			SongID:       song.ID,
			SongTitle:    song.Title,
			SongArtist:   song.Artist,
			RequestIDs:   []bson.ObjectId{},
			Singers:      []string{},
			RequestTimes: []time.Time{},
		}
		for j := 0; j < incompleteQueueLen; j++ {
			offset := fullQueueItemsCount*maxSingers + j
			item.RequestIDs = append(item.RequestIDs, reqs[offset].ReqID)
			item.Singers = append(item.Singers, reqs[offset].Singer)
			item.RequestTimes = append(item.RequestTimes, reqs[offset].ReqTime)
		}
		return completeRes, &item
	}
	return completeRes, nil
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
