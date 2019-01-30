package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"strings"
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

func getReservations(body io.ReadCloser) [][]string {
	if body == nil { // TODO: update error message
		fmt.Println("Body is nil")
		return nil
	}
	defer body.Close() // Don't forget to close the Reader

	tokenizer := html.NewTokenizer(body)
	reservations := [][]string{}

	// While have not hit the </html> tag
	counter := -2
	var text string
	for tokenizer.Token().Data != "html" {
		tagToken := tokenizer.Next()
		if tagToken == html.StartTagToken {
			token := tokenizer.Token()

			if token.Data == "tr" { // row
				counter++

			} else if token.Data == "td" && counter%8 == 0 { // colum
				inner := tokenizer.Next()
				if inner == html.TextToken {
					// Inside table cell
					reservation := []string{}
					// Step one: Get subject name
					text = (string)(tokenizer.Text())
					subject := strings.TrimSpace(text)
					// Step two: Get study (degree master) name
					// Ignore html tags
					inner = tokenizer.Next()
					for inner != html.TextToken {
						inner = tokenizer.Next()
					}
					// Get study name
					text = (string)(tokenizer.Text())
					study := strings.TrimSpace(text)
					// Append to the reservations slice
					reservation = append(reservation, subject)
					reservation = append(reservation, study)
					reservations = append(reservations, reservation)
				}
			}
		}
	}
	return reservations
}

func main() {
	reservations := getReservations(fecthURL(reservationsWebPage))

	for i1, e1 := range hours {
		for i2, e2 := range classrooms {
			if reservations[i1+i2][0] != "" {
				fmt.Printf("Clase de %s del %s en el aula %s a las %s\n", reservations[i1+i2][0], reservations[i1+i2][1], e2, e1)
			}
		}
	}
}
