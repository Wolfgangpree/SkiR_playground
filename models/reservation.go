package models

import "time"

// Reservation repräsentiert eine Ski-Ausleihe
type Reservation struct {
	ID         int        // Auto-Increment ID
	SkiID      string     // Referenz zu Ski
	PersonID   int        // Referenz zu Person
	RentDate   time.Time  // Ausleihdatum
	ReturnDate *time.Time // Rückgabedatum (nil wenn noch nicht zurückgegeben)
}
