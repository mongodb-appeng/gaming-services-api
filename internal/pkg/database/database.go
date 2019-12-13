package database

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MetaT is
type MetaT struct {
	Created time.Time `json:"created" bson:"created"`
	Updated time.Time `json:"updated" bson:"updated"`
	Version int32     `json:"version" bson:"version"`
	Closed  *bool     `json:"closed" bson:"closed"`
	Deleted *bool     `json:"deleted" bson:"deleted"`
}

// AtlasClientService struct
type AtlasClientService struct {
	// goctx       context.Context
	// goctxCancel context.CancelFunc
	uri    string
	client *mongo.Client
}

// NewAtlasClientService is
func NewAtlasClientService(uri string) *AtlasClientService {
	atlas := new(AtlasClientService)
	// atlas.goctx, atlas.goctxCancel = context.WithTimeout(context.Background(), 10*time.Second)
	atlas.uri = uri
	return atlas
}

// Connect to atlas mongo deployment
func (a *AtlasClientService) Connect() error {
	// https://docs.mongodb.com/ecosystem/drivers/go/#connect-to-mongodb-atlas
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(a.uri))
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Debug("Got a Mongo Client!")
	a.client = client
	return nil
}

// Ping the atlas mongo deployment, mostly for testing
func (a *AtlasClientService) Ping() error {
	err := a.client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Debug("Pinged Atlas!")
	return nil

}
