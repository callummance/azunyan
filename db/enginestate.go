package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/callummance/azunyan/models"
	"github.com/globalsign/mgo/bson"
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

func InitialiseState(env databaseConfig) {
	UpdateEngineState(env, models.InitSession(env.GetConfig()))
}

func getStateCollection(env databaseConfig) *mongo.Collection {
	return env.GetClient().Database(env.GetConfig().DbConfig.DatabaseName).Collection("state")
}
