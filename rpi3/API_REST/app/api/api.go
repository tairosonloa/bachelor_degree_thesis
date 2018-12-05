package api

import (
	// "bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	// "strconv"

	// "rpi3/API_REST/app/controllers"
	"rpi3/API_REST/app/models"

	// Blank import because we need the sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

var (
	authorized string  // Bearer token for authorized POST
	db         *sql.DB // Database driver (scribble JSON database)
)

const (
	jsonDatetimeLayout   = "2-1-2006 15:04"   // Datetime layout to/from JSON
	sqliteDatetimeLayout = "2006-01-02 15:04" // Datetime layout to insert into sqlite
)

// Initialize initializes the API server handlers and inner state
func Initialize(apiAuthorizedToken string, databaseFile string) *http.ServeMux {
	var err error
	authorized = apiAuthorizedToken

	// Database
	log.Printf("Initializeing database on %s\n", databaseFile)
	db, err = sql.Open("sqlite3", databaseFile)
	if err != nil {
		log.Printf("Error on initialize the database: %v\n", err.Error())
		os.Exit(1)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS reservation (id INTEGER PRIMARY KEY AUTOINCREMENT, classroom TEXT, subject TEXT, professor TEXT, datetime DATETIME)")
	if err != nil {
		log.Printf("Error on initialize the database: %v\n", err.Error())
		os.Exit(1)
	}

	// Handlers
	mux := http.NewServeMux()

	mux.HandleFunc("/", index)
	mux.HandleFunc("/reservation", reservation)
	mux.HandleFunc("/reservation/{[0-9]+}", reservationID)

	mux.HandleFunc("/favicon.ico", func(_ http.ResponseWriter, _ *http.Request) {})

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
	err := decoder.Decode(&reservation)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		log.Printf("%s /reservation %s status %d\n", r.Method, r.RemoteAddr, http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	// Generate ID by concatenate date and classrom
	// generateReservationID(&reservation)

	// Parse Time
	var datetime time.Time
	datetime, err = time.Parse(jsonDatetimeLayout, reservation.Datetime)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		log.Printf("%s /reservation %s status %d\n", r.Method, r.RemoteAddr, http.StatusBadRequest)
		log.Println(err.Error())
		return
	}
	// Check if datetime is before today
	if datetime.Unix() < time.Now().Unix() {
		respondWithError(w, http.StatusBadRequest, "Reservation datetime is before today")
		log.Printf("%s /reservation %s status %d: Reservation datetime is before today\n", r.Method, r.RemoteAddr, http.StatusBadRequest)
		return
	}

	// Insert into database
	statement := fmt.Sprintf("INSERT INTO reservation (classroom, subject, professor, datetime) VALUES ('%s', '%s', '%s', '%s')",
		reservation.Classroom, reservation.Subject, reservation.Professor, datetime.Format(sqliteDatetimeLayout))
	_, err = db.Exec(statement)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		log.Printf("%s /reservation %s status %d\n", r.Method, r.RemoteAddr, http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	// Respond 201 created and return ID on a JSON
	response := models.TransactionInfoAPIM{
		ID:     reservation.ID,
		Status: "Reservation created",
	}
	log.Printf("%s /reservation %s status %d\n", r.Method, r.RemoteAddr, http.StatusCreated)
	respondWithJSON(w, http.StatusCreated, response)
}

// getReservations responds with a JSON cantaining all reservations
func getReservations(w http.ResponseWriter, r *http.Request) {
	// TODO: filter by url params
	var classroom, subject, professor string
	var datetime time.Time
	var reservation models.Reservation
	payload := []models.Reservation{}
	rows, err := db.Query("SELECT classroom, subject, professor, datetime FROM reservation")
	defer rows.Close()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		log.Printf("%s /reservation %s status %d\n", r.Method, r.RemoteAddr, http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	for rows.Next() {
		err = rows.Scan(&classroom, &subject, &professor, &datetime)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			log.Printf("%s /reservation %s status %d\n", r.Method, r.RemoteAddr, http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		reservation = models.Reservation{"", classroom, subject, professor, datetime.Format(jsonDatetimeLayout)}
		payload = append(payload, reservation)
	}
	log.Printf("%s /reservation %s status %d\n", r.Method, r.RemoteAddr, http.StatusOK)
	respondWithJSON(w, http.StatusOK, payload)
}

// getReservationByID responds with a JSON cantaining a reservation info
func getReservationByID(w http.ResponseWriter, r *http.Request) {
	// TODO: to do
	respondWithError(w, http.StatusNotImplemented, "Not implemented")
}

// updateReservationByID updates a existing reservation from a JSON PUT request
func updateReservationByID(w http.ResponseWriter, r *http.Request) {
	// TODO: to do
	respondWithError(w, http.StatusNotImplemented, "Not implemented")
}

// validateToken checks if the request is authenticated (bearer token) and authorized
func validateToken(w http.ResponseWriter, r *http.Request) bool {
	token := r.Header.Get("Authorization")
	if token != authorized {
		var err string
		log.Println(token)
		log.Println(authorized)
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

// generateReservationID generates a reservation ID and saves it on the reservation struct
// func generateReservationID(reservation *models.Reservation) {
// 	var id bytes.Buffer
// 	id.WriteString(strconv.Itoa(reservation.Year))
// 	id.WriteString(".")
// 	id.WriteString(strconv.Itoa(reservation.Month))
// 	id.WriteString(".")
// 	id.WriteString(strconv.Itoa(reservation.Day))
// 	id.WriteString(".")
// 	id.WriteString(strconv.Itoa(reservation.Hour))
// 	id.WriteString("h.")
// 	id.WriteString(reservation.Classroom)
// 	reservation.ID = id.String()
// }

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
