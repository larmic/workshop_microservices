# Microservices Workshop

[![Build Go-Services](https://github.com/larmic/workshop_microservices/actions/workflows/build.yml/badge.svg)](https://github.com/larmic/workshop_microservices/actions/workflows/build.yml)
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

Die Anleitung zur Einrichtung des Arbeitsplatzes findet sich in [docs/vorbereitung.md](docs/vorbereitung.md).

## Projektstruktur

```
├── docs/                  # Workshop-Dokumentation
│   ├── vorbereitung.md    # Arbeitsplatz einrichten
│   ├── themen.md          # Themenübersicht & Moderationsleitfaden
│   ├── stories/           # User Stories (story-01 … story-07)
│   ├── instructions/      # Trainer-Hinweise
│   └── questions/         # Fragen & Diskussionsimpulse
├── services/              # Backend-Services (Go)
│   ├── booking/           # BookingService (story1 … story7)
│   ├── flight/            # FlightService
│   ├── hotel/             # HotelService
│   ├── car/               # CarService
│   ├── dashboard/         # Dashboard-UI
│   ├── traefik/           # API-Gateway-Konfiguration
│   ├── shared/            # Gemeinsame Bibliotheken
│   ├── docker-compose.yml # Lokales Setup
│   └── Makefile
└── README.md
```

## Lizenz & Nutzung

Dieses Repository ist **source-available, aber nicht klassisch Open Source**. Lesen, klonen und persönliche Nutzung sind ausdrücklich erwünscht — kommerzielle Nutzung erfordert eine separate Vereinbarung.

| Inhalt | Lizenz |
|---|---|
| Quellcode in `services/` | [PolyForm Noncommercial 1.0.0](LICENSE) |
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
