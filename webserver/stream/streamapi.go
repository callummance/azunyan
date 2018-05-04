package stream

import (
	"io"

	"github.com/callummance/azunyan/db"
	"github.com/callummance/azunyan/models"
	"github.com/callummance/azunyan/state"
	"github.com/gin-gonic/gin"
)

func GetSub(c *gin.Context) {
	env, ok := c.MustGet("env").(*state.Env)
	if !ok {
		env.Logger.Printf("Failed to grab environment from Context variable")
		c.String(500, "{\"message\": \"internal failure\"")
	}

	listener := SubscribeToChanges()
	defer Unsubscribe(listener)
	sendInitial(env, listener)

	c.Stream(func(w io.Writer) bool {
		var update BroadcastData
		update, ok := (<-listener).(BroadcastData)
		if ok {
			c.SSEvent(update.Name, update.Content)
			return true
		} else {
			env := c.MustGet("env").(*state.Env)
			env.Logger.Printf("Some junk got into the channel!")
			return true
		}
	})
}

func sendInitial(env *state.Env, listener chan interface{}) {
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

	state, err := db.GetEngineState(env, env.Config.KaraokeConfig.SessionName)
	if err != nil {
		env.Logger.Printf("Failed to get session data due to error %q", err)
	}

	listener <- BroadcastData{
		Name:    "queue",
		Content: abbQueue,
	}
	listener <- BroadcastData{
		Name:    "active",
		Content: state.IsActive,
	}
	if state != nil && state.NowPlaying != nil {
		song, err := db.GetSongById(env, state.NowPlaying.Song)
		if err != nil {
			env.Logger.Printf("Error whilst retrieving queued songs list: %q", err)
		}
		listener <- BroadcastData{
			Name:    "cur",
			Content: state.NowPlaying.Abbreviate(*song),
		}
	}
}
