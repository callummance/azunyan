package models

import (
	"github.com/callummance/azunyan/config"
)

type State struct {
	SessionName		string			`json:"sessionname" bson:"sessionname"`
	NowPlaying		*Request		`json:"nowplaying" bson:"nowplaying,omitempty"`
	NoSingers		int				`json:"nosingers" bson:"nosingers"`
	IsActive		bool			`json:"isactive" bson:"isactive"`
	RequestsActive	bool			`json:"reqactive" bson:"reqactive"`
}

func InitSession(conf config.Config) State {
	return State {
		SessionName: conf.KaraokeConfig.SessionName,
		NowPlaying: nil,
		NoSingers: conf.KaraokeConfig.NoSingers,
		IsActive: false,
		RequestsActive: false,
	}
}