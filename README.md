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

## Format

- **Dauer:** 2 Tage
- **Mix aus:** Kurzvorträgen und Gruppenarbeit
- **Schwerpunkt:** Praktische Hands-On-Übungen
- **Aufgabenformat:** User Stories

## Themen

Die vollständige Themenübersicht findet sich in [docs/themen.md](docs/themen.md).

### Grundlagen
- Twelve-Factor-App
- Health-Checks
- 1 DB pro Service
- External Configuration

### Resilience Patterns
- Circuit Breaker
- Bulkhead Pattern
- Saga Pattern

### Kommunikation & Routing
- API-Gateway
- Service Discovery / Service Registry
- Backends for Frontends (BFF)

### Daten & Events
- Eventsourcing / Event-Driven Architecture
- CQRS

### Deployment & Betrieb
- Downtimeless Deployment
- API First Ansatz

### Kultur & Organisation
- Microservices ohne DevOps?
- "You build it, you run it"
- Cloud Native

## Einstieg

Die Anleitung zur Einrichtung des Arbeitsplatzes findet sich in [docs/vorbereitung.md](docs/vorbereitung.md).

## Projektstruktur

```
├── docs/             # Workshop-Dokumentation
│   ├── vorbereitung.md # Arbeitsplatz einrichten
│   ├── stories.md    # User Stories (Index)
│   ├── stories/      # Einzelne Story-Dateien
│   ├── idea.md       # Workshop-Konzept
│   └── themen.md     # Themenübersicht
├── services/         # Backend-Services (Go)
│   ├── go.mod        # Go-Modul
│   ├── flight/       # FlightService (Port 8081)
│   ├── hotel/        # HotelService (Port 8082)
│   └── car/          # CarService (Port 8083)
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
