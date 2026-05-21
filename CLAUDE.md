# Microservices-Workshop

Workshop-Projekt zur Vermittlung von Architekturmustern und Best Practices für
Microservices. Sprache des gesamten Workshop-Materials: Deutsch.

## Zielgruppe & Sprache

- Zielgruppe: Software-Architekt:innen, Entwickler:innen auf dem Weg zur
  Architekturrolle, Tech-Leads mit Architekturverantwortung
- Alle Workshop-Inhalte (Dokumentation, Stories, Trainer-Hinweise) sind auf
  Deutsch verfasst
- Code, Bezeichner und Commit-Messages: Englisch

## Tech-Stack

- Die Referenz-Implementierung in `services/` ist in Go geschrieben (aktuell 1.25)
- Der Tech-Stack der Teilnehmenden ist bewusst **frei wählbar**: jede Sprache
  oder jedes Framework, das zur Aufgabe passt, ist erlaubt
- Beispiel-Domäne: Reise-Buchungssystem (Hotel, Flug, Mietwagen,
  Buchungs-Orchestrierung)

## Workshop-Themen

Die vollständige Themenübersicht liegt in `docs/themen.md`. Die Themen umfassen
Grundlagen, Resilience-Patterns, Kommunikation & Routing, Daten & Events,
Deployment & Betrieb sowie Kultur & Organisation.

Die User-Stories für die Hands-on-Aufgaben liegen in `docs/stories/`
(`story-01` bis `story-07`).

## Projektstruktur

- `docs/`: Workshop-Dokumentation
  - `themen.md`: Themenübersicht und Moderations-Leitfaden
  - `vorbereitung.md`: Setup der Arbeitsumgebung
  - `stories/`: User-Stories
  - `instructions/`: Trainer-Hinweise
  - `questions/`: Diskussionsfragen
- `services/`: Go-Referenz-Implementierung
  - `booking/`: BookingService mit einem Ordner pro Story (`story1` bis `story7`)
  - `flight/`, `hotel/`, `car/`: Domain-Services
  - `dashboard/`: Dashboard-UI
  - `traefik/`: API-Gateway-Konfiguration
  - `shared/`: gemeinsame Bibliotheken
  - `docker-compose*.yml`, `Makefile`: lokales Setup
