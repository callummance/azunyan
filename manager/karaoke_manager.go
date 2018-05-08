package manager

import (
	"log"
	"os"

	"github.com/callummance/azunyan/config"
	"github.com/callummance/azunyan/db"
	"github.com/callummance/fuwafuwasearch/levenshteinmatrix"
	broadcast "github.com/dustin/go-broadcast"
	mgo "gopkg.in/mgo.v2"
)

type KaraokeManager struct {
	DbSession         *mgo.Session
	Logger            *log.Logger
	Config            config.Config
	UpdateSubscribers broadcast.Broadcaster
	TitleSearch       *levenshteinmatrix.LMatrixSearch
	ArtistSearch      *levenshteinmatrix.LMatrixSearch
	SourceSearch      *levenshteinmatrix.LMatrixSearch
}

func Initialize(configLoc string) KaraokeManager {
	var newManager KaraokeManager
	newManager.Logger = log.New(os.Stderr, "AZUNYAN: ", log.Lshortfile|log.Ldate|log.Ltime)
	newManager.Config = config.LoadConfig(configLoc, newManager.Logger)
	newManager.DbSession = db.InitDB(newManager.Config, newManager.Logger)
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

func (m *KaraokeManager) UpdateSession() *KaraokeManager {
	newSession := db.GetNewSession(m)
	return &KaraokeManager{
		DbSession:         newSession,
		Logger:            m.Logger,
		Config:            m.Config,
		UpdateSubscribers: m.UpdateSubscribers,
		TitleSearch:       m.TitleSearch,
		ArtistSearch:      m.ArtistSearch,
		SourceSearch:      m.SourceSearch}
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
