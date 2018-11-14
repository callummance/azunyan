package manager

import (
	"io/ioutil"
	"sort"
	"strings"

	"github.com/callummance/azunyan/db"
	"github.com/callummance/azunyan/models"
	//"gopkg.in/mgo.v2/bson"
	"github.com/globalsign/mgo/bson"
)

//GetSongCoverImage returns a bytestring containing the cover image for the
//given song id
func GetSongCoverImage(id string, m *KaraokeManager) []byte {
	sid := bson.ObjectIdHex(id)
	bs, err := db.GetSongCoverByID(m, sid)
	if err != nil {
		m.Logger.Printf("Failed to get cover image for song id %q: %v", id, err)
	}
	if bs == nil {
		path := m.Config.KaraokeConfig.DefaultAlbumCover
		res, err := ioutil.ReadFile(path)
		if err != nil {
			m.Logger.Fatalf("Failed to get default cover image for song id %q: %v", id, err)
		}
		return res
	}
	return bs
}

//SearchResults implemeents sort.Interface for []bson.Objectid
type SearchResults []bson.ObjectId

func (r SearchResults) Len() int           { return len(r) }
func (r SearchResults) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r SearchResults) Less(i, j int) bool { return r[i] < r[j] }

//GetSearchResults retrieves a list of songs matching the given search
func (m *KaraokeManager) GetSearchResults(searchTerm string) []*models.Song {
	matches := m.SearchSongLibrary(searchTerm)
	res := []*models.Song{}
	for _, match := range matches {
		matchingSong, err := db.GetSongByID(m, match)
		if err != nil {
			m.Logger.Printf("Could not find a song matching id %v", match)
		} else {
			res = append(res, matchingSong)
		}
	}
	return res
}

//SearchSongLibrary returns a list of song IDs which match one or more words of
//the given search term, with results matching the most words at the start.
func (m *KaraokeManager) SearchSongLibrary(searchTerm string) []bson.ObjectId {
	itemsToSearch := strings.Split(searchTerm, " ")
	m.Logger.Printf("Now searching for %v", itemsToSearch)
	matchCount := make(map[bson.ObjectId]int)
	exactMatchCount := make(map[bson.ObjectId]int)
	for _, word := range itemsToSearch {
		if word != "" {
			matches := m.TitleSearch.GetMatchingKeys(word, 1)
			matches = mergeMatchSets(matches, m.ArtistSearch.GetMatchingKeys(word, 1))
			matches = mergeMatchSets(matches, m.SourceSearch.GetMatchingKeys(word, 1))
			for match, distance := range matches {
				if distance == 0 {
					exactMatchCount[match.(bson.ObjectId)]++
				} else {
					matchCount[match.(bson.ObjectId)]++
				}
			}
		}
	}

	matches := make(map[int]([]bson.ObjectId))
	exactMatches := make(map[int]([]bson.ObjectId))
	for id, count := range matchCount {
		matches[count] = append(matches[count], id)
	}
	for id, count := range exactMatchCount {
		exactMatches[count] = append(exactMatches[count], id)
	}

	var res []bson.ObjectId
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
