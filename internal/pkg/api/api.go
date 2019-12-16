package api

import (
	"net/http"
	log "github.com/sirupsen/logrus"
	"github.com/mongodb-appeng/gaming-services-api/internal/pkg/database"
)


var serverSettings struct {
	atlas *database.AtlasClientService
}
// Atlas database is a free tier MongoDB deployment
var atlas *database.AtlasClientService

// SetUpBackend establishes a connection to atlas
func SetUpBackend(uri string) {
	log.Debug("Setting up the backend")
	atlas = database.NewAtlasClientService(uri) // change to https://golang.org/pkg/os/#LookupEnv
	atlas.Connect()
	atlas.Ping()
}

// Generic/Shared across verions handlers

// RootHandler is
func RootHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Game Platform Services API. Latest is /v1")
}

// OptionsHandler for CORS bullshit
func OptionsHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	// w.Header().Set("Access-Control-Expose-Headers", "Access-Control-*")
	w.Header().Set("Access-Control-Allow-Headers", "Access-Control-*,  Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	// w.Header().Set("Allow", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
	// w.Header().Set("Vary", "Origin")
	w.WriteHeader(http.StatusOK)
	return

}
