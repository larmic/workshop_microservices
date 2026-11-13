# Microservices-Workshop

Workshop-Projekt zur Vermittlung von Architekturmustern und Best Practices für
Microservices. Der Workshop heißt offiziell **"Microservices richtig gemacht"**
(so in Titeln, Pages und Orga-Texten zu verwenden). Sprache des gesamten
Workshop-Materials: Deutsch.

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
(`story-01` bis `story-08`).

## Projektstruktur

- `docs/`: Workshop-Dokumentation
  - `themen.md`: Themenübersicht und Moderations-Leitfaden
  - `vorbereitung.md`: Setup der Arbeitsumgebung
  - `stories/`: User-Stories
  - `instructions/`: Trainer-Hinweise
  - `questions/`: Diskussionsfragen
  - `orga/`: Vorlagen für Intranet-Text und E-Mail-Einladung
- `pages/`: statische Quellen der GitHub-Pages-Seite (Landing, Vorbereitung,
  Feedback-Formular; die Folien kommen aus `services/slides/`)
- `services/`: Go-Referenz-Implementierung
  - `booking/`: BookingService mit einem Ordner pro Code-Story (`story1`, `story3`
    bis `story8`; Story 2 ist eine Design-Session ohne Code, Port 8086 bleibt frei)
  - `flight/`, `hotel/`, `car/`: Domain-Services
  - `dashboard/`: Dashboard-UI
  - `traefik/`: API-Gateway-Konfiguration
  - `shared/`: gemeinsame Bibliotheken
  - `docker-compose*.yml`, `Makefile`: lokales Setup

## Guard-Gates

```bash
scripts/gates.sh            # alle Gates in der Reihenfolge der Pipeline (auch: make -C services gates)
scripts/gates.sh fast       # schnelle Gates wie im pre-commit-Hook
scripts/gates.sh gofmt      # Guard-Gate: Go-Formatierung (cd services && gofmt -w . behebt es)
scripts/gates.sh vet        # Guard-Gate: go vet ./...
scripts/gates.sh arch       # Guard-Gate: Architekturregeln (services/archtest)
scripts/gates.sh openapi    # Guard-Gate: api/openapi.yaml ⇄ registrierte Routen
scripts/gates.sh palette    # Guard-Gate: Farbpalette slides/theme.css ⇄ dashboard/static/index.html, keine fremden Farben in Dashboard und pages/
scripts/gates.sh test       # Guard-Gate: go test (ohne archtest)
scripts/gates.sh js         # Guard-Gate: Syntax der Inline-Skripte (node --check)
scripts/gates.sh unicode    # Guard-Gate: keine unsichtbaren Unicode-Zeichen (gewollte als Escape, z. B. \u200B)
scripts/gates.sh html       # Guard-Gate: html-validate für Dashboard, Folien, Pages
scripts/gates.sh compose    # Guard-Gate: docker compose config -q für die Makefile-Kombinationen
# Git-Hook .githooks/pre-commit (aktivieren: git config core.hooksPath .githooks) prüft gofmt, vet, archtest, JS, Unicode
```

Alle Gates müssen grün sein; die Pipeline `.github/workflows/gates.yml` führt
jedes als eigenen Schritt aus. Vor jedem Commit mindestens `scripts/gates.sh fast`
ausführen, nach Änderungen an Compose-Files, HTML oder Specs `scripts/gates.sh`.

Regeln, die die Gates durchsetzen:

- Jede Booking-Story (`booking/storyN`), `flight`, `hotel`, `car` und
  `dashboard` importiert nur sich selbst und `shared`; `shared` importiert
  keinen Service. Neue Stories fallen automatisch unter die Regel.
- Jede Route in `main.go`/`handler/` steht in `api/openapi.yaml` und umgekehrt,
  Muster immer mit Methode (`"GET /pfad"`). Neue Routen also immer mit Spec.
- Farbänderungen an Slides, Dashboard oder Pages laufen gegen die Palette in
  `services/CLAUDE.md` (`voxel.css`, das Easter-Egg-Theme, ist ausgenommen).
