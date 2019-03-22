package manager

import (
	"io/ioutil"
	"sort"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/callummance/azunyan/db"
	"github.com/callummance/azunyan/models"
	//"gopkg.in/mgo.v2/bson"
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
	var matches []primitive.ObjectID
	matches = append(matches, m.SearchSongLibrary(searchTerm)...)
	res, _ := db.GetSongsByIDs(m, matches)
	return &res
}

//SearchSongLibrary returns a list of song IDs which match one or more words of
//the given search term, with results matching the most words at the start.
func (m *KaraokeManager) SearchSongLibrary(searchTerm string) []primitive.ObjectID {
	itemsToSearch := strings.Split(searchTerm, " ")
	m.Logger.Printf("Now searching for %v", itemsToSearch)
	matchCount := make(map[primitive.ObjectID]int)
	exactMatchCount := make(map[primitive.ObjectID]int)
	for _, word := range itemsToSearch {
		if word != "" {
			matches := m.TitleSearch.GetMatchingKeys(word, 1)
			matches = mergeMatchSets(matches, m.ArtistSearch.GetMatchingKeys(word, 1))
			matches = mergeMatchSets(matches, m.SourceSearch.GetMatchingKeys(word, 1))
			for match, distance := range matches {
				if distance == 0 {
					exactMatchCount[match.(primitive.ObjectID)]++
				} else {
					matchCount[match.(primitive.ObjectID)]++
				}
			}
		}
	}

	matches := make(map[int]([]primitive.ObjectID))
	exactMatches := make(map[int]([]primitive.ObjectID))
	for id, count := range matchCount {
		matches[count] = append(matches[count], id)
	}
	for id, count := range exactMatchCount {
		exactMatches[count] = append(exactMatches[count], id)
	}

	var res []primitive.ObjectID
	for i := len(itemsToSearch); i > 0; i-- {
		sort.Sort(SearchResults(exactMatches[i]))
		res = append(res, exactMatches[i]...)
	}
	sort.Sort(SearchResults(matches[len(itemsToSearch)]))
	res = append(res, matches[len(itemsToSearch)]...)
	return res
}

func mergeMatchSets(a, b map[interface{}]int) map[interface{}]int {
	for k, va := range a {
		if vb, ok := b[k]; ok {
			//Key is also in b, so take the lower
			if vb < va {
				b[k] = vb
			} else {
				b[k] = va
			}
		} else {
			//Key is only in a, so add to b
			b[k] = va
		}
	}
	return b
}
