package api

import (
	"encoding/json"
	
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"

	database "github.com/mongodb-appeng/gaming-services-api/internal/pkg/database"
	"github.com/gorilla/mux"
)

// HTTP Handlers CRUD+ for GamerProfile.

//ReadRandomGamerHandleHandler is
func ReadRandomGamerHandleHandler(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if r != nil {
			r.Body.Close()
		}
	}()

	result, err := atlas.GetRandomGamerHandle("gamePlatformServices", "gamerhandles")
	if err != nil {
		log.Error("atlas - \n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// success response
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(&result)

}

// CreateGamerProfileHandler is
func CreateGamerProfileHandler(w http.ResponseWriter, r *http.Request) {
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

	// for now accept GamerProfileT, TODO - limit/control fields
	var doc database.GamerProfileT
	err = json.Unmarshal(body, &doc)
	if err != nil {
		log.Error("unmarshal - ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := atlas.CreateGamerProfile("gamePlatformServices", "gamerprofiles", &doc)
	if err != nil {
		log.Error("atlas - ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// success response
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Vary", "Origin")
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Location", r.URL.String()+"/"+doc.ID) // relative ref for Created -- TODO test this works
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&result)
}

//ReadGamerProfileHandler is
func ReadGamerProfileHandler(w http.ResponseWriter, r *http.Request) {
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
	result, err := atlas.FindGamerProfileByID(ID, "gamePlatformServices", "gamerprofiles")
	if err != nil {
		log.Error("atlas - \n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// success response
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	// w.Header().Set("Vary", "Origin")
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(&result)
}

//ReadGamerProfileByAccountIDHandler is
func ReadGamerProfileByAccountIDHandler(w http.ResponseWriter, r *http.Request) {
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
	result, err := atlas.FindGamerProfileByAccountID(ID, "gamePlatformServices", "gamerprofiles")
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

//UpdateGamerProfileHandler is
func UpdateGamerProfileHandler(w http.ResponseWriter, r *http.Request) {
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

	// for now accept AccountT, TODO - limit/control fields
	var temp map[string]interface{}
	err = json.Unmarshal(body, &temp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// var doc database.GamerProfileT
	// err = json.Unmarshal(body, &doc)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	result, err := atlas.UpdateGamerProfileByID(ID, "gamePlatformServices", "gamerprofiles", temp)
	if err != nil {
		log.Error("atlas - ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// success response
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Vary", "Origin")
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK) // For a PUT request: HTTP 200 or HTTP 204 should imply "resource updated successfully".
	json.NewEncoder(w).Encode(&result)
}

//UpdatePlayedGameHandler is
func UpdatePlayedGameHandler(w http.ResponseWriter, r *http.Request) {
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

	gameID, ok := params["game"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var doc database.PlayedGameT
	err = json.Unmarshal(body, &doc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	doc.ID = ""
	doc.Title = ""
	result, err := atlas.UpdatePlayedGame(ID, gameID, "gamePlatformServices", "gamerprofiles", &doc)
	if err != nil {
		log.Error("atlas - ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// success response
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Vary", "Origin")
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK) // For a PUT request: HTTP 200 or HTTP 204 should imply "resource updated successfully".
	json.NewEncoder(w).Encode(&result)
}

//DeleteGamerProfileHandler is
func DeleteGamerProfileHandler(w http.ResponseWriter, r *http.Request) {
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
	result, err := atlas.DeleteGamerProfile(ID, "gamePlatformServices", "gamerprofiles")
	if err != nil {
		log.Error("atlas - ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// DELETE has been enacted && the response includes an entity describing the status
	// success response
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Vary", "Origin")
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&result)

}

//RunAggregateGamerProfileHandler is
func RunAggregateGamerProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
