package models

// Ski repräsentiert einen Ski im Verleih
type Ski struct {
	ID     string // 3-stellige Nummer (z.B. "001", "042", "999")
	Status string // "available" oder "rented"
	Brand  string // Marke (z.B. Atomic, Fischer, OriginalPlus)
	Model  string // Modell (optional)
	Length int    // Länge in cm (z.B. 150, 160, 170, 180)
}
