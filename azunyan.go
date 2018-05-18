package main

import (
	"fmt"

	"github.com/callummance/azunyan/db"
	"github.com/callummance/azunyan/manager"
	"github.com/callummance/azunyan/webserver"
)

const (
	configLoc = "./azunyan.conf"
)

func main() {
	env := manager.Initialize(configLoc)
	(&env).UpdateSession()

	db.InitialiseState(&env)

	router := webserver.Route(env)
	router.Run(fmt.Sprintf(":%d", env.Config.WebConfig.Port))
}
