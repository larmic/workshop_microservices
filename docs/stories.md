# Workshop Stories - Reisebuchungsplattform

## Worum geht es?

In diesem Workshop bauen die Teilnehmer Schritt für Schritt eine Reisebuchungsplattform auf Basis von Microservices. Kunden können Flüge, Hotels und Mietwagen buchen — der zentrale **BookingService** orchestriert dabei die spezialisierten Backend-Services. Entlang von User Stories werden die wichtigsten Microservice-Patterns praktisch erarbeitet: von Health-Checks über Circuit Breaker bis hin zu Event-Driven Architecture.

## Systemlandschaft

```
                         ┌─────────────┐
                         │   Client    │
                         │ (Web/Mobile)│
                         └──────┬──────┘
                                │
                         ┌──────▼──────┐
                         │ API Gateway │
                         │  (Port 8080)│
                         └──────┬──────┘
                                │
                    ┌───────────▼───────────┐
                    │   BookingService      │
                    │   (Teilnehmer bauen)  │◄────► Consul
                    └───┬───────┬───────┬───┘      (Service
                        │       │       │          Registry)
              ┌─────────▼┐ ┌───▼─────┐ ┌▼─────────┐
              │  Flight   │ │  Hotel  │ │   Car    │
              │  Service  │ │ Service │ │  Service │
              │ (Port 8081│ │(Port 8082│ │(Port 8083│
              └───────────┘ └─────────┘ └──────────┘
```

## Workshop-Setup

**Bereitgestellte Backend-Services (Docker-Container vom Trainer):**
- `FlightService` - Verwaltet Flugbuchungen
- `HotelService` - Verwaltet Hotelbuchungen
- `CarService` - Verwaltet Mietwagenbuchungen

**Was die Teilnehmer bauen:**
- Die Buchungs-/Bestellstrecke (Orchestrierung der Backend-Services)

**Technische Rahmenbedingungen:**
- Framework: Frei wählbar
- Service Discovery: Consul (Docker-Container wird bereitgestellt)
- Messaging: REST-Webhooks (kein separater Message-Broker erforderlich)

Die Schritt-für-Schritt-Anleitung zur Einrichtung des Arbeitsplatzes findet sich in der [Vorbereitung](vorbereitung.md).

---

## Story-Index

### Tag 1: Grundlagen & Kommunikation

| # | Story | Thema | Zeitrahmen |
|---|-------|-------|------------|
| 1 | [Der erste cloud-native Booking-Service](stories/story-01-cloud-native-booking-service.md) | Twelve-Factor-App, Health-Checks, Externe Konfiguration — der BookingService wird aufgesetzt und ruft erstmals die Backend-Services auf. | ca. 90 Min. |
| 2 | [Services dynamisch finden](stories/story-02-service-discovery.md) | Service Discovery mit Consul — statische URLs werden durch dynamische Service-Registrierung ersetzt. | ca. 60 Min. |

### Tag 2: Resilience & Events

| # | Story | Thema | Zeitrahmen |
|---|-------|-------|------------|
| 3 | [Wenn der Flug ausfällt](stories/story-03-circuit-breaker.md) | Circuit Breaker — graceful Degradation bei Ausfall des FlightService. | ca. 60 Min. |
| 4 | [Isolation ist Stärke](stories/story-04-bulkhead.md) | Bulkhead Pattern — Ressourcen-Isolation verhindert Kaskadenausfälle. | ca. 60 Min. |
| 5 | [Alles oder nichts - aber richtig](stories/story-05-saga.md) | Saga Pattern — verteilte Transaktionen mit Kompensationslogik für Komplettbuchungen. | ca. 60 Min. |
| 6 | [Die Saga wird leise](stories/story-06-choreography-saga.md) | Choreography-Saga — Kompensation wandert vom Booking-Service in die Backends via Events. | ca. 60 Min. |

### Optional (bei Zeit)

| # | Story | Thema | Zeitrahmen |
|---|-------|-------|------------|
| 8 | [Mobile First](stories/story-08-bff.md) | Backends for Frontends — ein spezialisiertes Backend für die mobile App. | ca. 60 Min. |
| 9 | [Ohne Unterbrechung](stories/story-09-downtimeless-deployment.md) | Downtimeless Deployment — Rolling Updates ohne Ausfallzeit. | ca. 60 Min. |
| 10 | [Den roten Faden im Log](stories/story-10-tracing.md) | Distributed Tracing — Geschäftsvorgänge über Service-Grenzen hinweg verfolgen. | ca. 60 Min. |

---

## Anhang

### Service-API-Dokumentation

#### FlightService (Port 8081)

| Methode | Endpoint | Beschreibung |
|---------|----------|--------------|
| GET | `/flights` | Liste verfügbarer Flüge |
| GET | `/flights/{id}` | Flugdetails |
| POST | `/bookings` | Flug buchen |
| GET | `/bookings/{id}` | Buchungsdetails |
| DELETE | `/bookings/{id}` | Buchung stornieren |
| GET | `/health` | Health-Check |

#### HotelService (Port 8082)

| Methode | Endpoint | Beschreibung |
|---------|----------|--------------|
| GET | `/hotels` | Liste verfügbarer Hotels |
| GET | `/hotels/{id}` | Hoteldetails |
| POST | `/bookings` | Hotel buchen |
| GET | `/bookings/{id}` | Buchungsdetails |
| DELETE | `/bookings/{id}` | Buchung stornieren |
| GET | `/health` | Health-Check |

#### CarService (Port 8083)

| Methode | Endpoint | Beschreibung |
|---------|----------|--------------|
| GET | `/cars` | Liste verfügbarer Mietwagen |
| GET | `/cars/{id}` | Mietwagendetails |
| POST | `/bookings` | Mietwagen buchen |
| GET | `/bookings/{id}` | Buchungsdetails |
| DELETE | `/bookings/{id}` | Buchung stornieren |
| GET | `/health` | Health-Check |

### Story-Abhängigkeiten

```
Story 1 (Grundsetup + Konfiguration)
    │
    ├── Story 2 (Service Discovery)
    │
    └── Story 3 (Circuit Breaker)
            │
            ├── Story 4 (Bulkhead)
            │
            └── Story 5 (Saga, sync Kompensation)
                    │
                    └── Story 6 (Choreography-Saga, async Kompensation via Events)

Optional:
    Story 1 + Story 2 → Story 8 (BFF)
    Story 1 + Story 2 → Story 9 (Downtimeless)
    Story 5 + Story 6 → Story 10 (Distributed Tracing)
```

### Glossar

| Begriff | Erklärung |
|---------|-----------|
| Circuit Breaker | Schutzmechanismus, der bei wiederholten Fehlern weitere Aufrufe unterbricht |
| Bulkhead | Isolierung von Ressourcen, um Kaskadenausfälle zu verhindern |
| Saga | Pattern für verteilte Transaktionen mit Kompensationslogik |
| Orchestration-Saga | Zentraler Koordinator (Booking) steuert Forward- und Kompensationsschritte synchron |
| Choreography-Saga | Services reagieren auf Events anderer Services — Saga-Wissen ist verteilt |
| CQRS | Trennung von Lese- und Schreibmodellen zur unabhängigen Optimierung — kann, muss aber nicht mit Events kombiniert werden |
| BFF | Backend for Frontend - spezialisiertes Backend pro Client-Typ |
| Eventual Consistency | Daten werden irgendwann konsistent, nicht sofort |
| Graceful Shutdown | Kontrolliertes Herunterfahren ohne Abbruch laufender Operationen |
