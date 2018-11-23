package webserver

import (
	"github.com/callummance/azunyan/manager"
	"github.com/gin-gonic/gin"
)

func RouteAdmin(group *gin.RouterGroup) {
	group.POST("/active", activateEndpoint)
	group.POST("/req_active", activateReqEndpoint)
	group.POST("/advance", advanceEndpoint)
	group.POST("/remove_singer", removeSingerEndpoint)
}

type activeRequest struct {
	Active bool `json:"active" form:"active" binding:"required"`
}

func activateEndpoint(c *gin.Context) {
	env, ok := c.MustGet("manager").(*manager.KaraokeManager)
	if !ok {
		env.Logger.Printf("Failed to grab environment from Context variable")
		c.String(500, "{\"message\": \"internal failure\"")
	}

	var payload activeRequest
	c.Bind(&payload)
	err := manager.SetActive(env, payload.Active)
	if err != nil {
		c.Error(err)
	} else {
		c.Status(201)
	}
}

type removeSingerRequest struct {
	Singer string `json:"singer" form:"singer" binding:"required"`
}

func removeSingerEndpoint(c *gin.Context) {
	env, ok := c.MustGet("manager").(*manager.KaraokeManager)
	if !ok {
		env.Logger.Printf("Failed to grab environment from Context variable")
		c.String(500, "{\"message\": \"internal failure\"")
	}

	var payload removeSingerRequest
	c.Bind(&payload)
	err := manager.RemoveSinger(env, payload.Singer)
	if err != nil {
		c.Error(err)
	} else {
		c.Status(201)
	}
}

func activateReqEndpoint(c *gin.Context) {
	env, ok := c.MustGet("manager").(*manager.KaraokeManager)
	if !ok {
		env.Logger.Printf("Failed to grab environment from Context variable")
		c.String(500, "{\"message\": \"internal failure\"")
	}

	var payload activeRequest
	c.Bind(&payload)
	err := manager.SetReqActive(env, payload.Active)
	if err != nil {
		c.Error(err)
	} else {
		c.Status(201)
	}
}

func advanceEndpoint(c *gin.Context) {
	env, ok := c.MustGet("manager").(*manager.KaraokeManager)
	if !ok {
		env.Logger.Printf("Failed to grab environment from Context variable")
		c.String(500, "{\"message\": \"internal failure\"")
	}

	manager.PopNextSong(env)
	c.Status(201)
}
