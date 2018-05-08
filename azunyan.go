package main

import (
	"fmt"
	"time"

	"github.com/callummance/azunyan/db"
	"github.com/callummance/azunyan/manager"
	"github.com/callummance/azunyan/models"
	"github.com/callummance/azunyan/webserver"
	"gopkg.in/mgo.v2/bson"
)

const (
	configLoc = "./azunyan.conf"
)

func main() {
	env := manager.Initialize(configLoc)
	(&env).UpdateSession()

	db.InitialiseState(&env)

	testReq := models.Request{ReqId: bson.NewObjectId(),
		ReqTime:     time.Now(),
		Singers:     []string{"Tail Red"},
		Song:        bson.NewObjectId(),
		PriorityMod: 0,
		PlayedTime:  nil}

	db.InsertRequest(&env, testReq)

	router := webserver.Route(env)
	router.Run(fmt.Sprintf(":%d", env.Config.WebConfig.Port))
}
