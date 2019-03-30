package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/callummance/azunyan/config"
)

// Karaoke Manager implements this interface
type databaseConfig interface {
	GetClient() *mongo.Client
	GetConfig() config.Config
	GetLog() *log.Logger
}

// InitDB connects to the MongoDB database and returns a client object
func InitDB(config config.Config, log *log.Logger) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.DbConfig.DatabaseAddress))
	if err != nil {
		log.Fatalf("Failed to connect to database %v due to error '%v'", config.DbConfig.DatabaseAddress, err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("Failed to connect to database %v due to error '%v'", config.DbConfig.DatabaseAddress, err)
	}
	return client
}

// CloseSession disconnects the client
func CloseSession(conf databaseConfig) {
	conf.GetClient().Disconnect(context.TODO())
}
