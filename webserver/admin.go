package webserver

import (
	"github.com/gin-gonic/gin"
	"github.com/callummance/azunyan/state"
	"github.com/callummance/azunyan/manager"
)

func RouteAdmin(group *gin.RouterGroup) {
	group.POST("/active", activateEndpoint)
	group.POST("/req_active", activateReqEndpoint)
	group.POST("/advance", advanceEndpoint)
}

type activeRequest struct {
	Active	bool 	`json:"active" form:"active" binding:"required"`
}

func activateEndpoint (c *gin.Context) {
	env, ok := c.MustGet("env").(*state.Env)
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

func activateReqEndpoint (c *gin.Context) {
	env, ok := c.MustGet("env").(*state.Env)
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
	env, ok := c.MustGet("env").(*state.Env)
	if !ok {
		env.Logger.Printf("Failed to grab environment from Context variable")
		c.String(500, "{\"message\": \"internal failure\"")
	}

	manager.PopNextSong(*env)
	c.Status(201)
}