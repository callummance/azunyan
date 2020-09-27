package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/callummance/azunyan/models"
	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Updates the state collection with new state fields
func UpdateEngineState(env databaseConfig, newState models.State) error {
	col := getStateCollection(env)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	updateOptions := options.Update()
	updateOptions.SetUpsert(true)
	_, err := col.UpdateOne(ctx, bson.M{"sessionname": newState.SessionName}, bson.M{"$set": newState}, updateOptions)
	if err != nil {
		env.GetLog().Printf("Database failure whist updating state: %s", err)
		return fmt.Errorf("could not update session state due to error %s", err)
	}
	return nil
}

// ClearEngineState clears the state collection
func ClearEngineState(env databaseConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := getStateCollection(env).DeleteMany(ctx, bson.M{})
	return err
}

func GetEngineState(env databaseConfig, sessionName string) (*models.State, error) {
	var res models.State
	col := getStateCollection(env)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	singleResult := col.FindOne(ctx, bson.M{"sessionname": sessionName})
	err := singleResult.Err()
	if err != nil {
		env.GetLog().Printf("Database failure whist fetching state: %s", err)
		return nil, fmt.Errorf("could not retrieve session state due to error %s", err)
	} else {
		singleResult.Decode(&res)
		return &res, nil
	}
}

// InitialiseState initialises the session state in the database. If an existing state collection already exists
// then it re-uses the SongsLastUpdated field, otherwise it uses the current time.
func InitialiseState(env databaseConfig) {
	var res models.State
	hasLastUpdatedField := false
	initialSessionStateConfig := models.InitSession(env.GetConfig())

	col := getStateCollection(env)
	if col != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		singleResult := col.FindOne(ctx, bson.M{})
		err := singleResult.Err()
		if err != nil {
			env.GetLog().Printf("Database failure whist fetching state: %s", err)
		} else {
			singleResult.Decode(&res)
			if res.SongsLastUpdated != 0 {
				hasLastUpdatedField = true
			}
		}
	}
	if hasLastUpdatedField {
		initialSessionStateConfig.SongsLastUpdated = res.SongsLastUpdated
	} else {
		initialSessionStateConfig.SongsLastUpdated = primitive.NewDateTimeFromTime(time.Now())
	}
	UpdateEngineState(env, initialSessionStateConfig)
}

func getStateCollection(env databaseConfig) *mongo.Collection {
	return env.GetClient().Database(env.GetConfig().DbConfig.DatabaseName).Collection("state")
}
