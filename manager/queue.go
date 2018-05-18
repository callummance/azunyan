package manager

import (
	"sort"
	"time"

	"github.com/callummance/azunyan/db"
	"github.com/callummance/azunyan/models"
)

//GetQueue returns an ordered list of wating queue items and an ordered list of
//requests still waiting for more singers.
func GetQueue(m *KaraokeManager) ([]models.QueueItem, []models.QueueItem, error) {
	enqueuedSongs, err := db.GetLiveAggregatedSongRequests(m)
	if err != nil {
		m.Logger.Printf("Failed to get list of queued songs due to error %v", err)
		return nil, nil, err
	}
	state, err := db.GetEngineState(m, m.Config.KaraokeConfig.SessionName)
	if err != nil {
		m.Logger.Printf("Failed to get engine statedue to error %v", err)
		return nil, nil, err
	}

	//Convert reqs to queueItems and sore into incomplete and complete queue items
	var filledItems []models.QueueItem
	var waitingItems []models.QueueItem
	for sid, rs := range enqueuedSongs {
		song, err := db.GetSongByID(m, sid)
		if err != nil {
			m.Logger.Printf("Failed to get details of queued song %v due to error %v", sid, err)
			continue
		}
		qis, incompleteQi := models.CompileQueueItems(rs, song, state.NoSingers)
		if qis != nil && len(qis) != 0 {
			filledItems = append(filledItems, qis...)
		}
		if incompleteQi != nil {
			waitingItems = append(waitingItems, *incompleteQi)
		}
	}

	//Sort the waiting lists
	now := time.Now()
	sort.Slice(filledItems, func(i int, j int) bool {
		return getWaitingTime(&filledItems[i], now) > getWaitingTime(&filledItems[j], now)
	})
	sort.Slice(waitingItems, func(i int, j int) bool {
		return getWaitingTime(&waitingItems[i], now) > getWaitingTime(&waitingItems[j], now)
	})

	return filledItems, waitingItems, nil
}

func getWaitingTime(req *models.QueueItem, now time.Time) int {
	totalTime := 0
	for _, t := range req.RequestTimes {
		diff := now.Sub(t)
		secondsWaited := int(diff.Seconds())
		totalTime += secondsWaited
	}
	return totalTime
}

//GetNextSong returns the next song which should be played
func GetNextSong(m *KaraokeManager) (*models.QueueItem, error) {
	cq, iq, err := GetQueue(m)
	if err != nil {
		m.Logger.Printf("Failed to get next queued song due to error %v", err)
		return nil, err
	}
	return getNext(cq, iq), nil
}

func getNext(cq []models.QueueItem, iq []models.QueueItem) *models.QueueItem {
	if cq != nil && len(cq) != 0 {
		return &cq[0]
	} else if iq != nil && len(iq) != 0 {
		return &iq[0]
	} else {
		return nil
	}
}

//PopNextSong fetches the next song to be played and marks it as played. It
//also updates all listeners with a new song list and now playing.
func PopNextSong(m *KaraokeManager) (*models.QueueItem, error) {
	completeQueue, partialQueue, err := GetQueue(m)
	if err != nil {
		m.Logger.Printf("Failed to get song queue state due to error %v", err)
		return nil, err
	}
	next := getNext(completeQueue, partialQueue)
	markQueueItemPlayed(m, next)
	var origState models.State
	curState, err := db.GetEngineState(m, m.GetConfig().KaraokeConfig.SessionName)
	if err != nil {
		m.Logger.Printf("Failed to get karaoke engine state due to error %v", err)
		return nil, err
	} else {
		origState = *curState
	}
	origState.NowPlaying = next
	db.UpdateEngineState(m, origState)
	UpdateListenersQueue(m, completeQueue, partialQueue)
	UpdateListenersCur(m, next)
	return next, nil
}

func markQueueItemPlayed(m *KaraokeManager, i *models.QueueItem) {
	for _, rid := range i.RequestIDs {
		db.SetRequestPlayed(m, rid, time.Now())
	}
}
