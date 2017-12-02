package manager

import (
	"github.com/callummance/azunyan/models"
	"github.com/callummance/azunyan/db"
	"github.com/callummance/azunyan/state"
	"time"
)

func GetPriority(env state.Env, request models.Request) float64 {
	var prio float64
	prio = float64(request.PriorityMod)

	previousSongReqs := db.GetPreviousRequestsBySong(&env, request.Song, request.ReqTime)

	var previousSingerReqs []models.Request
	for _, singer := range request.Singers {
		previousSingerReqs = append(previousSingerReqs, db.GetPreviousRequestsBySinger(&env, singer, request.ReqTime)...)
	}

	matches := append(previousSongReqs, previousSingerReqs...)

	for _, match := range matches {
		timeDiff := request.ReqTime.Sub(match.ReqTime).Minutes()
		_ = float64(env.GetConfig().KaraokeConfig.TimeMultiplier)/timeDiff
		//prio -= prioDecrement
	}

	prio += time.Now().Sub(request.ReqTime).Minutes() * float64(env.GetConfig().KaraokeConfig.WaitMultiplier)
	return prio
}

func UpdatePriorities(env state.Env) []models.Request {
	queuedSongs := db.GetQueued(&env)
	for _, song := range queuedSongs {
		song.Priority = GetPriority(env, song)
		db.UpdateReqPrio(&env, song.ReqId, song.Priority)
	}
	return queuedSongs
}

func GetQueue(env state.Env) []models.Request {
	return db.GetQueued(&env)
}

func GetNextSong(env state.Env) models.Request {
	UpdatePriorities(env)
	next := GetQueue(env)[0]
	return next
}

func PopNextSong(env state.Env) models.Request {
	res := GetNextSong(env)
	db.SetRequestPlayed(&env, res.ReqId, time.Now())
	var origState models.State
	curState, err := db.GetEngineState(&env, env.GetConfig().KaraokeConfig.SessionName)
	if err != nil {
		return res
	} else {
		origState = *curState
	}
	origState.NowPlaying = &res
	db.UpdateEngineState(&env, origState)
	UpdateListenersQueue(&env)
	UpdateListenersCur(&env)
	return res
}
