package webserver

import (
	"github.com/gin-gonic/gin"
	"github.com/callummance/azunyan/state"
	"github.com/callummance/azunyan/db"
	"fmt"
	"github.com/callummance/azunyan/webserver/stream"
	"github.com/callummance/azunyan/manager"
)

func RouteApi (group *gin.RouterGroup) {
	group.GET("/getsongslist", songListEndpoint)
	group.GET("/nosingers", noSingersEndpoint)
	group.GET("/queuestream", stream.GetSub)
	group.POST("/addrequest", makeRequestEndpoint)
}

func songListEndpoint (c *gin.Context) {
	env, ok := c.MustGet("env").(*state.Env)
	if !ok {
		env.Logger.Printf("Failed to grab environment from Context variable")
		c.String(500, "{\"message\": \"internal failure\"")
	}
	c.JSON(200, db.GetSongs(env))
}

func noSingersEndpoint (c *gin.Context) {
	env, ok := c.MustGet("env").(*state.Env)
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

func makeRequestEndpoint (c *gin.Context) {
	env, ok := c.MustGet("env").(*state.Env)
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