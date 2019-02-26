package api

import (
	"encoding/json"
	"log"
	"net/http"

	"rpi2/API_REST/app/controllers"
	"rpi2/API_REST/app/models"
)

var (
	cpd        models.CPD
	authorized string // Bearer token for authorized POST
	alarmSound string // Path tho the alarm sound
)

// Initialize initializes the API server handlers and inner state
func Initialize(apiAuthorizedToken string, hueBridgeAddress string, hueBridgeToken string, alarmSoundPath string) *http.ServeMux {
	// Handlers
	mux := http.NewServeMux()

	mux.HandleFunc("/", index)
	mux.HandleFunc("/cpd-status", cpdStatus)
	mux.HandleFunc("/cpd-update", cpdUpdate)

	mux.HandleFunc("/favicon.ico", func(_ http.ResponseWriter, _ *http.Request) {})

	// Inner state
	cpd = models.CPD{Temp: -1.0, Hum: -1.0, Light: false, UPSStatus: "online", WarningTemp: false, WarningUPS: false}
	authorized = apiAuthorizedToken
	alarmSound = alarmSoundPath

	// Connect to philips hue bridge
	controllers.InitializeHue(hueBridgeAddress, hueBridgeToken)
	controllers.LightOff() // Reset light

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

// cpdStatus is the API server handler for "/cpd-status"
// If method is GET, it responds with a JSON containing CPD values and info
// else, it responds with a JSON containing an error message
func cpdStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// Only method GET is allowed
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		log.Printf("%s /cpd-status from %s status %d\n", r.Method, r.RemoteAddr, http.StatusMethodNotAllowed)
		return
	}

	// Check if URL params
	params := r.URL.Query()
	if len(params) == 0 {
		// Respond all
		response := models.CPDStatusAPIM{
			Temperature: cpd.Temp,
			Humidity:    cpd.Hum,
			UPSStatus:   cpd.UPSStatus,
		}
		respondWithJSON(w, http.StatusOK, response)
	} else {
		// Filter response by params
		response := make(map[string]interface{})
		if _, ok := params["ups"]; ok {
			response["ups status (LDI rack)"] = cpd.UPSStatus
		}
		if _, ok := params["hum"]; ok {
			response["humidity"] = cpd.Hum
		}
		if _, ok := params["temp"]; ok {
			response["temperature"] = cpd.Temp
		}
		respondWithJSON(w, http.StatusOK, response)
	}

	log.Printf("%s /cpd-status from %s status %d\n", r.Method, r.RemoteAddr, http.StatusOK)
}

// cpdUpdate is the API server handler for "/cpd-update"
// If method is POST and authentication is susccessful, it updates CPD values and info (inner state)
// else, it responds with a JSON containing an error message
func cpdUpdate(w http.ResponseWriter, r *http.Request) {
	// Only method POST is allowed
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		log.Printf("%s /cpd-update from %s status %d\n", r.Method, r.RemoteAddr, http.StatusMethodNotAllowed)
		return
	}
	// Check authentication
	if !validateToken(w, r) {
		log.Printf("%s /cpd-update from %s status %d\n", r.Method, r.RemoteAddr, http.StatusUnauthorized)
		return
	}

	// Read JSON from body request and update CPD values
	oldLightStatus := cpd.Light
	decoder := json.NewDecoder(r.Body)
	e := decoder.Decode(&cpd)
	if e != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		log.Printf("%s /cpd-update from %s status %d\n", r.Method, r.RemoteAddr, http.StatusInternalServerError)
		log.Println(e.Error())
		return
	}

	// Check if light inside CPD has changed
	checkLightStatusChanged(oldLightStatus)
	// Check if in a warning state after inner state update
	checkNewWarningStatus()

	// Respond with JSON and 200
	m := models.MessageAPIM{
		Message: "OK",
	}
	respondWithJSON(w, http.StatusOK, m)
	log.Printf("%s /cpd-update from %s status %d %+v\n", r.Method, r.RemoteAddr, http.StatusOK, cpd)
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

// checkLightStatusChanged checks if light inside CPD has changed and sets proper visual identifier
func checkLightStatusChanged(oldLightStatus bool) {
	if cpd.Light != oldLightStatus && !cpd.IsWarning() {
		if cpd.Light {
			controllers.LightON()
			log.Println("INFO: light on")
		} else {
			controllers.LightOff()
			log.Println("INFO: light off")
		}
	}
}

// checkNewWarningStatus checks if were are in a warning state and fires alarm
func checkNewWarningStatus() {
	onAlarm := cpd.IsWarning()
	// Check high temperature
	if cpd.Temp >= 30 {
		cpd.WarningTemp = true
		log.Println("WARNING: CPD temperature is high (>= 30ºC)")
	} else if cpd.WarningTemp {
		cpd.WarningTemp = false
		log.Println("WARNING-UPDATE: CPD temperature is safe (< 30ºC)")
	}

	// Check power cut
	if cpd.UPSStatus == "battery" {
		cpd.WarningUPS = true
		log.Println("WARNING: UPS on battery")
	} else if cpd.WarningUPS {
		cpd.WarningUPS = false
		log.Println("WARNING-UPDATE: UPS online")
	}

	// If there is an alarm and we are not on a previous alarm, fire it
	if !onAlarm && cpd.IsWarning() {
		controllers.FireAlarm(&cpd, alarmSound)
	}

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
