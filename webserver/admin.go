package webserver

import (
	"net"
	"strconv"

	"github.com/callummance/azunyan/manager"
	"github.com/gin-gonic/gin"
)

func RouteAdmin(group *gin.RouterGroup) {
	group.POST("/active", activateEndpoint)
	group.POST("/req_active", activateReqEndpoint)
	group.POST("/advance", advanceEndpoint)
	group.POST("/remove_singer", removeSingerEndpoint)
	group.POST("/reset_queue", resetQueueEndpoint)
	group.POST("/singers/:number", changeNumberOfSingersEndpoint)
	group.POST("/allowdupes/:bool", allowDuplicatesEndpoint)
	group.GET("/ipaddress", ipAddressEndpoint)
}

func resetQueueEndpoint(c *gin.Context) {
	env, ok := c.MustGet("manager").(*manager.KaraokeManager)
	if !ok {
		env.Logger.Printf("Failed to grab environment from Context variable")
		c.String(500, "{\"message\": \"internal failure\"")
	}

	err := manager.Reset(env)
	if err != nil {
		c.Error(err)
	} else {
		c.Status(201)
	}
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

func changeNumberOfSingersEndpoint(c *gin.Context) {
	env, ok := c.MustGet("manager").(*manager.KaraokeManager)
	if !ok {
		env.Logger.Printf("Failed to grab environment from Context variable")
		c.String(500, "{\"message\": \"internal failure\"")
	}
	singersString := c.Param("number")
	singers, _ := strconv.Atoi(singersString)
	manager.ChangeNumberOfSingers(env, singers)
	c.Status(201)
}

func allowDuplicatesEndpoint(c *gin.Context) {
	env, ok := c.MustGet("manager").(*manager.KaraokeManager)
	if !ok {
		env.Logger.Printf("Failed to grab environment from Context variable")
		c.String(500, "{\"message\": \"internal failure\"")
	}
	allowDupesString := c.Param("bool")
	allowDupes, _ := strconv.ParseBool(allowDupesString)
	manager.ChangeAllowDuplication(env, allowDupes)
	c.Status(201)
}

func ipAddressEndpoint(c *gin.Context) {
	env, ok := c.MustGet("manager").(*manager.KaraokeManager)
	if !ok {
		env.Logger.Printf("Failed to grab environment from Context variable")
		c.String(500, "{\"message\": \"internal failure\"")
	}
	c.String(200, getIPAddress(env))
}

func getIPAddress(env *manager.KaraokeManager) string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		env.Logger.Printf("%q", err)
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
