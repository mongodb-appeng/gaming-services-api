package api

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"

	database "github.com/desteves/babysteps/internal/pkg/database"
)

// HTTP Handlers CRUD+ for Leaderboard.

// PostGameEventsHandler is
func PostGameEventsHandler(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if r != nil {
			r.Body.Close()
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("body - %+v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO -- support array
	var doc database.GameEventT
	err = json.Unmarshal(body, &doc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	doc.CreatedAt = time.Now()
	result, err := atlas.AddGameEvents("gamePlatformServices", "gameevents", &doc)
	if err != nil {
		log.Error("atlas - %+v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// success response
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Vary", "Origin")
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(&result)
}
