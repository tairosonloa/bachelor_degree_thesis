package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", "9000", "Port where web server will be listening")
	root := flag.String("root", "/srv/rpi3", "Web app static files root path")
	flag.Parse()

	fs := http.FileServer(http.Dir(string(*root)))
	http.Handle("/", fs)

	log.Printf("Web server listening on port %v...\n", string(*port))
	http.ListenAndServe(":"+string(*port), nil)
}
