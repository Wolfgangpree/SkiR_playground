
# Verfassung - Reservierungssystem für Skiverleih

## 1. Zweck und Zielsetzung

Das Reservierungssystem dient der Verwaltung von Skiausrüstung und deren Ausleihe an Personen. Die Anwendung soll minimalistisch und benutzerfreundlich sein.

## 2. Funktionale Anforderungen

### 2.1 Kernfunktionalität
- Verwaltung von Ski mit 3-stelliger Identifikationsnummer
- Verwaltung von Personen (Ausleihenden)
- Zuordnung von Ski zu Personen (Ausleihe)
- Rückgabe von Ski

### 2.2 Datenverwaltung
- Ski: ID (3-stellig), Status (verfügbar/ausgeliehen)
- Person: Name, Kontaktdaten
- Ausleihe: Ski-ID, Person, Ausleihdatum

## 3. Nicht-funktionale Anforderungen

### 3.1 Technologie
- **Programmiersprache**: Go
- **Architektur**: Webanwendung mit HTTP-Server
- **Datenhaltung**: In-Memory (einfache Lösung ohne Datenbank)
- **Frontend**: HTML-Templates, minimales CSS

### 3.2 Qualitätsmerkmale
- **Einfachheit**: Minimale Funktionalität, keine unnötigen Features
- **Benutzerfreundlichkeit**: Intuitive Weboberfläche
- **Wartbarkeit**: Klare Struktur, gut lesbar
- **Zuverlässigkeit**: Fehlerbehandlung bei ungültigen Eingaben

## 4. Systemarchitektur

### 4.1 Schichtenmodell
```
┌─────────────────────────────┐
│   Präsentationsschicht      │ HTML Templates, Static Files
├─────────────────────────────┤
│   Anwendungsschicht         │ HTTP Handlers, Business Logic
├─────────────────────────────┤
│   Datenschicht              │ In-Memory Storage
└─────────────────────────────┘
```

### 4.2 Komponenten
- **Models**: Datenstrukturen (Ski, Person, Reservation)
- **Storage**: In-Memory Datenverwaltung mit Concurrency-Sicherheit
- **Handlers**: HTTP Request Handler
- **Templates**: HTML-Vorlagen für die Benutzeroberfläche
- **Main**: Server-Setup und Routing

## 5. Einschränkungen und Annahmen

### 5.1 Einschränkungen
- Keine persistente Datenspeicherung (Daten gehen bei Server-Neustart verloren)
- Keine Benutzerauthentifizierung
- Keine mehrstufigen Buchungsprozesse
- Keine Zahlungsabwicklung

### 5.2 Annahmen
- Single-User-Betrieb oder wenige gleichzeitige Benutzer
- Vertrauenswürdige Benutzer (kein Missbrauch)
- Lokales Netzwerk oder kontrollierte Umgebung

## 6. Entwicklungsprinzipien

1. **KISS** (Keep It Simple, Stupid): Einfachste Lösung wählen
2. **YAGNI** (You Aren't Gonna Need It): Nur implementieren, was benötigt wird
3. **DRY** (Don't Repeat Yourself): Code-Wiederverwendung
4. **Separation of Concerns**: Klare Trennung der Verantwortlichkeiten

## 7. Erweiterbarkeit

Bei zukünftigem Bedarf können folgende Features ergänzt werden:
- Persistente Datenspeicherung (z.B. SQLite)
- Zusätzliche Ausrüstung (Skischuhe, Stöcke, Snowboards)
- Reservierungszeiträume und Kalender
- Preisgestaltung und Rechnungen
- Reporting und Statistiken