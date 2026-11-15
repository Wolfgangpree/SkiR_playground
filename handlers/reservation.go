package handlers

import (
	"html/template"
	"net/http"
	"reservation/models"
	"reservation/storage"
	"strconv"
)

type ReservationHandler struct {
	store     *storage.Store
	templates *template.Template
}

func NewReservationHandler(store *storage.Store, templates *template.Template) *ReservationHandler {
	return &ReservationHandler{
		store:     store,
		templates: templates,
	}
}

// ReservationView kombiniert Reservation mit Ski und Person Daten
type ReservationView struct {
	Reservation *models.Reservation
	Ski         *models.Ski
	Person      *models.Person
}

// ReservationsPage zeigt die Reservierungsseite an
func (h *ReservationHandler) ReservationsPage(w http.ResponseWriter, r *http.Request) {
	reservations := h.store.GetActiveReservations()
	skis := h.store.GetAllSkis()
	persons := h.store.GetAllPersons()

	// Verfügbare Ski filtern
	availableSkis := make([]*models.Ski, 0)
	for _, ski := range skis {
		if ski.Status == "available" {
			availableSkis = append(availableSkis, ski)
		}
	}

	// Reservierungen mit vollständigen Daten anreichern
	reservationViews := make([]*ReservationView, 0, len(reservations))
	for _, res := range reservations {
		ski, _ := h.store.GetSki(res.SkiID)
		person, _ := h.store.GetPerson(res.PersonID)
		reservationViews = append(reservationViews, &ReservationView{
			Reservation: res,
			Ski:         ski,
			Person:      person,
		})
	}

	data := struct {
		Reservations  []*ReservationView
		AvailableSkis []*models.Ski
		Persons       []*models.Person
		Error         string
	}{
		Reservations:  reservationViews,
		AvailableSkis: availableSkis,
		Persons:       persons,
	}

	h.templates.ExecuteTemplate(w, "reservations.html", data)
}

// CreateReservation erstellt eine neue Reservierung
func (h *ReservationHandler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Fehler beim Parsen des Formulars", http.StatusBadRequest)
		return
	}

	skiID := r.FormValue("ski_id")
	personIDStr := r.FormValue("person_id")

	personID, err := strconv.Atoi(personIDStr)
	if err != nil {
		http.Error(w, "Ungültige Personen-ID", http.StatusBadRequest)
		return
	}

	_, err = h.store.CreateReservation(skiID, personID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/reservations", http.StatusSeeOther)
}

// ReturnSki markiert einen Ski als zurückgegeben
func (h *ReservationHandler) ReturnSki(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID fehlt", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Ungültige ID", http.StatusBadRequest)
		return
	}

	err = h.store.ReturnSki(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/reservations", http.StatusSeeOther)
}
