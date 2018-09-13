package app

import (
	"app/api"
	"flag"
	"fmt"
	"log"
	"net/http"
)

// App represents the core of the application (server and API)
type App struct {
	Addr     string
	Port     int
	server   *http.Server
	handlers *http.ServeMux
}

// readCmd reads command line arguments (addr for address and port for port)
// Default is address localhost and port 3000
func (a *App) readCmd() {
	addrPtr := flag.String("addr", "localhost", "a string")
	portPtr := flag.Int("port", 3000, "a int")
	flag.Parse()
	a.Addr = string(*addrPtr)
	a.Port = int(*portPtr)
}

// Initialize initializes the API server address, port and hadlers
func (a *App) Initialize() {
	a.readCmd()

	log.Println("Initializating server")
	a.handlers = api.Handlers()
	a.server = &http.Server{Handler: a.handlers, Addr: fmt.Sprintf("%s:%d", a.Addr, a.Port)}
}

// Run runs the API server
func (a *App) Run() {
	log.Printf("Now listening on %s:%d\n", a.Addr, a.Port)
	log.Fatal(a.server.ListenAndServe())
}
