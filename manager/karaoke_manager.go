package manager

import (
	"log"
	"os"

	"github.com/callummance/azunyan/config"
	"github.com/callummance/azunyan/db"
	"github.com/callummance/fuwafuwasearch/levenshteinmatrix"
	broadcast "github.com/dustin/go-broadcast"

	// mgo "gopkg.in/mgo.v2"
	// mgo "github.com/globalsign/mgo"
	mgo "go.mongodb.org/mongo-driver/mongo"
)

// KaraokeManager object
type KaraokeManager struct {
	DbClient          *mgo.Client
	Logger            *log.Logger
	Config            config.Config
	UpdateSubscribers broadcast.Broadcaster
	TitleSearch       *levenshteinmatrix.LMatrixSearch
	ArtistSearch      *levenshteinmatrix.LMatrixSearch
	SourceSearch      *levenshteinmatrix.LMatrixSearch
}

// Initialize creates a new KaraokeManager object
func Initialize(configLoc string) KaraokeManager {
	var newManager KaraokeManager
	newManager.Logger = log.New(os.Stderr, "AZUNYAN: ", log.Lshortfile|log.Ldate|log.Ltime)
	newManager.Config = config.LoadConfig(configLoc, newManager.Logger)
	newManager.DbClient = db.InitDB(newManager.Config, newManager.Logger)
	newManager.UpdateSubscribers = broadcast.NewBroadcaster(broadcastBufLen)
	songSearchData := db.GetSongTAS(&newManager)
	var ids []interface{}
	titles := []string{}
	artists := []string{}
	sources := []string{}
	for _, song := range songSearchData {
		ids = append(ids, song.ID)
		titles = append(titles, song.Title)
		artists = append(artists, song.Artist)
		sources = append(sources, song.Source)
	}
	newManager.TitleSearch = levenshteinmatrix.NewLMatrixSearch(titles, ids, false)
	newManager.ArtistSearch = levenshteinmatrix.NewLMatrixSearch(artists, ids, false)
	newManager.SourceSearch = levenshteinmatrix.NewLMatrixSearch(sources, ids, false)

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
