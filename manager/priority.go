package manager

import (
	"time"

	"github.com/callummance/azunyan/db"
	"github.com/callummance/azunyan/models"
)

func GetPriority(m *KaraokeManager, request models.Request) float64 {
	var prio float64
	prio = float64(request.PriorityMod)

	previousSongReqs := db.GetPreviousRequestsBySong(m, request.Song, request.ReqTime)

	for _, match := range previousSongReqs {
		timeDiff := request.ReqTime.Sub(match.ReqTime).Seconds()
		if timeDiff < 1 {
			timeDiff = 1
		}
		prioDecrement := float64(m.GetConfig().KaraokeConfig.TimeMultiplier*60) / float64(timeDiff)
		prio -= prioDecrement
	}

	prio += time.Now().Sub(request.ReqTime).Minutes() * float64(m.GetConfig().KaraokeConfig.WaitMultiplier)
	return prio
}

func UpdatePriorities(m *KaraokeManager) []models.Request {
	queuedSongs := db.GetQueued(m)
	for _, song := range queuedSongs {
		song.Priority = GetPriority(m, song)
		db.UpdateReqPrio(m, song.ReqId, song.Priority)
	}
	return queuedSongs
}

func GetQueue(m *KaraokeManager) []models.Request {
	return db.GetQueued(m)
}

func GetNextSong(m *KaraokeManager) models.Request {
	UpdatePriorities(m)
	next := GetQueue(m)[0]
	return next
}

func PopNextSong(m *KaraokeManager) models.Request {
	res := GetNextSong(m)
	db.SetRequestPlayed(m, res.ReqId, time.Now())
	var origState models.State
	curState, err := db.GetEngineState(m, m.GetConfig().KaraokeConfig.SessionName)
	if err != nil {
		return res
	} else {
		origState = *curState
	}
	origState.NowPlaying = &res
	db.UpdateEngineState(m, origState)
	UpdateListenersQueue(m)
	UpdateListenersCur(m)
	return res
}
