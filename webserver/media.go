package webserver

import (
	"github.com/callummance/azunyan/manager"
	"github.com/gin-gonic/gin"
)

//RouteMedia registers URL handlers for all multimedia functionality to the
//given RouterGroup
func RouteMedia(group *gin.RouterGroup) {
	group.GET("/cover/:albumid", songCoverEndpoint)
}

func songCoverEndpoint(c *gin.Context) {
	env, ok := c.MustGet("manager").(*manager.KaraokeManager)
	if !ok {
		env.Logger.Printf("Failed to grab environment from Context variable")
		c.String(500, "{\"message\": \"internal failure\"")
	}
	albumID := c.Param("albumid")
	cover := manager.GetSongCoverImage(albumID, env)
	c.Data(200, "image/jpeg", cover)
}
