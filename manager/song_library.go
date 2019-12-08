package manager

import (
	"io/ioutil"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/callummance/azunyan/db"
	"github.com/callummance/azunyan/models"
)

//GetSongCoverImage returns a bytestring containing the cover image for the
//given song id
func GetSongCoverImage(id string, m *KaraokeManager) []byte {
	if id == "undefined" {
		return defaultAlbumCover(m)
	}
	sid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		m.Logger.Printf("Failed to convert hex string to object ID for song id %q: %v", id, err)
	}
	bs, err := db.GetSongCoverByID(m, sid)
	if err != nil {
		m.Logger.Printf("Failed to get cover image for song id %q: %v", id, err)
	}
	if bs == nil {
		return defaultAlbumCover(m)
	}
	return bs
}

func defaultAlbumCover(m *KaraokeManager) []byte {
	path := m.Config.KaraokeConfig.DefaultAlbumCover
	res, err := ioutil.ReadFile(path)
	if err != nil {
		m.Logger.Fatalf("Failed to read default album cover from local file storage")
	}
	return res
}

//SearchResults implements sort.Interface for []primitive.ObjectID
type SearchResults []primitive.ObjectID

func (r SearchResults) Len() int           { return len(r) }
func (r SearchResults) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r SearchResults) Less(i, j int) bool { return r[i].String() < r[j].String() }

//GetSearchResults retrieves a list of songs matching the given search. matches is a list of objectIDs
func (m *KaraokeManager) GetSearchResults(searchTerm string) *[]models.Song {
	res, _ := db.GetSongsByTextSearch(m, searchTerm)
	return &res
}
