package webserver

import (
	//"net/http"

	"github.com/callummance/azunyan/state"
	"github.com/callummance/azunyan/webserver/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/callummance/azunyan/webserver/stream"
)

func Route(env state.Env) *gin.Engine {
	stream.InitBroadcaster()
	router := gin.Default()

	//Attach environment struct
	router.Use(func(context *gin.Context) { middlewares.AttachEnvironment(env, context) })

	//Favicon
	router.StaticFile("favicon.ico", "./static/frontend/favicon.ico")
	//Static Files
	router.StaticFS("/static", http.Dir("./static/frontend"))

	//Forward root
	router.GET("/", ForwardRoot)

	//API group
	apig := router.Group("/api")
	RouteApi(apig)

	//Admin group
	adming := router.Group("/admin", gin.BasicAuth(gin.Accounts{
		"admin": "twintailsniconiconi",
	}))
	RouteAdmin(adming)

	return router
}
