package middlewares

import (
	"github.com/callummance/azunyan/manager"
	"github.com/gin-gonic/gin"
)

func AttachEnvironment(manager *manager.KaraokeManager, c *gin.Context) {
	c.Set("manager", manager)
	c.Next()
}
