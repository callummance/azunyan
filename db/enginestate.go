package db

import (
	"fmt"

	"github.com/callummance/azunyan/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func UpdateEngineState(env databaseConfig, newState models.State) error {
	col := getStateCollection(env)
	_, err := col.Upsert(bson.M{"sessionname": newState.SessionName}, newState)
	if err != nil {
		env.GetLog().Printf("Database failure whist updating state: %s", err)
		return fmt.Errorf("could not update session state due to error %s", err)
	} else {
		return nil
	}
}

func GetEngineState(env databaseConfig, sessionName string) (*models.State, error) {
	var res models.State
	col := getStateCollection(env)
	q := col.Find(bson.M{"sessionname": sessionName})
	err := q.One(&res)
	if err != nil {
		env.GetLog().Printf("Database failure whist fetching state: %s", err)
		return nil, fmt.Errorf("could not retrieve session state due to error %s", err)
	} else {
		return &res, nil
	}
}

func InitialiseState(env databaseConfig) {
	UpdateEngineState(env, models.InitSession(env.GetConfig()))
}

func getStateCollection(env databaseConfig) *mgo.Collection {
	return env.GetSession().DB(env.GetConfig().DbConfig.DatabaseName).C("state")
}
