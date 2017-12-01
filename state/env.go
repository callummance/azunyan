package state

import (
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"github.com/callummance/azunyan/db"
	"github.com/callummance/azunyan/config"
)

type Env struct {
	DbSession 		*mgo.Session
	Logger    		*log.Logger
	Config    		config.Config
}

func (e *Env) GetDbConfig() config.DbConfig {
	return e.Config.DbConfig
}

func (e *Env) GetLog() *log.Logger {
	return e.Logger
}

func (e *Env) GetSession() *mgo.Session {
	return e.DbSession
}

func (e *Env) GetConfig() config.Config {
	return e.Config
}

func Initialize(configLoc string) Env {
	logger := log.New(os.Stderr, "AZUNYAN: ", log.Lshortfile|log.Ldate|log.Ltime)
	conf := config.LoadConfig(configLoc, logger)
	session := db.InitDB(conf, logger)

	return Env{DbSession: session, Logger: logger, Config: conf}
}

func (e *Env) UpdateSession() *Env {
	newSession := db.GetNewSession(e)
	return &Env{DbSession: newSession, Logger: e.Logger, Config: e.Config}
}

func (e *Env) CloseSession() {
	db.CloseSession(e)
}

