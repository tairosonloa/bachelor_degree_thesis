package api

import (
	"encoding/json"
	"log"
	"net/http"

	// "rpi3/API_REST/app/controllers"
	"rpi3/API_REST/app/models"
)

var (
	authorized string // Bearer token for authorized POST
)

// Initialize initializes the API server handlers and inner state
func Initialize(apiAuthorizedToken string) *http.ServeMux {
	// Handlers
	mux := http.NewServeMux()

	mux.HandleFunc("/", index)
	mux.HandleFunc("/reservation", reservation)
	mux.HandleFunc("/reservation/{[0-9]+}", reservationID)

	mux.HandleFunc("/favicon.ico", func(_ http.ResponseWriter, _ *http.Request) {})

	// Inner state
	// TODO: cargar sqlite

	return mux
}

// index is the API server handler for "/"
// If method is GET, it responds with a JSON containing the message "API is up and running"
// else, it responds with a JSON containing an error message
func index(w http.ResponseWriter, r *http.Request) {
	// Only method GET is allowed
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		log.Printf("%s / from %s status %d\n", r.Method, r.RemoteAddr, http.StatusMethodNotAllowed)
		return
	}
	// Responds with API status
	response := models.MessageAPIM{
		Message: "API is up and running",
	}
	respondWithJSON(w, http.StatusOK, response)
	log.Printf("%s / from %s status %d\n", r.Method, r.RemoteAddr, http.StatusOK)
}

// reservation is the API server handler for "/reservation"
// If method is GET, it responds with a JSON containing all reservations info
// If method is POST and authentication is susccessful, it creates a reservation values and info
// else, it responds with a JSON containing an error message
func reservation(w http.ResponseWriter, r *http.Request) {
	// Check http method
	if r.Method == http.MethodPost {
		// Create new reservation
		createReservation(w, r)
	} else if r.Method == http.MethodGet {
		// Get all reservations
		getReservations(w, r)
	} else {
		respondWithError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		log.Printf("%s /reservation from %s status %d\n", r.Method, r.RemoteAddr, http.StatusMethodNotAllowed)
		return
	}
}

// reservationID is the API server handler for "/reservation/<id>"
// If method is GET, it responds with a JSON containing the reservation with id <id> info
// If method is PUT and authentication is susccessful, it updates a reservation values and info
// else, it responds with a JSON containing an error message
func reservationID(w http.ResponseWriter, r *http.Request) {
	// Check http method
	if r.Method == http.MethodGet {
		// Get reservation info
		getReservationByID(w, r)
	} else if r.Method == http.MethodPut {
		// Update reservation info
		updateReservationByID(w, r)
	} else {
		respondWithError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		log.Printf("%s /reservation/<id> from %s status %d\n", r.Method, r.RemoteAddr, http.StatusMethodNotAllowed)
		return
	}
}

// createReservation creates a new reservation from a JSON POST request
func createReservation(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	if !validateToken(w, r) {
		log.Printf("%s /reservation from %s status %d\n", r.Method, r.RemoteAddr, http.StatusUnauthorized)
		return
	}

	// Read JSON from body request and create reservation
	reservation := models.Reservation{}
	decoder := json.NewDecoder(r.Body)
	e := decoder.Decode(&reservation)
	if e != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		log.Printf("%s /reservation %s status %d\n", r.Method, r.RemoteAddr, http.StatusInternalServerError)
		log.Println(e.Error())
		return
	}
}

// getReservations responds with a JSON cantaining all reservations
func getReservations(w http.ResponseWriter, r *http.Request) {
	// TODO: filter by url params
	// TODO: to do
}

// getReservationByID responds with a JSON cantaining a reservation info
func getReservationByID(w http.ResponseWriter, r *http.Request) {
	// TODO: to do
}

// updateReservationByID updates a existing reservation from a JSON PUT request
func updateReservationByID(w http.ResponseWriter, r *http.Request) {
	// TODO: to do
}

// validateToken checks if the request is authenticated (bearer token) and authorized
func validateToken(w http.ResponseWriter, r *http.Request) bool {
	token := r.Header.Get("Authorization")
	if token != authorized {
		var err string
		if token == "" {
			err = "Authorization header not provided or empty"
		} else {
			err = http.StatusText(http.StatusUnauthorized)
		}
		respondWithError(w, http.StatusUnauthorized, err)
		return false
	}
	return true
}

// respondWithError responds to a request with a http code and a JSON containing an error message
func respondWithError(w http.ResponseWriter, code int, message string) {
	err := models.ErrorAPIM{
		Error: message,
	}
	respondWithJSON(w, code, err)
}

// respondWithJSON responds to a request with a http code and a JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
