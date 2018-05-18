package stream

import (
	"io"

	"github.com/callummance/azunyan/db"
	"github.com/callummance/azunyan/manager"
	"github.com/callummance/azunyan/models"
	"github.com/gin-gonic/gin"
)

func GetSub(c *gin.Context) {
	m, ok := c.MustGet("manager").(*manager.KaraokeManager)
	if !ok {
		m.Logger.Printf("Failed to grab environment from Context variable")
		c.String(500, "{\"message\": \"internal failure\"")
	}

	listener := m.SubscribeToChanges()
	defer m.Unsubscribe(listener)
	sendInitial(m, listener)

	c.Stream(func(w io.Writer) bool {
		var update manager.BroadcastData
		update, ok := (<-listener).(manager.BroadcastData)
		if ok {
			c.SSEvent(update.Name, update.Content)
			return true
		} else {
			env := c.MustGet("manager").(*manager.KaraokeManager)
			env.Logger.Printf("Some junk got into the channel!")
			return true
		}
	})
}

func sendInitial(m *manager.KaraokeManager, listener chan interface{}) {
	completeQueue, partialQueue, err := manager.GetQueue(m)
	if err != nil {
		m.Logger.Printf("Failed to get song queue state due to error %v", err)
	}

	state, err := db.GetEngineState(m, m.Config.KaraokeConfig.SessionName)
	if err != nil {
		m.Logger.Printf("Failed to get session data due to error %q", err)
	}

	listener <- manager.BroadcastData{
		Name: "queue",
		Content: map[string][]models.QueueItem{
			"complete": completeQueue,
			"partial":  partialQueue,
		},
	}
	listener <- manager.BroadcastData{
		Name:    "active",
		Content: state.IsActive,
	}
	if state != nil && state.NowPlaying != nil {
		listener <- manager.BroadcastData{
			Name:    "cur",
			Content: state.NowPlaying,
		}
	}
}
