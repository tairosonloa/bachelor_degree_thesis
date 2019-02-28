package models

// Reservation represents a classroom reservation
type Reservation struct {
	Subject     string
	Study       string
	Classroom   string
	StartHour   int
	StartMinute int
	EndHour     int
	EndMinute   int
	Professor   string
}
