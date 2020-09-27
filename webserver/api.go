package webserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/callummance/azunyan/db"
	"github.com/callummance/azunyan/manager"
	"github.com/callummance/azunyan/webserver/stream"
	"github.com/gin-gonic/gin"
)

func RouteApi(group *gin.RouterGroup) {
	group.GET("/getsongslist", songListEndpoint)
	group.GET("/nosingers", noSingersEndpoint)
	group.GET("/queuestream", stream.GetSub)
	group.GET("/searchsongs", searchSongsEndpoint)
	group.POST("/addrequest", makeRequestEndpoint)
}

func songListEndpoint(c *gin.Context) {
	env, ok := c.MustGet("manager").(*manager.KaraokeManager)
	if !ok {
		env.Logger.Printf("Failed to grab environment from Context variable")
		c.String(500, "{\"message\": \"internal failure\"")
	}
	c.JSON(200, db.GetSongs(env))
}

func noSingersEndpoint(c *gin.Context) {
	env, ok := c.MustGet("manager").(*manager.KaraokeManager)
	if !ok {
		env.Logger.Printf("Failed to grab environment from Context variable")
		c.String(500, "{\"message\": \"internal failure\"")
	}
	state, err := db.GetEngineState(env, env.Config.KaraokeConfig.SessionName)
	if err != nil {
		env.Logger.Printf("Failed to get singer count due to error %q", err)
	}
	c.String(200, fmt.Sprintf("%d", state.NoSingers))
}

func makeRequestEndpoint(c *gin.Context) {
	env, ok := c.MustGet("manager").(*manager.KaraokeManager)
	if !ok {
		env.Logger.Printf("Failed to grab environment from Context variable")
		c.String(500, "{\"message\": \"internal failure\"")
	}
	state, err := db.GetEngineState(env, env.Config.KaraokeConfig.SessionName)
	if err != nil {
		env.Logger.Printf("Failed to get singer count due to error %q", err)
	}

	if state.RequestsActive {
		var reqData struct {
			SongID string `json:"songid"`
			Singer string `json:"singer"`
		}
		err := c.BindJSON(&reqData)
		if err != nil {
			c.AbortWithError(404, fmt.Errorf("invalid request data sent to server"))
			env.Logger.Printf("Failed to make request %v due to error %v", c, err)
		}
		//TODO: Check if duplicate name requesting song, also check if song is already requested.
		err = manager.AddRequest(env, reqData.Singer, reqData.SongID)
		if err != nil {
			c.AbortWithError(500, fmt.Errorf("internal server error encountered: %v", err))
			env.Logger.Printf("Failed to make request due to error %v", err)
		}
	} else {
		c.AbortWithError(403, fmt.Errorf("requests are not open yet"))
	}
}

//TODO: Delete this as currently unused?
//Retrieves JSON object containing the singer name and the song id from a request,
//then adds the request to the queue.
//If a request for the same song by the same person has already been made,
//return an error message unless the relevant flag has been included.
//Otherwise, return a JSON object containing details on the new request's place
//in the queue.
func makeSinglePlayerRequestEndpoint(c *gin.Context) {
	env, ok := c.MustGet("manager").(*manager.KaraokeManager)
	if !ok {
		env.Logger.Printf("Failed to grab environment from Context variable")
		c.String(500, "{\"message\": \"internal failure\"")
	}
	var requestData struct {
		PlayerName  string `json:"singer"`
		SongID      string `json:"songid"`
		ForceRepeat string `json:"force_repeat"`
	}
	c.BindJSON(&requestData)
}

func searchSongsEndpoint(c *gin.Context) {
	m, ok := c.MustGet("manager").(*manager.KaraokeManager)
	if !ok {
		m.Logger.Printf("Failed to grab environment from Context variable")
		c.String(500, "{\"message\": \"internal failure\"")
	}
	state, err := db.GetEngineState(m, m.Config.KaraokeConfig.SessionName)
	if err != nil {
		m.Logger.Printf("searchSongsEndpoint: Failed to get state due to error %q", err)
	}
	songsLastUpdated := state.SongsLastUpdated.Time().Round(time.Second)
	httpDate := songsLastUpdated.Format(time.RFC1123)
	requestModifiedSinceHeader := c.Request.Header.Get("If-Modified-Since")
	if requestModifiedSinceHeader != "" {
		requestModifiedSince, err := time.Parse(time.RFC1123, requestModifiedSinceHeader)
		if err != nil {
			m.Logger.Printf("searchSongsEndpoint: Error parsing the If-Modified-Since date in the request header %q", err)
		} else if songsLastUpdated.Before(requestModifiedSince.Round(time.Second)) ||
			songsLastUpdated.Equal(requestModifiedSince.Round(time.Second)) {
			c.String(http.StatusNotModified, "")
			return
		}
	}
	searchString := c.Request.URL.Query().Get("q")
	c.Header("cache-control", "public, max-age=86400")
	c.Header("Last-Modified", httpDate)
	c.JSON(201, m.GetSearchResults(searchString))
}
