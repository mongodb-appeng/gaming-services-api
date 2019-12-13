package database

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AuthIdentityT is
type AuthIdentityT struct {
	ID           string `json:"id" bson:"id"`
	ProviderType string `json:"provider_type" bson:"provider_type"`
}

// AuthDataT is
type AuthDataT struct {
	Name      string `json:"name" bson:"name"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Picture   string `json:"picture" bson:"picture"`
	Email     string `json:"email" bson:"email"`
}

// AuthProfileT is
type AuthProfileT struct {
	Data *AuthDataT       `json:"data" bson:"data"`
	Type string           `json:"type" bson:"type"`
	IDs  []*AuthIdentityT `json:"identities" bson:"identities"`
}

//AuthProviderT is data provided by a 3rd party auth service
type AuthProviderT struct {
	ID                   string        `json:"id" bson:"id"`
	LoggedInProviderType string        `json:"loggedInProviderType" bson:"loggedInProviderType"`
	LoggedInProviderName string        `json:"loggedInProviderName" bson:"loggedInProviderName"`
	AuthProfile          *AuthProfileT `json:"profile" bson:"profile"`
	IsLoggedIn           bool          `json:"isLoggedIn" bson:"isLoggedIn"`
	LastAuthActivity     string        `json:"lastAuthActivity" bson:"lastAuthActivity"`
}

// AffiliationT is
type AffiliationT struct {
	Environment string `json:"environment" bson:"environment"` // school, work, etc
	Title       string `json:"title" bson:"title"`
	Place       string `json:"place" bson:"place"`
	Function    string `json:"function" bson:"function"`
}

// SocialT is
type SocialT struct {
	Service string `json:"service" bson:"service"`
	Handle  string `json:"handle" bson:"handle"`
}

// AccountT is
type AccountT struct {
	ObjectID     *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ID           string              `json:"id,omitempty" bson:"id,omitempty"`
	Name         string              `json:"name,omitempty" bson:"name,omitempty"`
	AuthProvider *AuthProviderT      `json:"auth,omitempty" bson:"auth,omitempty"`
	Social       []*SocialT          `json:"social,omitempty" bson:"social,omitempty"`
	Affiliation  []*AffiliationT     `json:"affiliation,omitempty" bson:"affiliation,omitempty"`
	Location     string              `json:"location,omitempty" bson:"location,omitempty"`
	Meta         *MetaT              `json:"meta,omitempty" bson:"meta,omitempty"`
}

// CreateAccount is
func (a *AtlasClientService) CreateAccount(argDb string, argColl string, account *AccountT) (result *mongo.InsertOneResult, err error) {

	defer func() {
		if r := recover(); r != nil {
			log.Error("panic -  %+v\n", r)
			// log.Fatal(err)
			err = r.(error)
		}
	}()

	result, err = a.client.Database(argDb).Collection(argColl).InsertOne(context.TODO(), *account)
	if err != nil {
		log.Error("insert -  %+v\n", err)
		// log.Fatal(err)
		return nil, err
	}
	return result, nil
}

// FindAccountByID returns an account matching the id provided, if any.
func (a *AtlasClientService) FindAccountByID(id, argDb, argColl string) (result *AccountT, err error) {

	defer func() {
		if r := recover(); r != nil {
			log.Error("panic -  %+v\n", r)
			// log.Fatal(err)
			err = r.(error)
		}
	}()

	filter := bson.D{primitive.E{Key: "id", Value: id}}
	err = a.client.Database(argDb).Collection(argColl).FindOne(context.TODO(), filter).Decode(&result)
	return
}

// UpdateAccountByID attempts to update the matching account with the provided update
func (a *AtlasClientService) UpdateAccountByID(id, argDb, argColl string, account *AccountT) (result *mongo.UpdateResult, err error) {

	defer func() {
		if r := recover(); r != nil {
			log.Error("panic -  %+v\n", r)
			err = r.(error)
		}
	}()

	log.Debug("UpdateAccountByIDUpdateAccountByIDUpdateAccountByID  %+v\n", account)

	filter := bson.D{primitive.E{Key: "id", Value: id}}
	update := bson.D{primitive.E{Key: "$set", Value: *account}}
	result, err = a.client.Database(argDb).Collection(argColl).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error("insert -  %+v\n", err)
		return nil, err
	}
	return result, nil
}

// NewStitchLogin attempts to update the matching account with the provided update
func (a *AtlasClientService) NewStitchLogin(id, argDb, argColl string, auth *AuthProviderT) (result *mongo.UpdateResult, err error) {

	defer func() {
		if r := recover(); r != nil {
			log.Error("panic -  %+v\n", r)
			err = r.(error)
		}
	}()

	var doc AccountT
	t := time.Now()
	var meta MetaT
	// var social SocialT
	// social.Email = "missing"
	// social.Phone = "000-000-0000"
	// social.Instagram = "missing"
	// social.Twitter = "missing"
	meta.Created = t
	meta.Updated = t
	oid, _ := primitive.ObjectIDFromHex(id)
	doc.ID = id
	doc.Name = "missing"
	doc.ObjectID = &oid
	doc.AuthProvider = auth
	doc.Location = "US"
	// doc.Social = &social
	doc.Meta = &meta

	log.Debug("new stitch log attempting to insert  %+v\n", doc)
	inserted, err := a.CreateAccount(argDb, argColl, &doc)
	//first time user
	if err != nil {
		// duplicate key most likely, attempt to do update instead
		filter := bson.D{primitive.E{Key: "id", Value: id}}
		update := bson.D{primitive.E{Key: "$set", Value: primitive.E{Key: "auth", Value: *auth}}}
		// opts := options.Update().SetUpsert(true)
		result, err = a.client.Database(argDb).Collection(argColl).UpdateOne(context.TODO(), filter, update)

		log.Debug("new stitch log in update result is %+v, %+v\n", result, err)

	} else {
		var gamer GamerProfileT
		gamer.ObjectID = &oid                                     // use the same
		gamer.ID = inserted.InsertedID.(primitive.ObjectID).Hex() // _id
		gamer.Location = doc.Location
		gamer.Meta = doc.Meta
		gamer.Handle = "playerhandle"
		var image ImageT
		image.Avatar = AvatarILoveMongo
		image.Banner = BannerWarioWares

		log.Debug("Now adding a gamer profile for the first time log in user...%+v \n", gamer)
		gamer.Image = &image
		creategamer, err := a.CreateGamerProfile("gamePlatformServices", "gamerprofiles", &gamer)
		if err != nil {
			log.Error("gamer err -  %+v\n", err)
		} else {
			log.Debug("inserted gamer profile result is %+v\n", creategamer)
		}

	}
	return
}

// account.ID = id // do not allow the user to update this value!
// filter := bson.D{primitive.E{Key: "id", Value: id}}
// update := bson.D{primitive.E{Key: "$set", Value: *account}}
// opts := options.Update().SetUpsert(true)
// result, err = a.client.Database(argDb).Collection(argColl).UpdateOne(context.TODO(), filter, update, opts)
// if err != nil {
// 	log.Error("upsert -  %+v\n", err)
// 	return
// }
// // auto generate a gamer profile to associate with the new account - game logic
// if result.UpsertedCount == 1 {
// 	var gamer GamerProfileT
// 	gamer.ID = result.UpsertedID.(primitive.ObjectID).Hex() // _id
// 	gamer.Location = account.Location
// 	// log.Debug("attempting to insert gamer profile for account  %+v with %+v", account.ID, gamer.ID)
// 	_, err = a.CreateGamerProfile("gamePlatformServices", "gamerprofiles", &gamer)
// 	if err != nil {
// 		log.Error("upsert -  %+v\n", err)
// 		return
// 	}
// }
// return

// DeleteAccount is
func (a *AtlasClientService) DeleteAccount(id, argDb, argColl string) (result *mongo.DeleteResult, err error) {

	defer func() {
		if r := recover(); r != nil {
			log.Error("panic -  %+v\n", r)
			err = r.(error)
		}
	}()

	filter := bson.D{primitive.E{Key: "id", Value: id}}
	result, err = a.client.Database(argDb).Collection(argColl).DeleteOne(context.TODO(), filter)
	return
}
