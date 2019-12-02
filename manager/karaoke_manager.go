package manager

import (
	"log"
	"os"

	"github.com/callummance/azunyan/config"
	"github.com/callummance/azunyan/db"
	broadcast "github.com/dustin/go-broadcast"

	mgo "go.mongodb.org/mongo-driver/mongo"
)

// KaraokeManager object
type KaraokeManager struct {
	DbClient          *mgo.Client
	Logger            *log.Logger
	Config            config.Config
	UpdateSubscribers broadcast.Broadcaster
}

// Initialize creates a new KaraokeManager object
func Initialize(configLoc string) KaraokeManager {
	var newManager KaraokeManager
	newManager.Logger = log.New(os.Stderr, "AZUNYAN: ", log.Lshortfile|log.Ldate|log.Ltime)
	newManager.Config = config.LoadConfig(configLoc, newManager.Logger)
	newManager.DbClient = db.InitDB(newManager.Config, newManager.Logger)
	newManager.UpdateSubscribers = broadcast.NewBroadcaster(broadcastBufLen)
	return newManager
}

// CloseSession closes database connection
func (m *KaraokeManager) CloseSession() {
	db.CloseSession(m)
}

// GetClient returns the database
func (m *KaraokeManager) GetClient() *mgo.Client {
	return m.DbClient
}

// GetConfig returns the configuration stored in the configuration file
func (m *KaraokeManager) GetConfig() config.Config {
	return m.Config
}

// GetLog returns the logger
func (m *KaraokeManager) GetLog() *log.Logger {
	return m.Logger
}
