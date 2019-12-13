package api

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"

	database "github.com/desteves/babysteps/internal/pkg/database"
	"github.com/gorilla/mux"
)

// HTTP Handlers CRUD+ for Account.

// CreateAccountHandler is
func CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
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

	// for now accept AccountT, TODO - limit/control fields
	var doc database.AccountT
	err = json.Unmarshal(body, &doc)
	if err != nil {
		log.Error("unmarshal - %+v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := atlas.CreateAccount("gamePlatformServices", "accounts", &doc)
	if err != nil {
		log.Error("atlas - %+v", err.Error())
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

//ReadAccountHandler is
func ReadAccountHandler(w http.ResponseWriter, r *http.Request) {
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
	result, err := atlas.FindAccountByID(ID, "gamePlatformServices", "accounts")
	if err != nil {
		log.Error("atlas - %+v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// success response
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Vary", "Origin")
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(&result)
}

//UpdateAccountHandler is
func UpdateAccountHandler(w http.ResponseWriter, r *http.Request) {
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
	var doc database.AccountT
	err = json.Unmarshal(body, &doc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := atlas.UpdateAccountByID(ID, "gamePlatformServices", "accounts", &doc)
	if err != nil {
		log.Error("atlas - %+v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// success response
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Vary", "Origin")
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&result)
}

//NewStitchLoginHandler is
func NewStitchLoginHandler(w http.ResponseWriter, r *http.Request) {
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
	// log.Debug("NewStitchLoginHandler ID %+v \n ", ID)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// var temp map[string]interface{}
	var doc database.AuthProviderT
	err = json.Unmarshal(body, &doc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Debug("NewStitchLoginHandler doc %+v \n ", doc)
	// delete(temp, "_id")
	// var doc database.AccountT
	// mapstructure.Decode(temp, &doc)

	result, err := atlas.NewStitchLogin(ID, "gamePlatformServices", "accounts", &doc)

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

//DeleteAccountHandler is
func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
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
	result, err := atlas.DeleteAccount(ID, "gamePlatformServices", "accounts")
	if err != nil {
		log.Error("atlas - %+v", err.Error())
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

//RunAggregateAccountHandler is
func RunAggregateAccountHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
