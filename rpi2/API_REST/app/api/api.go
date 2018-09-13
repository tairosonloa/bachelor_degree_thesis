package api

import (
	"app/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	cpd        models.CPD
	authorized string // Bearer token for authorized POST
)

// Initialize initializes the API server handlers and inner state
func Initialize(tokenFile string) *http.ServeMux {
	// Handlers
	mux := http.NewServeMux()

	mux.HandleFunc("/", index)
	mux.HandleFunc("/cpd-status", cpdStatus)
	mux.HandleFunc("/cpd-update", cpdUpdate)

	mux.HandleFunc("/favicon.ico", func(_ http.ResponseWriter, _ *http.Request) {})

	// Inner state
	cpd = models.CPD{Temp: -1.0, Hum: -1.0, Light: false, UPSStatus: "online", WarningTemp: false, WarningUPS: false}
	content, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		log.Printf("ERROR: %v\n", err.Error())
		authorized = ""
	} else {
		log.Printf("Token loaded from %s\n", tokenFile)
		authorized = strings.TrimRight("Bearer "+string(content), "\n\r")
	}

	return mux
}

// index is the API server handler for "/"
// If method is GET, it responds with a JSON containing the message "API is up and running"
// else, it responds with a JSON containing an error message
func index(w http.ResponseWriter, r *http.Request) {
	// Only method GET is allowed
	if r.Method != http.MethodGet {
		err := models.ErrorAPIM{
			Error: http.StatusText(http.StatusMethodNotAllowed),
		}
		respondWithJSON(w, http.StatusMethodNotAllowed, err)
		log.Printf("%s / %d\n", r.Method, http.StatusMethodNotAllowed)
		return
	}
	// Responds with API status
	response := models.MessageAPIM{
		Message: "API is up and running",
	}
	respondWithJSON(w, http.StatusOK, response)
	log.Printf("%s / %d\n", r.Method, http.StatusOK)
}

// cpdStatus is the API server handler for "/cpd-status"
// If method is GET, it responds with a JSON containing CPD values and info
// else, it responds with a JSON containing an error message
func cpdStatus(w http.ResponseWriter, r *http.Request) {
	// Only method GET is allowed
	if r.Method != http.MethodGet {
		err := models.ErrorAPIM{
			Error: http.StatusText(http.StatusMethodNotAllowed),
		}
		respondWithJSON(w, http.StatusMethodNotAllowed, err)
		log.Printf("%s /cpd-status %d\n", r.Method, http.StatusMethodNotAllowed)
		return
	}
	response := models.CPDStatusAPIM{
		Temperature: cpd.Temp,
		Humidity:    cpd.Hum,
		UPSStatus:   cpd.UPSStatus,
	}
	respondWithJSON(w, http.StatusOK, response)
	log.Printf("%s /cpd-status %d\n", r.Method, http.StatusOK)
}

// cpdUpdate is the API server handler for "/cpd-update"
// If method is POST and authenticationis susccessful, it updates CPD values and info (inner state)
// else, it responds with a JSON containing an error message
func cpdUpdate(w http.ResponseWriter, r *http.Request) {
	// Only method POST is allowed
	if r.Method != http.MethodPost {
		err := models.ErrorAPIM{
			Error: http.StatusText(http.StatusMethodNotAllowed),
		}
		respondWithJSON(w, http.StatusMethodNotAllowed, err)
		log.Printf("%s /cpd-update %d\n", r.Method, http.StatusMethodNotAllowed)
		return
	}
	// Check authentication
	httpCode := validateToken(w, r)
	if httpCode != http.StatusOK {
		log.Printf("%s /cpd-update %d\n", r.Method, httpCode)
		return
	}
	// TODO comprobar si rpi1 or ultraheroe y rutinas
	log.Printf("%s /cpd-update %d\n", r.Method, httpCode)
}

// validateToken checks if the request is authenticated (bearer token) and authorized
func validateToken(w http.ResponseWriter, r *http.Request) int {
	token := r.Header.Get("Authorization")
	if token != authorized {
		var err models.ErrorAPIM
		if token == "" {
			err = models.ErrorAPIM{
				Error: "Authorization header not provided or empty",
			}
		} else {
			err = models.ErrorAPIM{
				Error: http.StatusText(http.StatusUnauthorized),
			}
		}
		respondWithJSON(w, http.StatusUnauthorized, err)
		return http.StatusUnauthorized
	}
	return http.StatusOK
}

// respondWithJSON responds to a request with a http code and a JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
