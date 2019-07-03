package controllers

import (
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"rpi3/API_REST/app/models"
)

const (
	tableSize = 48
)

// join concats strings
func join(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

func capitalize(str string) string {
	return strings.Title(strings.ToLower(str))
}

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

// getReservations returns an slice with all the reservations for today (structs)
func getReservations(body io.ReadCloser) []*models.Reservation {
	if body == nil { // TODO: update error message
		log.Println("Body is nil")
		return nil
	}
	defer body.Close() // Don't forget to close the Reader

	tokenizer := html.NewTokenizer(body)
	reservations := []*models.Reservation{}
	oneTimeReservations := []*models.Reservation{}
	oneTimeReservation := false
	classrooms := [...]string{"4.0.F16", "4.0.F18", "2.2.C05", "2.2.C06"}

	var text string
	var colum int
	row := -1
	toCheck := 0

	// While have not hit the </html> tag
	for tokenizer.Token().Data != "html" {
		tagToken := tokenizer.Next()
		if tagToken == html.StartTagToken {
			token := tokenizer.Token()
			if token.Data == "tr" { // row
				row++
				colum = -1
			} else if token.Data == "td" { // colum
				colum++
				// Check rowspan attr to calculate endtime
				rowspan := -1
				for _, attr := range token.Attr {
					if (string)(attr.Key) == "rowspan" {
						rowspan, _ = strconv.Atoi(attr.Val)
						break
					}
				}
				if rowspan >= 0 { // If rowspan < 0 then there is not reservation for this hour
					inner := tokenizer.Next()
					if inner == html.TextToken {
						// Inside table cell
						// Step one: Get subject name
						text = (string)(tokenizer.Text())
						subject := strings.TrimSpace(text)
						if subject != "" { // Ignore empty cells
							study := ""
							group := -1
							if strings.Contains(strings.ToLower(subject), "reserva puntual") {
								// We detected a one-time reservation
								if strings.Contains(strings.ToLower(subject), "reserva puntual:") { // Two dots means it have subject info
									// Get subject for the one-time reservation and concat the two strings
									// Ignore html tags
									inner = tokenizer.Next()
									for inner != html.TextToken {
										inner = tokenizer.Next()
									}
									// Get subject
									text = (string)(tokenizer.Text())
									text = strings.TrimSpace(text)
									subject = join(subject, " ", text)

									// Get count of every one-time reservation to check professor later
									oneTimeReservation = true
									toCheck++
								}
							} else {
								// Step two: Get study (degree master) name
								// Ignore html tags
								inner = tokenizer.Next()
								for inner != html.TextToken {
									inner = tokenizer.Next()
								}
								// Get study name
								text = (string)(tokenizer.Text())
								study = strings.TrimSpace(text)
							}

							// Steep three, separate subject from group
							trimmed := strings.Split(subject, "(")
							subject = trimmed[0]
							if len(trimmed) > 1 {
								group, _ = strconv.Atoi(strings.TrimPrefix(strings.TrimSuffix(trimmed[1], ")"), "G"))
							}

							// Step four calculate end time
							log.Printf("Asignatura %v row %v", subject, row)
							startTimeH := 9 + (int)(row/4)            // start hour
							startTimeM := 15 * (row%4 - 1)            // start minutes
							endTimeH := startTimeH + (int)(rowspan/4) // end hour
							endTimeM := 15 * (rowspan % 4)            // end minutes

							// Step four: Create and append reservation object
							reservation := models.Reservation{
								Subject:     subject,
								Study:       study,
								Group:       group,
								Classroom:   classrooms[colum],
								StartHour:   startTimeH,
								StartMinute: startTimeM,
								EndHour:     endTimeH,
								EndMinute:   endTimeM,
								Professor:   "",
							}
							reservations = append(reservations, &reservation)

							if oneTimeReservation {
								// Check if one-time reservation, to check professor later
								oneTimeReservations = append(oneTimeReservations, reservations[len(reservations)-1])
								oneTimeReservation = false
							}
						}
					}
				}
			}
		}
		if row > tableSize {
			break
		}
	}
	// If we have one time reservations, get professors
	if toCheck > 0 {
		// While have not hit the </html> tag
		for tokenizer.Token().Data != "html" && toCheck > 0 {
			if tokenizer.Next() == html.StartTagToken && tokenizer.Token().Data == "td" {
				inner := tokenizer.Next()
				if inner == html.TextToken {
					// Inside table cell
					// Get subject name
					text = (string)(tokenizer.Text())
					subject := strings.TrimSpace(text)
					for i := range oneTimeReservations {
						// Search one-time reservation for that subject
						if strings.Contains(oneTimeReservations[i].Subject, subject) {
							// Search associated professor
							inner = tokenizer.Next()
							for inner != html.TextToken {
								inner = tokenizer.Next()
							}
							// Update reservation object
							text = (string)(tokenizer.Text())
							oneTimeReservations[i].Professor = capitalize(strings.TrimSpace(text))
							toCheck--
							break
						}
					}
				}
			}
		}
	}
	return reservations
}

// GetTodayReservations returns a slice of pointers to models.Reservation struct with today reservations
func GetTodayReservations(url string) []*models.Reservation {
	body := fecthURL(url)
	if body == nil {
		// TODO: better error handling
		log.Println("ERROR: cannot fecth reservations URL. Body is nil.")
		return nil
	}
	return getReservations(body)
}
