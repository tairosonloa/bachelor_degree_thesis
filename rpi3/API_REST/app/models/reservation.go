package models

// Reservation represents a classroom reservation
type Reservation struct {
	ID        int
	Classroom string
	Subject   string
	Professor string
	Day       int
	Month     int
	Year      int
	Hour      int
}
