package manager

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/callummance/azunyan/db"
	"github.com/callummance/azunyan/models"
	// "gopkg.in/mgo.v2/bson"
)

//AddRequest takes a singer and the song id as a string and adds it to the
//request queue, also updating any listeners.
func AddRequest(m *KaraokeManager, singer string, song string) error {
	objID, err := primitive.ObjectIDFromHex(song)
	if err != nil {
		m.Logger.Printf("Could not insert request for song %s due to error %s", song, err)
		err = FetchAndUpdateListenersQueue(m, 0)
		return err
	}
	req := models.Request{
		ReqID:       primitive.NewObjectID(),
		ReqTime:     time.Now(),
		Singer:      singer,
		Song:        objID,
		PriorityMod: 0,
		PlayedTime:  nil,
	}
	err = db.InsertRequest(m, req)
	if err != nil {
		m.Logger.Printf("Could not insert request for song %s due to error %s", song, err)
	}
	err = FetchAndUpdateListenersQueue(m, 0)
	return err
}

//RemoveSinger removes all requests by a singer with the given name
func RemoveSinger(m *KaraokeManager, singer string) error {
	_, err := db.RemoveRequests(m, singer)
	if err != nil {
		return err
	}
	err = FetchAndUpdateListenersQueue(m, 0)
	return err
}

//Reset removes ALL requests.
func Reset(m *KaraokeManager) error {
	err := db.ResetRequests(m)
	if err != nil {
		return err
	}
	err = FetchAndUpdateListenersQueue(m, 0)
	UpdateListenersCur(m, nil)
	return err
}

//SetActive sets the karaoke system active and supdates all listeners
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

//SetDupesAllowed modifies the internal flag allowing requests for duplicate
//songs.
func SetDupesAllowed(m *KaraokeManager, newDupePermissions bool) error {
	newstate, err := db.GetEngineState(m, m.GetConfig().KaraokeConfig.SessionName)
	if err != nil {
		m.Logger.Printf("Failed to get session data due to error %q", err)
	}
	newstate.AllowingDupes = newDupePermissions
	err = db.UpdateEngineState(m, *newstate)
	if err != nil {
		m.Logger.Printf("Failed to set duplicate request permissions due to error %q", err)
		return err
	} else {
		return nil
	}
}

//SetReqActive updates the internal flag allowing users to make requests.
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

//ChangeNumberOfSingers updates the number of singers required for each request.
func ChangeNumberOfSingers(m *KaraokeManager, noSingers int) error {
	state, err := db.GetEngineState(m, m.Config.KaraokeConfig.SessionName)
	if err != nil {
		m.Logger.Printf("Failed to get session data due to error %q", err)
	}
	state.NoSingers = noSingers
	err = db.UpdateEngineState(m, *state)
	if err != nil {
		m.Logger.Printf("Failed to update number of singers due to error %q", err)
		return err
	}
	UpdateNumberSingers(m, noSingers)
	FetchAndUpdateListenersQueue(m, 1)
	return nil
}
