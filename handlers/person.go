package handlers

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reservation/models"
	"reservation/storage"
	"strconv"
	"strings"
)

type PersonHandler struct {
	store     *storage.Store
	templates *template.Template
}

func NewPersonHandler(store *storage.Store, templates *template.Template) *PersonHandler {
	return &PersonHandler{
		store:     store,
		templates: templates,
	}
}

// PersonsPage zeigt die Personen-Verwaltungsseite an
func (h *PersonHandler) PersonsPage(w http.ResponseWriter, r *http.Request) {
	persons := h.store.GetAllPersons()

	data := struct {
		Persons []*models.Person
		Error   string
	}{
		Persons: persons,
	}

	h.templates.ExecuteTemplate(w, "persons.html", data)
}

// CreatePerson fügt eine neue Person hinzu
func (h *PersonHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Multipart-Form parsen (für Datei-Upload)
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Fehler beim Parsen des Formulars", http.StatusBadRequest)
		return
	}

	firstName := r.FormValue("firstname")
	lastName := r.FormValue("lastname")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	mobile := r.FormValue("mobile")

	// Validierung: Mobile ist Pflichtfeld
	if mobile == "" {
		http.Error(w, "Mobiltelefonnummer ist ein Pflichtfeld", http.StatusBadRequest)
		return
	}

	// Datei-Upload verarbeiten (optional)
	var documentPath string
	file, header, err := r.FormFile("id_document")
	if err == nil {
		defer file.Close()

		// Validiere Dateiformat
		ext := strings.ToLower(filepath.Ext(header.Filename))
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".pdf" {
			http.Error(w, "Ungültiges Dateiformat. Nur PNG, JPG und PDF erlaubt.", http.StatusBadRequest)
			return
		}

		// Erstelle Upload-Verzeichnis falls nicht vorhanden
		uploadDir := "uploads/id_documents"
		os.MkdirAll(uploadDir, 0755)

		// Generiere eindeutigen Dateinamen
		filename := fmt.Sprintf("%d_%s", len(h.store.GetAllPersons())+1, header.Filename)
		documentPath = filepath.Join(uploadDir, filename)

		// Speichere Datei
		dst, err := os.Create(documentPath)
		if err != nil {
			http.Error(w, "Fehler beim Speichern der Datei", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Fehler beim Speichern der Datei", http.StatusInternalServerError)
			return
		}
	}

	person := &models.Person{
		FirstName:      firstName,
		LastName:       lastName,
		Email:          email,
		Phone:          phone,
		Mobile:         mobile,
		IDDocumentPath: documentPath,
	}

	_, err = h.store.CreatePerson(person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/persons", http.StatusSeeOther)
}

// DeletePerson löscht eine Person
func (h *PersonHandler) DeletePerson(w http.ResponseWriter, r *http.Request) {
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

	err = h.store.DeletePerson(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/persons", http.StatusSeeOther)
}

// ViewDocument zeigt das ID-Dokument an
func (h *PersonHandler) ViewDocument(w http.ResponseWriter, r *http.Request) {
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

	person, err := h.store.GetPerson(id)
	if err != nil {
		http.Error(w, "Person nicht gefunden", http.StatusNotFound)
		return
	}

	if person.IDDocumentPath == "" {
		http.Error(w, "Kein Dokument vorhanden", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, person.IDDocumentPath)
}
