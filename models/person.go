package models

// Person repräsentiert einen Kunden
type Person struct {
	ID             int    // Auto-Increment ID
	FirstName      string
	LastName       string
	Email          string
	Phone          string // Festnetz (optional)
	Mobile         string // Mobiltelefonnummer (Pflichtfeld)
	IDDocumentPath string // Pfad zum Ausweisscan/Foto (optional)
}
