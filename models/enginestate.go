package models

import (
	"github.com/callummance/azunyan/config"
)

//State contains data on the current state of the karaoke manager.
type State struct {
	SessionName    string     `json:"sessionname" bson:"sessionname"`
	NowPlaying     *QueueItem `json:"nowplaying" bson:"nowplaying,omitempty"`
	NoSingers      int        `json:"nosingers" bson:"nosingers"`
	IsActive       bool       `json:"isactive" bson:"isactive"`
	RequestsActive bool       `json:"reqactive" bson:"reqactive"`
	AllowingDupes  bool       `json:"allowdupes" bson:"allowdupes"`
}

//InitSession creates a new State struct with sane defaults given a config struct.
func InitSession(conf config.Config) State {
	return State{
		SessionName:    conf.KaraokeConfig.SessionName,
		NowPlaying:     nil,
		NoSingers:      conf.KaraokeConfig.NoSingers,
		IsActive:       false,
		RequestsActive: false,
		AllowingDupes:  conf.KaraokeConfig.AllowDupes,
	}
}
