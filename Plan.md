# Implementierungsplan - Reservierungssystem

## Übersicht

Dieser Plan beschreibt die schrittweise Implementierung des Reservierungssystems gemäß der Verfassung.md.

## Phase 1: Grundstruktur und Datenmodelle

### 1.1 Projektstruktur
```
Reservation/
├── main.go                 # Hauptprogramm, Server-Setup
├── models/
│   ├── ski.go             # Ski-Datenmodell
│   ├── person.go          # Person-Datenmodell
│   └── reservation.go     # Reservierungs-Datenmodell
├── storage/
│   └── memory.go          # In-Memory Storage mit Mutex
├── handlers/
│   ├── ski.go             # Handler für Ski-Operationen
│   ├── person.go          # Handler für Personen-Operationen
│   └── reservation.go     # Handler für Reservierungs-Operationen
├── templates/
│   ├── layout.html        # Basis-Layout
│   ├── index.html         # Startseite/Dashboard
│   ├── skis.html          # Ski-Verwaltung
│   ├── persons.html       # Personen-Verwaltung
│   └── reservations.html  # Reservierungs-Übersicht
├── static/
│   └── style.css          # CSS-Styling
└── uploads/
    └── id_documents/      # Ausweisscans/Fotos
```

### 1.2 Datenmodelle (models/)

**Ski** (`models/ski.go`)
```go
type Ski struct {
    ID        string    // 3-stellige Nummer (z.B. "001", "042", "999")
    Status    string    // "available" oder "rented"
    Brand     string    // Marke (z.B. Atomic, Fischer, OriginalPlus)
    Model     string    // Modell (optional)
    Length    int       // Länge in cm (z.B. 150, 160, 170, 180)
}
```

**Person** (`models/person.go`)
```go
type Person struct {
    ID              int       // Auto-Increment ID
    FirstName       string
    LastName        string
    Email           string
    Phone           string    // Festnetz (optional)
    Mobile          string    // Mobiltelefonnummer
    IDDocumentPath  string    // Pfad zum Ausweisscan/Foto (optional)
}
```

**Reservation** (`models/reservation.go`)
```go
type Reservation struct {
    ID          int       // Auto-Increment ID
    SkiID       string    // Referenz zu Ski
    PersonID    int       // Referenz zu Person
    RentDate    time.Time
    ReturnDate  *time.Time // nil wenn noch nicht zurückgegeben
}
```

## Phase 2: Datenspeicher

### 2.1 In-Memory Storage (`storage/memory.go`)

Implementierung eines Thread-sicheren In-Memory Speichers:

```go
type Store struct {
    mu           sync.RWMutex
    skis         map[string]*models.Ski
    persons      map[int]*models.Person
    reservations map[int]*models.Reservation
    nextPersonID int
    nextResID    int
}
```

**Operationen:**
- CRUD für Ski (Create, Read, Update, Delete)
- CRUD für Personen
- CRUD für Reservierungen
- Geschäftslogik: Ski ausleihen/zurückgeben

## Phase 3: HTTP Handler

### 3.1 Handler-Funktionen (`handlers/`)

**Ski-Handler** (`handlers/ski.go`)
- `GET /skis` - Liste aller Ski anzeigen
- `POST /skis` - Neuen Ski hinzufügen
- `DELETE /skis/{id}` - Ski entfernen

**Personen-Handler** (`handlers/person.go`)
- `GET /persons` - Liste aller Personen
- `POST /persons` - Neue Person hinzufügen (mit Multipart-Form-Upload für ID-Dokument)
- `GET /persons/{id}/document` - ID-Dokument anzeigen
- `DELETE /persons/{id}` - Person entfernen

**Reservierungs-Handler** (`handlers/reservation.go`)
- `GET /reservations` - Aktuelle Reservierungen
- `POST /reservations` - Ski ausleihen
- `PUT /reservations/{id}/return` - Ski zurückgeben

## Phase 4: Frontend

### 4.1 HTML Templates (`templates/`)

**Layout** (`layout.html`)
- Basis-HTML-Struktur
- Navigation
- CSS-Einbindung

**Dashboard** (`index.html`)
- Übersicht: Anzahl verfügbarer/ausgeliehener Ski
- Schnellzugriff auf Hauptfunktionen

**Ski-Verwaltung** (`skis.html`)
- Tabelle mit allen Ski (ID, Status, Marke, Modell, Länge)
- Formular zum Hinzufügen neuer Ski mit Längenangabe
- Löschen-Button

**Personen-Verwaltung** (`persons.html`)
- Tabelle mit allen Personen (Name, Email, Festnetz, Mobil, Ausweis-Link)
- Formular zum Hinzufügen mit Feldern:
  - Vorname, Nachname
  - Email
  - Telefon (Festnetz, optional)
  - Mobiltelefon (Pflichtfeld)
  - Ausweisscan/Foto-Upload (optional, PNG/JPG/PDF)
- Link zum Anzeigen des ID-Dokuments
- Löschen-Button

