package main

import (
	"html/template"
	"log"
	"net/http"
	"reservation/handlers"
	"reservation/models"
	"reservation/storage"
)

func main() {
	// Store initialisieren
	store := storage.NewStore()

	// Seed-Daten einfügen
	seedData(store)

	// Templates laden
	templates := template.Must(template.ParseGlob("templates/*.html"))

	// Handler erstellen
	skiHandler := handlers.NewSkiHandler(store, templates)
	personHandler := handlers.NewPersonHandler(store, templates)
	reservationHandler := handlers.NewReservationHandler(store, templates)

	// Router konfigurieren
	mux := http.NewServeMux()

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Dashboard
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		skis := store.GetAllSkis()
		persons := store.GetAllPersons()

		availableCount := 0
		rentedCount := 0
		for _, ski := range skis {
			if ski.Status == "available" {
				availableCount++
			} else {
				rentedCount++
			}
		}

		data := struct {
			AvailableCount int
			RentedCount    int
			TotalSkis      int
			TotalPersons   int
		}{
			AvailableCount: availableCount,
			RentedCount:    rentedCount,
			TotalSkis:      len(skis),
			TotalPersons:   len(persons),
		}

		templates.ExecuteTemplate(w, "index.html", data)
	})

	// Ski Routes
	mux.HandleFunc("/skis", skiHandler.SkisPage)
	mux.HandleFunc("/skis/create", skiHandler.CreateSki)
	mux.HandleFunc("/skis/delete", skiHandler.DeleteSki)

	// Person Routes
	mux.HandleFunc("/persons", personHandler.PersonsPage)
	mux.HandleFunc("/persons/create", personHandler.CreatePerson)
	mux.HandleFunc("/persons/delete", personHandler.DeletePerson)
	mux.HandleFunc("/persons/document", personHandler.ViewDocument)

	// Reservation Routes
	mux.HandleFunc("/reservations", reservationHandler.ReservationsPage)
	mux.HandleFunc("/reservations/create", reservationHandler.CreateReservation)
	mux.HandleFunc("/reservations/return", reservationHandler.ReturnSki)

	// Server starten
	log.Println("Server läuft auf http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// seedData füllt die Datenbank mit Beispieldaten
func seedData(store *storage.Store) {
	// Ski hinzufügen - Atomic (001-007)
	atomicModels := []struct {
		ID     string
		Model  string
		Length int
	}{
		{"001", "Redster X9", 170},
		{"002", "Redster X7", 160},
		{"003", "Vantage 82", 175},
		{"004", "Vantage 90", 180},
		{"005", "Redster G9", 165},
		{"006", "Vantage 75", 150},
		{"007", "Redster S9", 170},
	}

	for _, ski := range atomicModels {
		store.CreateSki(&models.Ski{
			ID:     ski.ID,
			Status: "available",
			Brand:  "Atomic",
			Model:  ski.Model,
			Length: ski.Length,
		})
	}

	// Ski hinzufügen - Fischer (010-016)
	fischerModels := []struct {
		ID     string
		Model  string
		Length int
	}{
		{"010", "RC4 Worldcup", 185},
		{"011", "RC4 Superior", 175},
		{"012", "Progressor F19", 165},
		{"013", "Progressor F17", 155},
		{"014", "RC4 The Curv", 180},
		{"015", "Progressor F18", 170},
		{"016", "RC One 82", 175},
	}

	for _, ski := range fischerModels {
		store.CreateSki(&models.Ski{
			ID:     ski.ID,
			Status: "available",
			Brand:  "Fischer",
			Model:  ski.Model,
			Length: ski.Length,
		})
	}

	// Ski hinzufügen - OriginalPlus (020-025)
	originalPlusModels := []struct {
		ID     string
		Model  string
		Length int
	}{
		{"020", "Allround Pro", 170},
		{"021", "Allround Plus", 160},
		{"022", "Sport Classic", 180},
		{"023", "Sport Elite", 170},
		{"024", "Allround Expert", 175},
		{"025", "Sport Advanced", 165},
	}

	for _, ski := range originalPlusModels {
		store.CreateSki(&models.Ski{
			ID:     ski.ID,
			Status: "available",
			Brand:  "OriginalPlus",
			Model:  ski.Model,
			Length: ski.Length,
		})
	}

	// Personen hinzufügen
	persons := []struct {
		FirstName string
		LastName  string
		Email     string
		Phone     string
		Mobile    string
	}{
		{"Max", "Mustermann", "max.mustermann@email.at", "01 1234567", "+43 664 1234567"},
		{"Anna", "Schmidt", "anna.schmidt@email.at", "", "+43 676 2345678"},
		{"Peter", "Huber", "peter.huber@email.at", "01 2345678", "+43 650 3456789"},
		{"Maria", "Wagner", "maria.wagner@email.at", "", "+43 660 4567890"},
		{"Thomas", "Berger", "thomas.berger@email.at", "01 3456789", "+43 664 5678901"},
		{"Lisa", "Mayer", "lisa.mayer@email.at", "", "+43 676 6789012"},
		{"Michael", "Bauer", "michael.bauer@email.at", "01 4567890", "+43 650 7890123"},
		{"Sarah", "Hoffmann", "sarah.hoffmann@email.at", "", "+43 660 8901234"},
		{"David", "Steiner", "david.steiner@email.at", "", "+43 664 9012345"},
		{"Julia", "Koch", "julia.koch@email.at", "01 5678901", "+43 676 0123456"},
	}

	for _, p := range persons {
		store.CreatePerson(&models.Person{
			FirstName: p.FirstName,
			LastName:  p.LastName,
			Email:     p.Email,
			Phone:     p.Phone,
			Mobile:    p.Mobile,
		})
	}

	// Optional: 2-3 Beispiel-Reservierungen
	store.CreateReservation("001", 1) // Max Mustermann leiht Atomic Redster X9
	store.CreateReservation("010", 3) // Peter Huber leiht Fischer RC4 Worldcup
	store.CreateReservation("020", 5) // Thomas Berger leiht OriginalPlus Allround Pro

	log.Println("Seed-Daten erfolgreich geladen:")
	log.Printf("- %d Ski hinzugefügt", len(store.GetAllSkis()))
	log.Printf("- %d Personen hinzugefügt", len(store.GetAllPersons()))
	log.Printf("- %d Reservierungen erstellt", len(store.GetActiveReservations()))
}
