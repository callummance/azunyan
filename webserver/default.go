package webserver

import (
	"github.com/callummance/azunyan/db"
	"github.com/callummance/azunyan/manager"
	"github.com/gin-gonic/gin"
)

//ForwardRoot handles the root site, and forwards to either the song listing or
//the requests page depending on whether requests have been enabled.
func ForwardRoot(c *gin.Context) {
	env, ok := c.MustGet("manager").(*manager.KaraokeManager)
	if !ok {
		env.Logger.Printf("Failed to grab environment from Context variable")
		c.String(500, "{\"message\": \"internal failure\"")
	}
	state, err := db.GetEngineState(env, env.Config.KaraokeConfig.SessionName)
	if err != nil {
		env.Logger.Printf("Failed to get manager state due to error %q", err)
	}
	if state.RequestsActive {
		c.Redirect(302, "/static/request/index.html")
	} else {
		c.Redirect(302, "/static/songlist")
	}
}
