package database

import (
	"context"
	"time"

	"math/rand"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//ImageT is
type ImageT struct {
	S3     string `json:"url" bson:"url"`
	Avatar string `json:"avatar" bson:"avatar"`
	Banner string `json:"banner" bson:"banner"`
}

//OverallStatsT is
type OverallStatsT struct {
	PercentDone float32   `json:"completionRatePercentage" bson:"completionRatePercentage"`
	PlayTime    time.Time `json:"playTime" bson:"playTime"`
	WatchTime   time.Time `json:"watchTime" bson:"watchTime"`
}

//StatsT is
type StatsT struct {
	LastActive   time.Time      `json:"lastActive" bson:"lastActive"`
	OverallStats *OverallStatsT `json:"overall" bson:"overall"`
}

//PlayedGameT is
type PlayedGameT struct {
	ID    string `json:"id" bson:"id"`
	Title string `json:"title" bson:"title"`
	Wins  int32  `json:"wins" bson:"wins"`
	Draws int32  `json:"draw" bson:"draw"`
}

//AchievementT is
type AchievementT struct {
	GameID      string    `json:"gameId" bson:"gameId"`
	GameTitle   string    `json:"gameTitle" bson:"gameTitle"`
	ID          string    `json:"id" bson:"id"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"desc" bson:"desc"`
	Image       string    `json:"img" bson:"img"`
	DateEarned  time.Time `json:"earned" bson:"earned"`
	RarityLevel int       `json:"rarityLevel" bson:"rarityLevel"`
}

// GamerProfileT is
type GamerProfileT struct {
	ObjectID     *primitive.ObjectID `json:"_id" bson:"_id"`
	ID           string              `json:"id" bson:"id"`
	Handle       string              `json:"handle" bson:"handle"`
	Image        *ImageT             `json:"image" bson:"image"`
	Stats        *StatsT             `json:"stats" bson:"stats"`
	Sharing      []string            `json:"shareWith" bson:"shareWith"`
	PlayedGame   []*PlayedGameT      `json:"games" bson:"games"`
	Achievements []*AchievementT     `json:"achievements" bson:"achievements"`
	Location     string              `json:"location" bson:"location"`
	Meta         *MetaT              `json:"meta" bson:"meta"`
}

//GamerHandleT is
type GamerHandleT struct {
	ID   int    `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

// CreateGamerProfile is
func (a *AtlasClientService) CreateGamerProfile(argDb, argColl string, gamerprofile *GamerProfileT) (result *mongo.InsertOneResult, err error) {

	defer func() {
		if r := recover(); r != nil {
			log.Error("panic -  \n", r)
			err = r.(error)
		}
	}()

	// log.Error("gamer profile to insert is \n", *gamerprofile)
	result, err = a.client.Database(argDb).Collection(argColl).InsertOne(context.TODO(), *gamerprofile)
	if err != nil {
		log.Error("insert -  \n", err)
		return
	}
	return
}

// FindGamerProfileByID returns an gamerprofile matching the id provided, if any.
func (a *AtlasClientService) FindGamerProfileByID(id, argDb, argColl string) (result *GamerProfileT, err error) {

	defer func() {
		if r := recover(); r != nil {
			log.Error("panic -  \n", r)
			// log.Fatal(err)
			err = r.(error)
		}
	}()
	// log.Error("looking for gamer profile with id \n", id)
	filter := bson.D{primitive.E{Key: "id", Value: id}}
	err = a.client.Database(argDb).Collection(argColl).FindOne(context.TODO(), filter).Decode(&result)
	return
}

// GetRandomGamerHandle returns a random gamer handle
func (a *AtlasClientService) GetRandomGamerHandle(argDb, argColl string) (result *GamerHandleT, err error) {

	defer func() {
		if r := recover(); r != nil {
			log.Error("panic -  \n", r)
			err = r.(error)
		}
	}()

	total, err := a.client.Database(argDb).Collection(argColl).EstimatedDocumentCount(context.TODO())
	if err != nil {
		return
	}
	id := rand.Intn(int(total)-1) + 1 // 0-based
	filter := bson.D{primitive.E{Key: "id", Value: id}}
	err = a.client.Database(argDb).Collection(argColl).FindOne(context.TODO(), filter).Decode(&result)
	return
}

// FindGamerProfileByAccountID returns an gamerprofile matching the id provided, if any.
func (a *AtlasClientService) FindGamerProfileByAccountID(id, argDb, argColl string) (result *GamerProfileT, err error) {

	defer func() {
		if r := recover(); r != nil {
			log.Error("panic -  \n", r)
			// log.Fatal(err)
			err = r.(error)
		}
	}()

	filter := bson.D{primitive.E{Key: "id", Value: id}}
	err = a.client.Database(argDb).Collection(argColl).FindOne(context.TODO(), filter).Decode(&result)
	return
}

// UpdateGamerProfileByID attempts to update the matching gamerprofile with the provided update
func (a *AtlasClientService) UpdateGamerProfileByID(id, argDb, argColl string, gamerprofile map[string]interface{}) (result *mongo.UpdateResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("panic -  \n", r)
			err = r.(error)
		}
	}()

	filter := bson.D{primitive.E{Key: "id", Value: id}}
	update := bson.D{primitive.E{Key: "$set", Value: gamerprofile}}
	result, err = a.client.Database(argDb).Collection(argColl).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return
	}
	return
}

// UpdatePlayedGame attempts to update the matching gamerprofile with the provided update
func (a *AtlasClientService) UpdatePlayedGame(gamerID, gameID, argDb, argColl string, playedgame *PlayedGameT) (result *mongo.UpdateResult, err error) {

	defer func() {
		if r := recover(); r != nil {
			log.Error("panic -  \n", r)
			err = r.(error)
		}
	}()

	// TODO - test this
	filter := bson.D{primitive.E{Key: "id", Value: gamerID}, primitive.E{Key: "games.id", Value: gameID}}
	// { "$inc": { "game.$" : *playedgame } }
	update := bson.D{primitive.E{Key: "$inc", Value: bson.D{primitive.E{Key: "games.$.wins", Value: playedgame.Wins},
		primitive.E{Key: "games.$.draw", Value: playedgame.Draws}}}}
	result, err = a.client.Database(argDb).Collection(argColl).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error("update -  \n", err)
		return
	}
	return
}

// DeleteGamerProfile is
func (a *AtlasClientService) DeleteGamerProfile(id, argDb, argColl string) (result *mongo.DeleteResult, err error) {

	defer func() {
		if r := recover(); r != nil {
			log.Error("panic -  \n", r)
			err = r.(error)
		}
	}()

	filter := bson.D{primitive.E{Key: "id", Value: id}}
	result, err = a.client.Database(argDb).Collection(argColl).DeleteOne(context.TODO(), filter)
	return
}
