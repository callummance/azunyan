package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/BurntSushi/toml"
)

type WebConfig struct {
	Port int `toml:"serverport"`
}

type KaraokeConfig struct {
	SessionName       string `toml:"sessionname"`
	NoSingers         int    `toml:"nosingers"`
	TimeMultiplier    int    `toml:"timemultiplier"`
	WaitMultiplier    int    `toml:"waitmultiplier"`
	DefaultAlbumCover string `toml:"defaultcoverimage"`
}

type DbConfig struct {
	DatabaseAddress        string `toml:"dbaddr"`
	DatabaseName           string `toml:"dbname"`
	DatabaseCollectionName string `toml:"dbcollection"`
}

type Config struct {
	DbConfig      DbConfig      `toml:"dbconfig"`
	WebConfig     WebConfig     `toml:"webconfig"`
	KaraokeConfig KaraokeConfig `toml:"karaokeconfig"`
}

func LoadConfig(loc string, logger *log.Logger) Config {
	var res Config
	if _, err := toml.DecodeFile(loc, &res); err != nil {
		logger.Fatal(err)
	}
	err := godotenv.Load()
	if err == nil {
		res.DbConfig.DatabaseAddress = os.Getenv("dbaddr")
		res.DbConfig.DatabaseCollectionName = os.Getenv("dbcollection")
		res.DbConfig.DatabaseName = os.Getenv("dbname")
	}

	return res
}
