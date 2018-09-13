package api

import (
	"app/models"
	"encoding/json"
	"log"
	"net/http"
)

// Handlers initializes the API server handlers
func Handlers() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", index)

	mux.HandleFunc("/favicon.ico", func(_ http.ResponseWriter, _ *http.Request) {})

	return mux
}

// index is the API server handler for "/"
func index(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		err := models.ErrorAPI{
			Error: http.StatusText(http.StatusMethodNotAllowed),
		}
		respondWithJSON(w, http.StatusMethodNotAllowed, err)
		log.Printf("%s / %d\n", r.Method, http.StatusMethodNotAllowed)
		return
	}
	// Responds with API status
	response := models.MessageAPI{
		Message: "API is up and running",
	}
	respondWithJSON(w, http.StatusOK, response)
	log.Printf("%s / %d\n", r.Method, http.StatusOK)
}

// respondWithError responds to a request with a http code and a JSON with an error message
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON responds to a request with a http code and a JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
