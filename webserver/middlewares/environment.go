package middlewares

import (
	"github.com/callummance/azunyan/manager"
	"github.com/gin-gonic/gin"
)

func AttachEnvironment(manager *manager.KaraokeManager, c *gin.Context) {
	newMan := manager.UpdateSession()
	defer newMan.CloseSession()

	c.Set("manager", newMan)
	c.Next()
}
