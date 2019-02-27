package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("GUI/public"))
	http.Handle("/", fs)

	port := flag.String("port", "9000", "Port where web server will be listening")
	flag.Parse()

	log.Printf("Web server listening on port %v...\n", string(*port))
	http.ListenAndServe(":"+string(*port), nil)
}
