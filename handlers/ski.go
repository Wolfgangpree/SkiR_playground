package handlers

import (
	"html/template"
	"net/http"
	"regexp"
	"reservation/models"
	"reservation/storage"
	"strconv"
)

type SkiHandler struct {
	store     *storage.Store
	templates *template.Template
}

func NewSkiHandler(store *storage.Store, templates *template.Template) *SkiHandler {
	return &SkiHandler{
		store:     store,
		templates: templates,
	}
}

// SkisPage zeigt die Ski-Verwaltungsseite an
func (h *SkiHandler) SkisPage(w http.ResponseWriter, r *http.Request) {
	skis := h.store.GetAllSkis()

	data := struct {
		Skis  []*models.Ski
		Error string
	}{
		Skis: skis,
	}

	h.templates.ExecuteTemplate(w, "skis.html", data)
}

// CreateSki fügt einen neuen Ski hinzu
func (h *SkiHandler) CreateSki(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Formular parsen
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Fehler beim Parsen des Formulars", http.StatusBadRequest)
		return
	}

	id := r.FormValue("id")
	brand := r.FormValue("brand")
	model := r.FormValue("model")
	lengthStr := r.FormValue("length")

	// Validierung: ID muss 3-stellig sein
	matched, _ := regexp.MatchString(`^\d{3}$`, id)
	if !matched {
		http.Error(w, "ID muss eine 3-stellige Nummer sein", http.StatusBadRequest)
		return
	}

	// Länge parsen
	length, err := strconv.Atoi(lengthStr)
	if err != nil || length <= 0 {
		http.Error(w, "Ungültige Länge", http.StatusBadRequest)
		return
	}

	ski := &models.Ski{
		ID:     id,
		Status: "available",
		Brand:  brand,
		Model:  model,
		Length: length,
	}

	err = h.store.CreateSki(ski)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Zurück zur Ski-Seite
	http.Redirect(w, r, "/skis", http.StatusSeeOther)
}

// DeleteSki löscht einen Ski
func (h *SkiHandler) DeleteSki(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID fehlt", http.StatusBadRequest)
		return
	}

	err := h.store.DeleteSki(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/skis", http.StatusSeeOther)
}