**Reservierungen** (`reservations.html`)
- Aktive Ausleihen anzeigen
- Formular: Ski ausleihen (Ski-Auswahl, Personen-Auswahl)
- Zurückgeben-Button pro Reservierung

### 4.2 Styling (`static/style.css`)
- Einfaches, übersichtliches Design
- Responsive für Tablet-Nutzung
- Tabellen-Styling
- Formular-Styling

## Phase 5: Server-Setup

### 5.1 Main (`main.go`)

```go
func main() {
    // Store initialisieren
    store := storage.NewStore()

    // Beispieldaten einfügen
    seedData(store)

    // Router/Mux konfigurieren
    mux := http.NewServeMux()

    // Routes registrieren
    // Static files
    // Template routes
    // API routes

    // Server starten
    log.Println("Server läuft auf http://localhost:8080")
    http.ListenAndServe(":8080", mux)
}
```

### 5.2 Seed-Daten (`seedData` Funktion)

Die Datenbank wird beim Start mit Testdaten vorausgefüllt:

**Ski (15-20 Stück):**
- **Atomic** (6-7 Ski):
  - ID: 001-007
  - Modelle: Redster, Vantage
  - Längen: 150cm, 160cm, 170cm, 180cm

- **Fischer** (6-7 Ski):
  - ID: 010-016
  - Modelle: RC4, Progressor
  - Längen: 155cm, 165cm, 175cm, 185cm

- **OriginalPlus** (3-6 Ski):
  - ID: 020-025
  - Modelle: Allround, Sport
  - Längen: 160cm, 170cm, 180cm

**Personen (10 Stück):**
1. Max Mustermann (max.mustermann@email.at, +43 664 1234567)
2. Anna Schmidt (anna.schmidt@email.at, +43 676 2345678)
3. Peter Huber (peter.huber@email.at, +43 650 3456789)
4. Maria Wagner (maria.wagner@email.at, +43 660 4567890)
5. Thomas Berger (thomas.berger@email.at, +43 664 5678901)
6. Lisa Mayer (lisa.mayer@email.at, +43 676 6789012)
7. Michael Bauer (michael.bauer@email.at, +43 650 7890123)
8. Sarah Hoffmann (sarah.hoffmann@email.at, +43 660 8901234)
9. David Steiner (david.steiner@email.at, +43 664 9012345)
10. Julia Koch (julia.koch@email.at, +43 676 0123456)

**Optional:** 2-3 Beispiel-Reservierungen für Testzwecke

## Implementierungsreihenfolge

1. ✅ Go-Modul initialisieren
2. ⏳ **Models erstellen** (ski.go, person.go, reservation.go)
3. ⏳ **Storage implementieren** (memory.go mit allen CRUD-Operationen)
4. ⏳ **Handler implementieren** (ski.go, person.go, reservation.go)
5. ⏳ **Templates erstellen** (HTML-Dateien)
6. ⏳ **CSS-Styling** (style.css)
7. ⏳ **Main-Server** (Routing, Server-Start, Seed-Daten)
8. ⏳ **Testing** (Manuelle Tests aller Funktionen)
9. ⏳ **Dokumentation** (README.md aktualisieren mit Nutzungsanleitung)

## Testplan

### Manuelle Tests
1. **Ski-Verwaltung**
   - Ski hinzufügen (gültige 3-stellige ID)
   - Ski hinzufügen (ungültige ID) → Fehler
   - Ski löschen

2. **Personen-Verwaltung**
   - Person hinzufügen mit allen Pflichtfeldern (Name, Mobil)
   - Person hinzufügen mit ID-Dokument-Upload
   - ID-Dokument anzeigen (PNG/JPG/PDF)
   - Person ohne Mobiltelefon hinzufügen → Fehler
   - Ungültiges Dateiformat hochladen → Fehler
   - Person löschen
   - Person löschen mit aktiver Ausleihe → Fehler

3. **Reservierungen**
   - Ski ausleihen (verfügbarer Ski)
   - Ski ausleihen (bereits ausgeliehen) → Fehler
   - Ski zurückgeben
   - Mehrere gleichzeitige Ausleihen

4. **Edge Cases**
   - Server-Neustart (Daten weg - erwartet)
   - Leere Formulare absenden
   - Duplikate Ski-IDs

## Zeitabschätzung

- Phase 1 & 2 (Models + Storage): 1-2 Stunden
- Phase 3 (Handler): 2-3 Stunden
- Phase 4 (Frontend): 2-3 Stunden
- Phase 5 (Server): 1 Stunde
- Testing & Bugfixes: 1-2 Stunden

**Gesamt: 7-11 Stunden**

## Nächste Schritte

Nach Fertigstellung des MVP können folgende Erweiterungen in Betracht gezogen werden:
1. Persistente Speicherung (SQLite)
2. Suchfunktion
3. Export-Funktion (CSV)
4. Historische Daten (vergangene Ausleihen)