package webserver

import (
	"github.com/gin-gonic/gin"
	"github.com/callummance/azunyan/state"
	"github.com/callummance/azunyan/db"
)

func ForwardRoot(c *gin.Context) {
	env, ok := c.MustGet("env").(*state.Env)
	if !ok {
		env.Logger.Printf("Failed to grab environment from Context variable")
		c.String(500, "{\"message\": \"internal failure\"")
	}
	state, err := db.GetEngineState(env, env.Config.KaraokeConfig.SessionName)
	if err != nil {
		env.Logger.Printf("Failed to get singer count due to error %q", err)
	}
	if (state.RequestsActive) {
		c.Redirect(302, "/static/request/index.html")
	} else {
		c.Redirect(302, "/songlist")
	}

}