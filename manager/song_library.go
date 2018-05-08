package manager

import (
	"strings"

	"github.com/callummance/azunyan/db"
	"github.com/callummance/azunyan/models"
	"gopkg.in/mgo.v2/bson"
)

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
	matchCount := make(map[bson.ObjectId]int)
	exactMatchCount := make(map[bson.ObjectId]int)
	for _, word := range itemsToSearch {
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
		res = append(res, exactMatches[i]...)
	}
	for i := len(itemsToSearch); i > 0; i-- {
		res = append(res, matches[i]...)
	}
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
