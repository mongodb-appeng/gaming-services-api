package database

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

//GameDataT is
type GameDataT struct {
	Score int `json:"score,omitempty" bson:"score,omitempty"`
}

// GameEventT is
type GameEventT struct {
	// ObjectID  *primitive.ObjectID `json:"_id" bson:"_id"`
	ID        string    `json:"id" bson:"id"`     // Type of event
	Name      string    `json:"name" bson:"name"` // Type of event
	Type      string    `json:"type" bson:"type"` // Type of event
	GameID    string    `json:"gameId" bson:"gameId"`
	GamerID   string    `json:"gamerId" bson:"gamerId"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	Data      GameDataT `json:"data" bson:"data"`
}

// AddGameEvents is
func (a *AtlasClientService) AddGameEvents(argDb, argColl string, gameevents *GameEventT) (result *mongo.InsertOneResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("panic -  \n", r)
			err = r.(error)
		}
	}()

	log.Debug("InsertOne gameevents - \n", gameevents)
	result, err = a.client.Database(argDb).Collection(argColl).InsertOne(context.TODO(), *gameevents)
	// result, err = a.client.Database(argDb).Collection(argColl).InsertMany(context.TODO(), *gameevents)
	if err != nil {
		log.Error("insert -  \n", err)
		return
	}
	return
}
