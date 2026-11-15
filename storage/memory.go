package storage

import (
	"errors"
	"reservation/models"
	"sync"
	"time"
)

// Store ist ein Thread-sicherer In-Memory Datenspeicher
type Store struct {
	mu           sync.RWMutex
	skis         map[string]*models.Ski
	persons      map[int]*models.Person
	reservations map[int]*models.Reservation
	nextPersonID int
	nextResID    int
}

// NewStore erstellt einen neuen Store
func NewStore() *Store {
	return &Store{
		skis:         make(map[string]*models.Ski),
		persons:      make(map[int]*models.Person),
		reservations: make(map[int]*models.Reservation),
		nextPersonID: 1,
		nextResID:    1,
	}
}

// --- Ski CRUD Operationen ---

// CreateSki fügt einen neuen Ski hinzu
func (s *Store) CreateSki(ski *models.Ski) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.skis[ski.ID]; exists {
		return errors.New("ski mit dieser ID existiert bereits")
	}

	s.skis[ski.ID] = ski
	return nil
}

// GetSki gibt einen Ski zurück
func (s *Store) GetSki(id string) (*models.Ski, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ski, exists := s.skis[id]
	if !exists {
		return nil, errors.New("ski nicht gefunden")
	}
	return ski, nil
}

// GetAllSkis gibt alle Ski zurück
func (s *Store) GetAllSkis() []*models.Ski {
	s.mu.RLock()
	defer s.mu.RUnlock()

	skis := make([]*models.Ski, 0, len(s.skis))
	for _, ski := range s.skis {
		skis = append(skis, ski)
	}
	return skis
}

// DeleteSki löscht einen Ski
func (s *Store) DeleteSki(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.skis[id]; !exists {
		return errors.New("ski nicht gefunden")
	}

	delete(s.skis, id)
	return nil
}

// --- Person CRUD Operationen ---

// CreatePerson fügt eine neue Person hinzu
func (s *Store) CreatePerson(person *models.Person) (*models.Person, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	person.ID = s.nextPersonID
	s.nextPersonID++
	s.persons[person.ID] = person
	return person, nil
}

// GetPerson gibt eine Person zurück
func (s *Store) GetPerson(id int) (*models.Person, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	person, exists := s.persons[id]
	if !exists {
		return nil, errors.New("person nicht gefunden")
	}
	return person, nil
}

// GetAllPersons gibt alle Personen zurück
func (s *Store) GetAllPersons() []*models.Person {
	s.mu.RLock()
	defer s.mu.RUnlock()

	persons := make([]*models.Person, 0, len(s.persons))
	for _, person := range s.persons {
		persons = append(persons, person)
	}
	return persons
}

// DeletePerson löscht eine Person
func (s *Store) DeletePerson(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.persons[id]; !exists {
		return errors.New("person nicht gefunden")
	}

	// Prüfe ob Person aktive Reservierungen hat
	for _, res := range s.reservations {
		if res.PersonID == id && res.ReturnDate == nil {
			return errors.New("person hat aktive ausleihen und kann nicht gelöscht werden")
		}
	}

	delete(s.persons, id)
	return nil
}

// --- Reservation CRUD Operationen ---

// CreateReservation erstellt eine neue Reservierung
func (s *Store) CreateReservation(skiID string, personID int) (*models.Reservation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Prüfe ob Ski existiert
	ski, exists := s.skis[skiID]
	if !exists {
		return nil, errors.New("ski nicht gefunden")
	}

	// Prüfe ob Ski verfügbar ist
	if ski.Status != "available" {
		return nil, errors.New("ski ist bereits ausgeliehen")
	}

	// Prüfe ob Person existiert
	if _, exists := s.persons[personID]; !exists {
		return nil, errors.New("person nicht gefunden")
	}

	// Erstelle Reservierung
	reservation := &models.Reservation{
		ID:         s.nextResID,
		SkiID:      skiID,
		PersonID:   personID,
		RentDate:   time.Now(),
		ReturnDate: nil,
	}
	s.nextResID++
	s.reservations[reservation.ID] = reservation

	// Aktualisiere Ski-Status
	ski.Status = "rented"

	return reservation, nil
}

// GetReservation gibt eine Reservierung zurück
func (s *Store) GetReservation(id int) (*models.Reservation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	reservation, exists := s.reservations[id]
	if !exists {
		return nil, errors.New("reservierung nicht gefunden")
	}
	return reservation, nil
}

// GetAllReservations gibt alle Reservierungen zurück
func (s *Store) GetAllReservations() []*models.Reservation {
	s.mu.RLock()
	defer s.mu.RUnlock()

	reservations := make([]*models.Reservation, 0, len(s.reservations))
	for _, reservation := range s.reservations {
		reservations = append(reservations, reservation)
	}
	return reservations
}

// GetActiveReservations gibt alle aktiven Reservierungen zurück
func (s *Store) GetActiveReservations() []*models.Reservation {
	s.mu.RLock()
	defer s.mu.RUnlock()

	reservations := make([]*models.Reservation, 0)
	for _, reservation := range s.reservations {
		if reservation.ReturnDate == nil {
			reservations = append(reservations, reservation)
		}
	}
	return reservations
}

// ReturnSki markiert einen Ski als zurückgegeben
func (s *Store) ReturnSki(reservationID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	reservation, exists := s.reservations[reservationID]
	if !exists {
		return errors.New("reservierung nicht gefunden")
	}

	if reservation.ReturnDate != nil {
		return errors.New("ski wurde bereits zurückgegeben")
	}

	// Markiere Rückgabe
	now := time.Now()
	reservation.ReturnDate = &now

	// Aktualisiere Ski-Status
	if ski, exists := s.skis[reservation.SkiID]; exists {
		ski.Status = "available"
	}

	return nil
}
