package models

// Reservation represents a classroom reservation
type Reservation struct {
	ID        string
	Classroom string
	Subject   string
	Professor string
	Datetime  string
}
