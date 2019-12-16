package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/mongodb-appeng/gaming-services-api/internal/pkg/database"
	"github.com/gorilla/mux"
)

// HTTP Handlers CRUD+ for Leaderboard.

// CreateLeaderboardHandler is
func CreateLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r != nil {
			r.Body.Close()
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("body - ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var doc database.LeaderboardT
	err = json.Unmarshal(body, &doc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := atlas.CreateLeaderboard("gamePlatformServices", "leaderboards", &doc)
	if err != nil {
		log.Error("atlas - ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// success response
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Vary", "Origin")

	// TOD - use  render.JSON ???
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Location", r.URL.String()+"/"+doc.ID) // relative ref for Created -- TODO test this works
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&result)
}

//ReadLeaderboardsByGameIDHandler is
func ReadLeaderboardsByGameIDHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r != nil {
			r.Body.Close()
		}
	}()

	params := mux.Vars(r)
	log.Debug(params)
	ID, ok := params["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := atlas.FindLeaderboardsByGameID(ID, "gamePlatformServices", "leaderboards")
	if err != nil {
		log.Error("atlas - ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} // success response
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Vary", "Origin")
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(&result)
}

//ReadLeaderboardHandler is
func ReadLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r != nil {
			r.Body.Close()
		}
	}()

	params := mux.Vars(r)
	log.Debug(params)
	ID, ok := params["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := atlas.FindLeaderboardByID(ID, "gamePlatformServices", "leaderboards")
	if err != nil {
		log.Error("atlas - ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} // success response
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Vary", "Origin")
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(&result)
}

//CountLeaderboardsHandler is
func CountLeaderboardsHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r != nil {
			r.Body.Close()
		}
	}()

	result, err := atlas.CountLeaderboard(mux.Vars(r), "gamePlatformServices", "leaderboards")
	if err != nil {
		log.Error("atlas - ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// success response
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Vary", "Origin")
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(&result)
}

//UpdateLeaderboardHandler is
func UpdateLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r != nil {
			r.Body.Close()
		}
	}()

	params := mux.Vars(r)
	ID, ok := params["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// for now accept LeaderboardT, TODO - limit/control fields
	var doc database.LeaderboardT
	err = json.Unmarshal(body, &doc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := atlas.UpdateLeaderboardByID(ID, "gamePlatformServices", "leaderboards", &doc)
	if err != nil {
		log.Error("atlas - ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// success
	// success response
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Vary", "Origin")
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK) // For a PUT request: HTTP 200 or HTTP 204 should imply "resource updated successfully".
	json.NewEncoder(w).Encode(&result)
}

//DeleteLeaderboardHandler is
func DeleteLeaderboardHandler(w http.ResponseWriter, r *http.Request) {

	log.Debug("DeleteLeaderboardHandler call ", r)
	defer func() {
		if r != nil {
			r.Body.Close()
		}
	}()

	params := mux.Vars(r)
	ID, ok := params["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := atlas.DeleteLeaderboard(ID, "gamePlatformServices", "leaderboards")
	if err != nil {
		log.Error("atlas - ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// DELETE has been enacted && the response includes an entity describing the status
	// success response
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
	w.Header().Set("Vary", "Origin")
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&result)

}

//RunAggregateLeaderboardHandler is
func RunAggregateLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
