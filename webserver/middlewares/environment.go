package middlewares

import (
	"github.com/callummance/azunyan/state"
	"github.com/gin-gonic/gin"
)

func AttachEnvironment(env state.Env, c *gin.Context) {
	newEnv := env.UpdateSession()
	defer newEnv.CloseSession()

	c.Set("env", newEnv)
	c.Next()
}
