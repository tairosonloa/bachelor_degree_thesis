package app

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"rpi3/API_REST/app/api"
)

// configValues represents config values readed from JSON on initialization
type configValues struct {
	Rpi3APIAddress string
	Rpi3APIPort    int
	ControlServer  string
	OccupationCmd  string
	OccupationWeb  string
}

// App represents the core of the application (server and API)
type App struct {
	server     *http.Server
	handlers   *http.ServeMux
	configFile string
	config     configValues
}

// readCmd reads command line arguments (config file path)
// Default is file "config.json" in current directory
// Run the binary with --help or -h for more info
func (a *App) readCmd() {
	currentDir, _ := os.Getwd()
	configFile := flag.String("conf", currentDir+"/config.json", "Path to the file config.json")
	flag.Parse()
	a.configFile = string(*configFile)
}

// loadConfig loads the config file
func (a *App) loadConfig() {
	a.config = configValues{}
	fd, err := os.Open(a.configFile)
	defer fd.Close()
	if err != nil {
		log.Printf("ERROR app/loadConfig(): %v\nExiting...", err.Error())
		os.Exit(1)
	} else {
		decoder := json.NewDecoder(fd)
		err = decoder.Decode(&a.config)
		if err != nil {
			log.Printf("ERROR app/loadConfig(): %v\n", err.Error())
		}
	}
	log.Println(a.config)
}

// Initialize initializes the API server address, port and hadlers
func (a *App) Initialize() {
	a.readCmd()
	a.loadConfig()
	log.Println("Initializating server")
	a.handlers = api.Initialize(a.config.ControlServer, a.config.OccupationCmd, a.config.OccupationWeb)
	a.server = &http.Server{Handler: a.handlers, Addr: fmt.Sprintf("%s:%d", a.config.Rpi3APIAddress, a.config.Rpi3APIPort)}
}

// Run runs the API server
func (a *App) Run() {
	log.Printf("Now listening on %s:%d\n", a.config.Rpi3APIAddress, a.config.Rpi3APIPort)
	log.Fatal(a.server.ListenAndServe())
}
