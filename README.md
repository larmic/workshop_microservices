# Microservices richtig gemacht

[![Build Go-Services](https://github.com/larmic/workshop_microservices/actions/workflows/build.yml/badge.svg)](https://github.com/larmic/workshop_microservices/actions/workflows/build.yml)
[![Guard-Gates](https://github.com/larmic/workshop_microservices/actions/workflows/gates.yml/badge.svg)](https://github.com/larmic/workshop_microservices/actions/workflows/gates.yml)
[![Docker Build & Push Go-Services](https://github.com/larmic/workshop_microservices/actions/workflows/docker.yml/badge.svg)](https://github.com/larmic/workshop_microservices/actions/workflows/docker.yml)

Microservices sind auch nur Monolithen mit Netzwerkproblemen.

Ein 2-tägiger Workshop zu Microservice-Architektur-Patterns für (angehende) Software-Architekten.

## Zielgruppe

Dieser Workshop richtet sich an:

- Software-Architekten
- Entwickler, die Architekten werden möchten
- Tech Leads mit Architekturverantwortung

## Inhalt

Der Workshop behandelt die wichtigsten Patterns und Best Practices für Microservice-Architekturen. Anhand einer praktischen Beispieldomäne (Reisebuchung: Hotel, Flug, Mietwagen) werden die Konzepte hands-on erarbeitet.

Die vollständige Themenübersicht findet sich in [docs/themen.md](docs/themen.md).

## Format

- **Dauer:** 2 Tage
- **Mix aus:** Kurzvorträgen und Gruppenarbeit
- **Schwerpunkt:** Praktische Hands-On-Übungen
- **Aufgabenformat:** User Stories

## Einstieg

Die Anleitung zur Einrichtung des Arbeitsplatzes findet sich in [docs/vorbereitung.md](docs/vorbereitung.md), online lesbar unter <https://larmic.github.io/workshop_microservices/vorbereitung/>. Die Workshop-Seite mit Folien und Feedback: <https://larmic.github.io/workshop_microservices/>.

## Projektstruktur

```
├── docs/                  # Workshop-Dokumentation
│   ├── vorbereitung.md    # Pflicht-Vorbereitung vor dem Workshop
│   ├── custom-setup.md    # Eigenen Booking-Service einklinken (im Workshop)
│   ├── themen.md          # Themenübersicht & Moderationsleitfaden
│   ├── stories/           # User Stories (story-01 … story-08)
│   ├── instructions/      # Trainer-Hinweise
│   ├── questions/         # Fragen & Diskussionsimpulse
│   └── orga/              # Vorlagen für Intranet-Text & Einladungs-Mail
├── pages/                 # Statische Quellen der GitHub-Pages-Seite
│   ├── landing/           # Startseite (larmic.github.io/workshop_microservices)
│   ├── vorbereitung/      # Vorbereitung & Troubleshooting (rendert docs/*.md)
│   └── feedback/          # Anonymes Feedback-Formular
├── services/              # Backend-Services (Go)
│   ├── booking/           # BookingService (story1, story3 … story8; Story 2 hat keinen Code)
│   ├── flight/            # FlightService
│   ├── hotel/             # HotelService
│   ├── car/               # CarService
│   ├── dashboard/         # Dashboard-UI
│   ├── traefik/           # API-Gateway-Konfiguration
│   ├── shared/            # Gemeinsame Bibliotheken
│   ├── archtest/          # Guard-Gates als Go-Tests (Architektur, OpenAPI, Farbpalette)
│   ├── docker-compose.yml # Lokales Setup
│   └── Makefile
├── scripts/gates.sh       # Alle Guard-Gates, lokal wie in der Pipeline
├── .githooks/pre-commit   # Schnelle Guard-Gates vor jedem Commit
└── README.md
```

## Guard-Gates

Die Referenz-Implementierung ist Lehrmaterial: Was hier steht, übernehmen
Teilnehmende als Vorbild. Deshalb läuft jeder Pull Request und jeder Push auf
`main` durch dieselben deterministischen Prüfungen, egal ob ein Mensch oder
eine KI den Code geschrieben hat
([.github/workflows/gates.yml](.github/workflows/gates.yml), ein Schritt pro Gate).

| Gate | Prüft | Werkzeug |
| --- | --- | --- |
| `gofmt` | einheitlich formatierter Go-Code | gofmt |
| `vet` | typische Fehler im Go-Code | go vet |
| `arch` | Import-Graph: jede Booking-Story, Flight, Hotel, Car und das Dashboard importieren nur sich selbst und `shared`; `shared` kennt keinen Service | Go-Test in `services/archtest` |
| `openapi` | jede in `api/openapi.yaml` deklarierte Route ist registriert und umgekehrt; `/health`, `/info`, `/openapi` gibt es überall | Go-Test in `services/archtest` |
| `palette` | Farbrollen aus `services/CLAUDE.md` sind in `slides/theme.css` und im Dashboard vorhanden | Go-Test in `services/archtest` |
| `test` | Unit-Tests | go test |
| `js` | Syntax der Inline-Skripte von Dashboard, Folien und Pages | `node --check` |
| `html` | gültiges HTML von Dashboard, Folien und Pages | html-validate (Preset `standard`) |
| `compose` | die Compose-Kombinationen aus `Makefile` und Vorbereitung | `docker compose config -q` |

Lokal laufen alle zusammen mit `scripts/gates.sh` (oder `make -C services gates`),
einzelne mit `scripts/gates.sh openapi`. Ein Git-Hook (`.githooks/pre-commit`)
prüft vor jedem Commit die schnellen Gates (gofmt, vet, Architektur, OpenAPI,
Palette, Inline-JavaScript) in wenigen Sekunden. Aktiviert wird er einmal pro
Klon mit `git config core.hooksPath .githooks` (oder `make -C services hooks`);
im Notfall überspringt `git commit --no-verify` ihn. Fehlt lokal ein Werkzeug
(Go, Node, Docker Compose), überspringen Skript und Hook das Gate mit Hinweis;
verbindlich ist die Pipeline.

## Lizenz & Nutzung

Dieses Repository ist **source-available, aber nicht klassisch Open Source**. Lesen, klonen und persönliche Nutzung sind ausdrücklich erwünscht — kommerzielle Nutzung erfordert eine separate Vereinbarung.

| Inhalt | Lizenz |
|---|---|
| Quellcode in `services/` und `pages/` | [PolyForm Noncommercial 1.0.0](LICENSE) |
| Workshop-Inhalt in `docs/` | [CC BY-NC 4.0](docs/LICENSE) |

**Erlaubt (ohne Rückfrage):**

- Code und Doku lesen, klonen, forken
- Persönliche Nutzung zum Lernen, Experimentieren und Nachvollziehen — etwa als ehemaliger Workshop-Teilnehmer
- Nutzung im Rahmen von Bildungs- oder Forschungseinrichtungen

**Erfordert eine schriftliche Genehmigung:**

- Durchführung dieses Workshops (oder wesentlicher Teile davon) als bezahltes Angebot
- Integration substantieller Teile des Codes in kommerzielle Produkte oder Services
- Wiederverwendung der Workshop-Materialien in kommerziellen Trainings

Für kommerzielle Nutzung: Lars Michaelis &lt;l.michaelis@neusta.de&gt;
