package main

import (
	"flag"
	"fmt"

	"github.com/callummance/azunyan/db"
	"github.com/callummance/azunyan/manager"
	"github.com/callummance/azunyan/webserver"
)

const (
	configLoc = "./azunyan.conf"
)

func main() {
	//Get the conf file location from command line flags
	var confFileLoc string
	flag.StringVar(&confFileLoc, "c", configLoc, "location of the config file")
	flag.Parse()
	fmt.Printf("Starting azunyan with config file %q", confFileLoc)
	//Load the config file
	env := manager.Initialize(confFileLoc)
	(&env).UpdateSession()

	db.InitialiseState(&env)

	//Start listening for web requests
	router := webserver.Route(env)
	router.Run(fmt.Sprintf(":%d", env.Config.WebConfig.Port))
}
