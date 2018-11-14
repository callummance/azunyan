package db

import (
	"log"
	"github.com/callummance/azunyan/config"
	// "gopkg.in/mgo.v2"
	"github.com/globalsign/mgo"
)

type databaseConfig interface {
	GetSession() *mgo.Session
	GetConfig() config.Config
	GetLog() *log.Logger
}

func InitDB(config config.Config, log *log.Logger) *mgo.Session {
	session, err := mgo.Dial(config.DbConfig.DatabaseAddress)
	if err != nil {
		log.Fatalf("Failed to connect to database %v due to error '%v'", config.DbConfig.DatabaseAddress, err)
	}

	return session
}

func GetNewSession(conf databaseConfig) *mgo.Session {
	s := conf.GetSession()

	return s.Copy()
}

func CloseSession(conf databaseConfig) {
	conf.GetSession().Close()
}
