package database

import (
	"context"
	
	log "github.com/sirupsen/logrus"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//LimitsT is
type LimitsT struct {
	Min int `json:"min" bson:"min"`
	Max int `json:"max" bson:"max"`
}

//ScoringT is
type ScoringT struct {
	Type          string   `json:"type" bson:"type"`
	LowerIsBetter bool     `json:"lowerIsBetter" bson:"lowerIsBetter"`
	Limits        *LimitsT `json:"limits" bson:"limits"`
	MaxPlayers    int      `json:"maxPlayers" bson:"maxPlayers"`
}

//LeaderboardPlayerT is
type LeaderboardPlayerT struct {
	GamerID string `json:"gamerId" bson:"gamerId"`
	Score   int    `json:"score" bson:"score"`
	Handle  string `json:"handle" bson:"handle"`
}

//LeaderboardSummaryT is
type LeaderboardSummaryT struct {
	MinScore    int `json:"minScore" bson:"minScore"`
	MaxScore    int `json:"maxScore" bson:"maxScore"`
	PlayerCount int `json:"playerCount" bson:"playerCount"`
}

//LeaderboardFrameT is
type LeaderboardFrameT struct {
	EventID     string                `json:"_id" bson:"_id"`
	ID          string                `json:"id" bson:"id"`
	Title       string                `json:"title" bson:"title"`
	Summary     *LeaderboardSummaryT  `json:"summary" bson:"summary"`
	Players     *[]LeaderboardPlayerT `json:"players" bson:"players"`
	LastUpdated time.Time             `json:"lastUpdated" bson:"lastUpdated"`
}

// LeaderboardT is
type LeaderboardT struct {
	ID               string               `json:"id" bson:"id"`
	GameID           string               `json:"gameId" bson:"gameId"`
	Name             string               `json:"name" bson:"name"`
	Icon             string               `json:"icon" bson:"icon"`
	Scoring          *ScoringT            `json:"scoring" bson:"scoring"`
	TimeFrames       *[]string            `json:"timeFrames" bson:"timeFrames"`
	Visibility       string               `json:"visibility" bson:"visibility"`
	LeaderboardFrame []*LeaderboardFrameT `json:"data" bson:"data"`
	Location         string               `json:"location" bson:"location"`
	Meta             *MetaT               `json:"meta" bson:"meta"`
}

// CreateLeaderboard is
func (a *AtlasClientService) CreateLeaderboard(argDb, argColl string, leaderboard *LeaderboardT) (result *mongo.InsertOneResult, err error) {

	defer func() {
		if r := recover(); r != nil {
			log.Error("panic -  \n", r)
			// log.Fatal(err)
			err = r.(error)
		}
	}()

	result, err = a.client.Database(argDb).Collection(argColl).InsertOne(context.TODO(), *leaderboard)
	if err != nil {
		log.Error("insert -  \n", err)
		// log.Fatal(err)
		return nil, err
	}
	return result, nil
}

// FindLeaderboardByID returns an leaderboard matching the id provided, if any.
func (a *AtlasClientService) FindLeaderboardByID(id, argDb, argColl string) (result *LeaderboardT, err error) {

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

// FindLeaderboardsByGameID returns an leaderboard matching the id provided, if any.
func (a *AtlasClientService) FindLeaderboardsByGameID(id, argDb, argColl string) (results []*LeaderboardT, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("panic -  \n", r)
			err = r.(error)
		}
	}()

	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
	// TODO - don't hack this aggregation srly make it more generic..
	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
	// adding the $lookup on the materialiezed leaderboards - which are generated by the schedueld trigger
	filter := bson.D{primitive.E{Key: "gameId", Value: id}}
	match := bson.D{{"$match", filter}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "leaderboards_mv"},
		{"localField", "string"},
		{"foreignField", "string"},
		{"as", "data"}}}}

	pipeline := bson.A{match, lookupStage}
	cur, err := a.client.Database(argDb).Collection(argColl).Aggregate(context.TODO(), pipeline)
	// cur, err := a.client.Database(argDb).Collection(argColl).Find(context.TODO(), filter)
	if err != nil {
		return
	}
	defer cur.Close(context.TODO())
	err = cur.All(context.TODO(), &results)
	return

}

// UpdateLeaderboardByID attempts to update the matching leaderboard with the provided update
func (a *AtlasClientService) UpdateLeaderboardByID(id, argDb, argColl string, leaderboard *LeaderboardT) (result *mongo.UpdateResult, err error) {

	defer func() {
		if r := recover(); r != nil {
			log.Error("panic -  \n", r)
			err = r.(error)
		}
	}()

	filter := bson.D{primitive.E{Key: "id", Value: id}}
	update := bson.D{primitive.E{Key: "$set", Value: *leaderboard}}
	result, err = a.client.Database(argDb).Collection(argColl).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error("insert -  \n", err)
		return nil, err
	}
	return result, nil
}

// DeleteLeaderboard is
func (a *AtlasClientService) DeleteLeaderboard(id, argDb, argColl string) (result *mongo.DeleteResult, err error) {

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

// CountLeaderboard returns the count with a given filter
func (a *AtlasClientService) CountLeaderboard(filters map[string]string, argDb, argColl string) (result int64, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("panic -  \n", r)
			err = r.(error)
		}
	}()

	log.Error("got filters as \n", filters)
	if filters == nil || len(filters) == 0 {
		result, err = a.client.Database(argDb).Collection(argColl).EstimatedDocumentCount(context.TODO())
	} else {
		result, err = a.client.Database(argDb).Collection(argColl).CountDocuments(context.TODO(), filters)
	}

	return

}
