package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"strings"

	"../models"
)

const (
	reservationsWebPage = "http://www.lab.inf.uc3m.es/informacion/ocupacion-de-las-aulas/ocupacion-diaria/"
)

var ( // TODO: remove from global vars (?), Anyway, they are inmutable by nature
	classrooms = [...]string{"4.0.F16", "4.0.F18", "2.2.C05", "2.2.C06"}
	hours      = [...]string{"9:00", "11:00", "13:00", "15:00", "17:00", "20:00"}
)

// fecthURL gets the web page body and returns it
func fecthURL(url string) io.ReadCloser {
	resp, err := http.Get(url)
	// Check errors
	if err != nil {
		log.Printf("ERROR crawler/fecthURL(): %v\n", err.Error())
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("ERROR crawler/fecthURL(): http status code: %v\n", resp.StatusCode)
		return nil
	}
	return resp.Body
}

func getReservations(body io.ReadCloser) []models.Reservation {
	if body == nil { // TODO: update error message
		fmt.Println("Body is nil")
		return nil
	}
	defer body.Close() // Don't forget to close the Reader

	tokenizer := html.NewTokenizer(body)
	reservations := []models.Reservation{}

	var text string
	var colum int
	row := -1
	// While have not hit the </html> tag
	for tokenizer.Token().Data != "html" {
		tagToken := tokenizer.Next()
		if tagToken == html.StartTagToken {
			token := tokenizer.Token().Data
			if token == "tr" { // row
				row++
				colum = -1
			}
			if token == "td" { // colum
				inner := tokenizer.Next()
				if inner == html.TextToken {
					colum++
					// Inside table cell
					// Step one: Get subject name
					text = (string)(tokenizer.Text())
					subject := strings.TrimSpace(text)
					if subject != "" { // Ignore empty cells
						// Step two: Get study (degree master) name
						// Ignore html tags
						inner = tokenizer.Next()
						for inner != html.TextToken {
							inner = tokenizer.Next()
						}
						// Get study name
						text = (string)(tokenizer.Text())
						study := strings.TrimSpace(text)

						// Create and append reservation object
						reservation := models.Reservation{
							Subject:   subject,
							Study:     study,
							Classroom: classrooms[colum],
							Time:      hours[row/8],
							Professor: "",
						}
						reservations = append(reservations, reservation)
					}
				}
			}
		}
	}
	return reservations
}

func main() {
	body := fecthURL(reservationsWebPage)
	if body == nil {
		// TODO: better error handling
		log.Println("ERROR: cannot fecth reservations URL. Body is nil.")
		return
	}
	reservations := getReservations(body)

	fmt.Printf("\n")
	for i := range reservations {
		fmt.Printf("%+v\n", reservations[i])
	}
	fmt.Printf("\n")
}
