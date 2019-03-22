package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

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

	db.InitialiseState(&env)

	//Start listening for web requests
	router := webserver.Route(env)
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil || port == 0 {
		port = env.Config.WebConfig.Port
	}
	router.Run(fmt.Sprintf(":%d", port))
}
