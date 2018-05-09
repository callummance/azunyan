package webserver

import (
	//"net/http"

	"net/http"

	"github.com/callummance/azunyan/manager"
	"github.com/callummance/azunyan/webserver/middlewares"
	"github.com/gin-gonic/gin"
)

//Route sets up and returns a router for the webserver
func Route(man manager.KaraokeManager) *gin.Engine {
	router := gin.Default()

	//Attach environment struct
	router.Use(func(context *gin.Context) { middlewares.AttachEnvironment(&man, context) })

	//Favicon
	router.StaticFile("favicon.ico", "./static/frontend/favicon.ico")
	//Static Files
	router.StaticFS("/static", http.Dir("./static/frontend"))

	//Forward root
	router.GET("/", ForwardRoot)

	//API group
	apig := router.Group("/api")
	RouteApi(apig)

	//Image Group
	imgg := router.Group("/i")
	RouteMedia(imgg)

	//Admin group
	adming := router.Group("/admin", gin.BasicAuth(gin.Accounts{
		"admin": "twintailsniconiconi",
	}))
	RouteAdmin(adming)

	return router
}
