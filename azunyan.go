package main

import (
	"github.com/callummance/azunyan/state"
	"github.com/callummance/azunyan/webserver"
	"fmt"
	"github.com/callummance/azunyan/db"
	"github.com/callummance/azunyan/models"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	configLoc = "./azunyan.conf"
)

func main() {
	env := state.Initialize(configLoc)
	(&env).UpdateSession()

	db.ImportJSONSongList(&env, "songs.json")
	db.InitialiseState(&env)

	testReq := models.Request{ReqId:bson.NewObjectId(),
								ReqTime:time.Now(),
								Singers:[]string{"Tail Red"},
								Song: bson.NewObjectId(),
								PriorityMod:0,
								PlayedTime:nil}

	db.InsertRequest(&env, testReq)

	router := webserver.Route(env)
	router.Run(fmt.Sprintf(":%d", env.Config.WebConfig.Port))
}
