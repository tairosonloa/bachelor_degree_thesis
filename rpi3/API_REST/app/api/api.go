package api

import (
	"encoding/json"
	"log"
	"net/http"

	"rpi3/API_REST/app/controllers"
	"rpi3/API_REST/app/models"
)

var (
	controlServer string
	occupationCmd string
	occupationWeb string
)

// Initialize initializes the API server handlers and inner state
func Initialize(server, cmd, web string) *http.ServeMux {

	// Handlers
	mux := http.NewServeMux()

	mux.HandleFunc("/", index)
	mux.HandleFunc("/reservations", reservations)
	mux.HandleFunc("/classrooms", classrooms)
	mux.HandleFunc("/occupation", occupation)

	mux.HandleFunc("/favicon.ico", func(_ http.ResponseWriter, _ *http.Request) {})

	// Inner state
	controlServer = server
	occupationCmd = cmd
	occupationWeb = web

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

// reservations is the API server handler for "/reservations"
// If method is GET, it responds with a JSON containing all reservations info
// else, it responds with a JSON containing an error message
func reservations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// Check http method
	if r.Method == http.MethodGet {
		getReservations(w, r)
	} else {
		respondWithError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		log.Printf("%s /reservation from %s status %d\n", r.Method, r.RemoteAddr, http.StatusMethodNotAllowed)
		return
	}
}

// getReservations responds with a JSON cantaining all reservations
func getReservations(w http.ResponseWriter, r *http.Request) {
	// TODO: filter by url params
	reservations := controllers.GetTodayReservations(occupationWeb)
	if reservations == nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		log.Printf("%s /reservation %s status %d\n", r.Method, r.RemoteAddr, http.StatusInternalServerError)
		log.Println("Reservations slice is nil.")
		return
	}
	log.Printf("%s /reservation %s status %d\n", r.Method, r.RemoteAddr, http.StatusOK)
	respondWithJSON(w, http.StatusOK, reservations)
}

// classrooms is the API server handler for "/classrooms"
// If method is GET, it responds with a JSON containing all classrooms status
// else, it responds with a JSON containing an error message
// Valid status are: 0 - free, 1 - occupied, 2 - reserved within next two hours
func classrooms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// Check http method
	if r.Method == http.MethodGet {
		getClassroomsStatus(w, r)
	} else {
		respondWithError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		log.Printf("%s /classrooms from %s status %d\n", r.Method, r.RemoteAddr, http.StatusMethodNotAllowed)
		return
	}
}

// getClassroomsStatus responds with a JSON containing all classrooms status
func getClassroomsStatus(w http.ResponseWriter, r *http.Request) {
	// TODO: filter by url params
	reservations := controllers.GetTodayReservations(occupationWeb)
	if reservations == nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		log.Printf("%s /classrooms %s status %d: ", r.Method, r.RemoteAddr, http.StatusInternalServerError)
		log.Println("Reservations slice is nil.")
		return
	}
	classrooms := controllers.GetClassroomsStatus(reservations)
	if classrooms == nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		log.Printf("%s /classrooms %s status %d: ", r.Method, r.RemoteAddr, http.StatusInternalServerError)
		log.Println("Classrooms switch reached default tag.")
		return
	}
	log.Printf("%s /classrooms %s status %d\n", r.Method, r.RemoteAddr, http.StatusOK)
	respondWithJSON(w, http.StatusOK, classrooms)
}

// occupation is the API server handler for "/occupation"
// If method is GET, it responds with a JSON containing all classrooms occupation
// else, it responds with a JSON containing an error message
func occupation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// Check http method
	if r.Method == http.MethodGet {
		getClassroomsOccupation(w, r)
	} else {
		respondWithError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		log.Printf("%s /occupation from %s status %d\n", r.Method, r.RemoteAddr, http.StatusMethodNotAllowed)
		return
	}
}

func getClassroomsOccupation(w http.ResponseWriter, r *http.Request) {
	// TODO: filter by url params
	occupations := controllers.GetClassroomsOccupation(controlServer, occupationCmd)
	if occupations == nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		log.Printf("%s /occupation %s status %d: ", r.Method, r.RemoteAddr, http.StatusInternalServerError)
		log.Println("Command execution via ssh failed.")
		return
	}
	log.Printf("%s /occupation %s status %d\n", r.Method, r.RemoteAddr, http.StatusOK)
	respondWithJSON(w, http.StatusOK, occupations)
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
