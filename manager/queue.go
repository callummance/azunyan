package manager

import (
	"reflect"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/callummance/azunyan/db"
	"github.com/callummance/azunyan/models"
)

/*GetQueue returns an ordered list of waiting queue items and an ordered list of
  requests still waiting for more singers.
*/
func GetQueue(m *KaraokeManager) ([]models.QueueItem, []models.QueueItem, error) {
	requests := db.GetLiveRequests(m)
	upcomingSongs := db.GetUpcomingSongs(m)

	var filledItems []models.QueueItem
	var waitingItems []models.QueueItem
	requestsAdded := make(map[primitive.ObjectID]bool)
	leftoverRequests := make(map[primitive.ObjectID][]models.Request)

	for i := 0; i < len(requests); i++ {
		// Check if in upcoming queue already
		request := requests[i]
		if upcomingqueueitem, present := upcomingSongs[request.ReqID]; present {
			if !requestsAdded[upcomingqueueitem.QueueItem.QueueItemID] {
				queueItem := models.QueueItem{
					RequestIDs:   upcomingqueueitem.QueueItem.RequestIDs,
					QueueItemID:  upcomingqueueitem.QueueItem.QueueItemID,
					SongID:       upcomingqueueitem.QueueItem.SongID,
					SongTitle:    upcomingqueueitem.QueueItem.SongTitle,
					SongArtist:   upcomingqueueitem.QueueItem.SongArtist,
					Singers:      upcomingqueueitem.QueueItem.Singers,
					RequestTimes: upcomingqueueitem.QueueItem.RequestTimes,
				}
				filledItems = append(filledItems, queueItem)
				requestsAdded[upcomingqueueitem.QueueItem.QueueItemID] = true
			}

		} else {
			if rqs, present := leftoverRequests[request.Song]; present {
				newReqs := append(rqs, request)
				leftoverRequests[request.Song] = newReqs
			} else {
				reqSlice := []models.Request{request}
				leftoverRequests[request.Song] = reqSlice
			}
		}
	}
	completedQueue, incompleteQueue, _ := getQueueHelper(m, leftoverRequests)
	filledItems = append(filledItems, completedQueue...)
	waitingItems = append(waitingItems, incompleteQueue...)
	filledItems, waitingItems = sortQueues(filledItems, waitingItems)
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
//completedQueue represents songs that have filled the number of singers
//partialQueue are songs that are awaiting singers
func PopNextSong(m *KaraokeManager) (*models.QueueItem, error) {
	completeQueue, partialQueue, err := GetQueue(m)
	if err != nil {
		m.Logger.Printf("Failed to get song queue state due to error %v", err)
		return nil, err
	}
	next := getNext(completeQueue, partialQueue)
	var origState models.State
	curState, err := db.GetEngineState(m, m.GetConfig().KaraokeConfig.SessionName)
	if err != nil {
		m.Logger.Printf("Failed to get karaoke engine state due to error %v", err)
		return nil, err
	}
	origState = *curState
	if next != nil {
		m.Logger.Printf("Popping next song")
		markQueueItemPlayed(m, next)
		origState.NowPlaying = next
		db.UpdateEngineState(m, origState)
		db.RemoveSongFromUpcomingQueue(m, next.QueueItemID)
		UpdateListenersQueue(m, removeNowPlayingFromList(next, completeQueue), removeNowPlayingFromList(next, partialQueue), 0)
	} else {
		m.Logger.Printf("No more songs left in queue")
		newState := models.InitSession(m.GetConfig())
		newState.IsActive = true
		newState.RequestsActive = true
		newState.NoSingers = origState.NoSingers
		db.ClearEngineState(m)
		db.UpdateEngineState(m, newState)
	}
	UpdateListenersCur(m, next)
	return next, nil
}

//IsDupeRequest returns true iff the given request
//has been requested by enough people before to form a full party.
func IsDupeRequest(m *KaraokeManager, sid primitive.ObjectID) (bool, error) {
	prevRequests, err := db.GetPreviousRequestsBySong(m, sid, time.Now())
	if err != nil {
		m.Logger.Printf("Failed to get previous requests when searching for dupe requests for %v with error %v", sid, err)
	}
	return len(prevRequests) > m.Config.KaraokeConfig.NoSingers, nil
}

//CompileQueueItems takes a list of requests for a song, as well as the song's
//struct itself and the number of singers which can be contained in a single
//queued item and returns a slice of complete queue items along with an optional
//incomplete queue item.
func CompileQueueItems(m *KaraokeManager, reqs []models.Request, song *models.Song, maxSingers int) ([]models.QueueItem, *models.QueueItem) {
	if len(reqs) == 0 {
		return nil, nil
	}

	noRequests := len(reqs)
	incompleteQueueLen := noRequests % maxSingers
	fullQueueItemsCount := noRequests / maxSingers
	completeRes := make([]models.QueueItem, fullQueueItemsCount)

	for i := 0; i < fullQueueItemsCount; i++ {
		item := models.QueueItem{
			QueueItemID:  primitive.NewObjectID(),
			SongID:       song.ID,
			SongTitle:    song.Title,
			SongArtist:   song.Artist,
			RequestIDs:   []primitive.ObjectID{},
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

		for _, reqID := range item.RequestIDs {
			db.AddSongToUpcomingQueue(m, reqID, item)
		}

	}

	if incompleteQueueLen != 0 {
		item := models.QueueItem{
			SongID:       song.ID,
			SongTitle:    song.Title,
			SongArtist:   song.Artist,
			RequestIDs:   []primitive.ObjectID{},
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

func removeNowPlayingFromList(np *models.QueueItem, list []models.QueueItem) []models.QueueItem {
	updated := make([]models.QueueItem, 0)
	for _, v := range list {
		if !reflect.DeepEqual(np.RequestIDs, v.RequestIDs) {
			updated = append(updated, v)
		}
	}
	return updated
}

func markQueueItemPlayed(m *KaraokeManager, i *models.QueueItem) {
	for _, rid := range i.RequestIDs {
		db.SetRequestPlayed(m, rid, time.Now())
	}
}

func getQueueHelper(m *KaraokeManager, enqueuedSongs map[primitive.ObjectID][]models.Request) ([]models.QueueItem, []models.QueueItem, error) {
	//Convert reqs to queueItems and sort into incomplete and complete queue items
	var filledItems []models.QueueItem
	var waitingItems []models.QueueItem

	state, err := db.GetEngineState(m, m.Config.KaraokeConfig.SessionName)
	if err != nil {
		m.Logger.Printf("Failed to get engine statedue to error %v", err)
		return nil, nil, err
	}

	for sid, rs := range enqueuedSongs {
		shouldRemoveBecauseDupe := false
		if !state.AllowingDupes {
			isDupe, err := IsDupeRequest(m, sid)
			if err != nil {
				m.Logger.Printf("Could not check if song %q was a duplicate, allowing request...", sid)
			} else {
				shouldRemoveBecauseDupe = isDupe
			}
		}

		if !shouldRemoveBecauseDupe {
			song, err := db.GetSongByID(m, sid)
			if err != nil {
				m.Logger.Printf("Failed to get details of queued song %v due to error %v", sid, err)
				continue
			}
			completedQueue, incompleteQueue := CompileQueueItems(m, rs, song, state.NoSingers)
			if completedQueue != nil && len(completedQueue) != 0 {
				filledItems = append(filledItems, completedQueue...)
			}
			if incompleteQueue != nil {
				waitingItems = append(waitingItems, *incompleteQueue)
			}
		}
	}
	return filledItems, waitingItems, nil
}

func sortQueues(filledItems []models.QueueItem, waitingItems []models.QueueItem) ([]models.QueueItem, []models.QueueItem) {
	//Sort the waiting lists
	now := time.Now()
	sort.Slice(filledItems, func(i int, j int) bool {
		return getWaitingTime(&filledItems[i], now) > getWaitingTime(&filledItems[j], now)
	})
	sort.Slice(waitingItems, func(i int, j int) bool {
		return getWaitingTime(&waitingItems[i], now) > getWaitingTime(&waitingItems[j], now)
	})
	return filledItems, waitingItems
}
