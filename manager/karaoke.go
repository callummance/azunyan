package manager

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"fmt"
	"github.com/callummance/azunyan/webserver/stream"
	"github.com/callummance/azunyan/models"
	"github.com/callummance/azunyan/db"
	"github.com/callummance/azunyan/state"
)

func AddRequest (e *state.Env, singers []string, song string) {
	req := models.Request{
		ReqId:       bson.NewObjectId(),
		ReqTime:     time.Now(),
		Singers:     singers,
		Song:        bson.ObjectIdHex(song),
		PriorityMod: 0,
		Priority:	 0,
		PlayedTime:  nil,
	}
	err := db.InsertRequest(e, req)
	if err != nil {
		fmt.Errorf("Could not insert request for song %s due to error %s", song, err)
	}
	UpdatePriorities(*e)
	UpdateListenersQueue(e)
}

func SetActive(e *state.Env, newActiveState bool) error {
	newstate, err := db.GetEngineState(e, e.Config.KaraokeConfig.SessionName)
	if err != nil {
		e.Logger.Printf("Failed to get session data due to error %q", err)
	}
	newstate.IsActive = newActiveState
	err = db.UpdateEngineState(e, *newstate)
	if err != nil {
		e.Logger.Printf("Failed to activate manager due to error %q", err)
		return err
	} else {
		stream.SendBroadcast(stream.BroadcastData{
			Name: "active",
			Content: newActiveState,
		})
		return nil
	}
}


func SetReqActive(e *state.Env, newActiveState bool) error {
	state, err := db.GetEngineState(e, e.Config.KaraokeConfig.SessionName)
	if err != nil {
		e.Logger.Printf("Failed to get session data due to error %q", err)
	}
	state.RequestsActive = newActiveState
	err = db.UpdateEngineState(e, *state)
	if err != nil {
		e.Logger.Printf("Failed to activate manager due to error %q", err)
		return err
	} else {
		return nil
	}
}

func UpdateListenersQueue(env *state.Env) {
	queued := db.GetQueued(env)
	var abbQueue []models.AbbreviatedQueueItem

	for _, item := range queued {
		song, err := db.GetSongById(env, item.Song)
		if err != nil {
			env.Logger.Printf("Error whilst retrieving queued songs list: %q", err)
			continue
		}
		abbQueue = append(abbQueue, item.Abbreviate(*song))
	}

	//Send updates
	stream.SendBroadcast(stream.BroadcastData{
		Name: "queue",
		Content: abbQueue,
	})
}

func UpdateListenersCur(env *state.Env) {
	state, err := db.GetEngineState(env, env.Config.KaraokeConfig.SessionName)
	if err != nil {
		env.Logger.Printf("Failed to get session data due to error %q", err)
	}
	song, err := db.GetSongById(env, state.NowPlaying.Song)
	if err != nil {
		env.Logger.Printf("Error whilst retrieving queued songs list: %q", err)
	}

	//Send updates
	stream.SendBroadcast(stream.BroadcastData{
		Name: "cur",
		Content: state.NowPlaying.Abbreviate(*song),
	})
}
