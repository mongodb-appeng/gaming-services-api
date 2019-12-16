package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	api "github.com/mongodb-appeng/gaming-services-api/internal/pkg/api"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

//ServerOptionsT is
type ServerOptionsT struct {
	Port    string `yaml:"port"`
	Address string `yaml:"address"`
	LogFile     string `yaml:"logFile"`  // no file == stdout, else output to provided file
	LogLevel    int    `yaml:"logLevel"` // verbosity level control
	DatabaseURI string `yaml:"databaseURI"`
}

var (
	buildRelease    = "undefined"      // Release Version Tracking, run binary with -v flag. Set at compile time by Makefile.
	gitVersion      = "undefined"      // Release Version Tracking, run binary with -v flag. Set at compile time by Makefile.
	operatingSystem = "undefined"      // Release Version Tracking, run binary with -v flag. Set at compile time by Makefile.
	options         = ServerOptionsT{} // parsed YAML file config options stored here
)

func main() {

	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
	// TODO - load the configs/database.yaml intead of setting env var here
	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO

	// Only log the User-specified severity or above.
	log.SetLevel(log.DebugLevel)

	yamlFile := flag.String("c", "", "<path to .yaml configuration file>")
	version := flag.Bool("v", false, "prints the build version and exits ")
	flag.Parse()

	if *version {
		fmt.Printf("Build release %s. Git version %s compiled %s. \n", buildRelease, gitVersion, operatingSystem)
		os.Exit(0)
	}
	if yamlFile == nil || *yamlFile == "" {
		panic(fmt.Errorf("missing required yaml config file"))
	}
	y, err := ioutil.ReadFile(*yamlFile)
	if err != nil {
		panic(fmt.Errorf("cant parse yaml"))
	}

	err = yaml.Unmarshal(y, &options)
	if err != nil {
		panic(fmt.Errorf("cant unmarshal yaml"))
	}
	if options.LogFile != "" {
		f, err := os.OpenFile(options.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		if err != nil {
			panic(err)
		}
		// util.Log.UseFileLogger(f) - TODO change to logrus file
	}

	api.SetUpBackend(options.DatabaseURI)

	log.Debug("Adding routes\n")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", api.RootHandler)

	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
	// TODO - move this to internal/pkg/rest/routes.go
	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
	router.HandleFunc("/v1/account", api.CreateAccountHandler).Methods("POST")

	router.HandleFunc("/v1/gameevent/", api.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/v1/gameevent/", api.PostGameEventsHandler).Methods("PUT")

	router.HandleFunc("/v1/account/{id}", api.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/v1/account/{id}", api.ReadAccountHandler).Methods("GET")
	router.HandleFunc("/v1/account/{id}", api.UpdateAccountHandler).Methods("PATCH")
	router.HandleFunc("/v1/account/{id}", api.DeleteAccountHandler).Methods("DELETE")
	router.HandleFunc("/v1/account/{id}", api.NewStitchLoginHandler).Methods("PUT") // called only when user logs in via Stitch

	router.HandleFunc("/v1/gamerhandle/", api.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/v1/gamerhandle/", api.ReadRandomGamerHandleHandler).Methods("GET")
	router.HandleFunc("/v1/gamerprofile", api.CreateGamerProfileHandler).Methods("POST")

	router.HandleFunc("/v1/gamerprofile/{id}", api.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/v1/gamerprofile/{id}", api.ReadGamerProfileHandler).Methods("GET")
	router.HandleFunc("/v1/gamerprofile/{id}", api.UpdateGamerProfileHandler).Methods("PATCH")
	router.HandleFunc("/v1/gamerprofile/{id}", api.DeleteGamerProfileHandler).Methods("DELETE")

	router.HandleFunc("/v1/leaderboard", api.CreateLeaderboardHandler).Methods("POST")
	router.HandleFunc("/v1/leaderboard/{id}", api.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/v1/leaderboard/{id}", api.ReadLeaderboardHandler).Methods("GET")
	router.HandleFunc("/v1/leaderboard/{id}", api.UpdateLeaderboardHandler).Methods("PATCH")
	router.HandleFunc("/v1/leaderboard/{id}", api.DeleteLeaderboardHandler).Methods("DELETE")

	router.HandleFunc("/v1/leaderboard/findByGameID/{id}", api.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/v1/leaderboard/findByGameID/{id}", api.ReadLeaderboardsByGameIDHandler).Methods("GET")
	router.HandleFunc("/v1/leaderboard/findByGameID/{id}", api.UpdatePlayedGameHandler).Methods("PATCH")

	listenOn := ":8888" //options.Address +
	log.Debug("Starting server on " + listenOn)
	if err := http.ListenAndServe(listenOn, router); err != nil {
		panic(err)
	}
}
