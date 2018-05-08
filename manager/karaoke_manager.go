package manager

import (
	"log"
	"os"

	"github.com/callummance/azunyan/config"
	"github.com/callummance/azunyan/db"
	broadcast "github.com/dustin/go-broadcast"
	mgo "gopkg.in/mgo.v2"
)

type KaraokeManager struct {
	DbSession         *mgo.Session
	Logger            *log.Logger
	Config            config.Config
	UpdateSubscribers broadcast.Broadcaster
}

func Initialize(configLoc string) KaraokeManager {
	logger := log.New(os.Stderr, "AZUNYAN: ", log.Lshortfile|log.Ldate|log.Ltime)
	conf := config.LoadConfig(configLoc, logger)
	session := db.InitDB(conf, logger)
	subscriberBc := broadcast.NewBroadcaster(broadcastBufLen)

	return KaraokeManager{DbSession: session, Logger: logger, Config: conf, UpdateSubscribers: subscriberBc}
}

func (m *KaraokeManager) UpdateSession() *KaraokeManager {
	newSession := db.GetNewSession(m)
	return &KaraokeManager{DbSession: newSession, Logger: m.Logger, Config: m.Config, UpdateSubscribers: m.UpdateSubscribers}
}

func (m *KaraokeManager) CloseSession() {
	db.CloseSession(m)
}

func (m *KaraokeManager) GetSession() *mgo.Session {
	return m.DbSession
}

func (m *KaraokeManager) GetConfig() config.Config {
	return m.Config
}

func (m *KaraokeManager) GetLog() *log.Logger {
	return m.Logger
}
