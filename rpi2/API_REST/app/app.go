package app

import (
	"app/api"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

// App represents the core of the application (server and API)
type App struct {
	Addr      string
	Port      int
	server    *http.Server
	handlers  *http.ServeMux
	tokenFile string
}

// readCmd reads command line arguments (address, port, and token file)
// Default is address localhost, port 3000 and token file "token.txt" in current directory
// Run the binary with --help or -h for more info
func (a *App) readCmd() {
	addr := flag.String("addr", "localhost", "Address where API will be listening")
	port := flag.Int("port", 3000, "Port where APi will be listening")
	currentDir, _ := os.Getwd()
	tokenFile := flag.String("tokenf", currentDir+"/token.txt", "Path to the file containing authorized bearer token")
	flag.Parse()
	a.Addr = string(*addr)
	a.Port = int(*port)
	a.tokenFile = string(*tokenFile)
}

// Initialize initializes the API server address, port and hadlers
func (a *App) Initialize() {
	a.readCmd()

	log.Println("Initializating server")
	a.handlers = api.Initialize(a.tokenFile)
	a.server = &http.Server{Handler: a.handlers, Addr: fmt.Sprintf("%s:%d", a.Addr, a.Port)}
}

// Run runs the API server
func (a *App) Run() {
	log.Printf("Now listening on %s:%d\n", a.Addr, a.Port)
	log.Fatal(a.server.ListenAndServe())
}
