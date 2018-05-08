package manager

import (
	"time"

	"github.com/callummance/azunyan/db"
	"github.com/callummance/azunyan/models"
	"gopkg.in/mgo.v2/bson"
)

func AddRequest(m *KaraokeManager, singers []string, song string) {
	req := models.Request{
		ReqId:       bson.NewObjectId(),
		ReqTime:     time.Now(),
		Singers:     singers,
		Song:        bson.ObjectIdHex(song),
		PriorityMod: 0,
		Priority:    0,
		PlayedTime:  nil,
	}
	err := db.InsertRequest(m, req)
	if err != nil {
		m.Logger.Printf("Could not insert request for song %s due to error %s", song, err)
	}
	UpdatePriorities(m)
	UpdateListenersQueue(m)
}

func SetActive(m *KaraokeManager, newActiveState bool) error {
	newstate, err := db.GetEngineState(m, m.GetConfig().KaraokeConfig.SessionName)
	if err != nil {
		m.Logger.Printf("Failed to get session data due to error %q", err)
	}
	newstate.IsActive = newActiveState
	err = db.UpdateEngineState(m, *newstate)
	if err != nil {
		m.Logger.Printf("Failed to activate manager due to error %q", err)
		return err
	} else {
		m.SendBroadcast(BroadcastData{
			Name:    "active",
			Content: newActiveState,
		})
		return nil
	}
}

func SetReqActive(m *KaraokeManager, newActiveState bool) error {
	state, err := db.GetEngineState(m, m.Config.KaraokeConfig.SessionName)
	if err != nil {
		m.Logger.Printf("Failed to get session data due to error %q", err)
	}
	state.RequestsActive = newActiveState
	err = db.UpdateEngineState(m, *state)
	if err != nil {
		m.Logger.Printf("Failed to activate manager due to error %q", err)
		return err
	} else {
		return nil
	}
}

func UpdateListenersQueue(m *KaraokeManager) {
	queued := db.GetQueued(m)
	var abbQueue []models.AbbreviatedQueueItem

	for _, item := range queued {
		song, err := db.GetSongByID(m, item.Song)
		if err != nil {
			m.Logger.Printf("Error whilst retrieving queued songs list: %q", err)
			continue
		}
		abbQueue = append(abbQueue, item.Abbreviate(*song))
	}

	//Send updates
	m.SendBroadcast(BroadcastData{
		Name:    "queue",
		Content: abbQueue,
	})
}

func UpdateListenersCur(m *KaraokeManager) {
	state, err := db.GetEngineState(m, m.Config.KaraokeConfig.SessionName)
	if err != nil {
		m.Logger.Printf("Failed to get session data due to error %q", err)
	}
	song, err := db.GetSongByID(m, state.NowPlaying.Song)
	if err != nil {
		m.Logger.Printf("Error whilst retrieving queued songs list: %q", err)
	}

	//Send updates
	m.SendBroadcast(BroadcastData{
		Name:    "cur",
		Content: state.NowPlaying.Abbreviate(*song),
	})
}
