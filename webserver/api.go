package webserver

import (
	"fmt"

	"github.com/callummance/azunyan/db"
	"github.com/callummance/azunyan/manager"
	"github.com/callummance/azunyan/webserver/stream"
	"github.com/gin-gonic/gin"
)

func RouteApi(group *gin.RouterGroup) {
	group.GET("/getsongslist", songListEndpoint)
	group.GET("/nosingers", noSingersEndpoint)
	group.GET("/queuestream", stream.GetSub)
	group.GET("/searchsongs", searchSongsEndpoint)
	group.POST("/addrequest", makeRequestEndpoint)
}

func songListEndpoint(c *gin.Context) {
	env, ok := c.MustGet("manager").(*manager.KaraokeManager)
	if !ok {
		env.Logger.Printf("Failed to grab environment from Context variable")
		c.String(500, "{\"message\": \"internal failure\"")
	}
	c.JSON(200, db.GetSongs(env))
}

func noSingersEndpoint(c *gin.Context) {
	env, ok := c.MustGet("manager").(*manager.KaraokeManager)
	if !ok {
		env.Logger.Printf("Failed to grab environment from Context variable")
		c.String(500, "{\"message\": \"internal failure\"")
	}
	state, err := db.GetEngineState(env, env.Config.KaraokeConfig.SessionName)
	if err != nil {
		env.Logger.Printf("Failed to get singer count due to error %q", err)
	}
	c.String(200, fmt.Sprintf("%d", state.NoSingers))
}

func makeRequestEndpoint(c *gin.Context) {
	env, ok := c.MustGet("manager").(*manager.KaraokeManager)
	if !ok {
		env.Logger.Printf("Failed to grab environment from Context variable")
		c.String(500, "{\"message\": \"internal failure\"")
	}
	state, err := db.GetEngineState(env, env.Config.KaraokeConfig.SessionName)
	if err != nil {
		env.Logger.Printf("Failed to get singer count due to error %q", err)
	}

	if state.RequestsActive {
		singers := c.PostFormArray("singers[]")
		song := c.PostForm("songId")
		manager.AddRequest(env, singers, song)
	} else {
		c.AbortWithError(403, fmt.Errorf("Requests are not open yet!"))
	}
}

func searchSongsEndpoint(c *gin.Context) {
	m, ok := c.MustGet("manager").(*manager.KaraokeManager)
	if !ok {
		m.Logger.Printf("Failed to grab environment from Context variable")
		c.String(500, "{\"message\": \"internal failure\"")
	}

	searchString := c.Request.URL.Query().Get("q")
	c.JSON(200, m.GetSearchResults(searchString))
}
