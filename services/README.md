# Services

Dokumentation der Microservices für das Travel-Booking-System.

## Service-Übersicht

| Service    | Port (Docker) | Interner Port | Beschreibung                    |
|------------|---------------|---------------|---------------------------------|
| Flight     | 8081          | 8080          | Flugbuchungen                   |
| Hotel      | 8082          | 8080          | Hotelbuchungen                  |
| Car        | 8083          | 8080          | Mietwagenbuchungen              |
| Swagger UI | 8084          | 8080          | API-Dokumentation               |
| **Booking**| **8085**      | **8080**      | **Orchestrierungs-Service**     |

## API-Endpoints

### Basis-Services (Flight, Hotel, Car)

Alle Basis-Services haben eine identische Endpoint-Struktur:

| Endpoint            | Beschreibung                          |
|---------------------|---------------------------------------|
| `GET /health`       | Health Check                          |
| `GET /flights`      | Liste aller Flüge                     |
| `GET /hotels`       | Liste aller Hotels                    |
| `GET /cars`         | Liste aller Mietwagen                 |
| `GET /openapi`      | OpenAPI-Spezifikation                 |
| `POST /bookings`    | Buchung erstellen                     |
| `DELETE /bookings/{id}` | Buchung stornieren (Compensating Transaction) |

### Booking Service (Orchestrator)

Der Booking Service koordiniert Buchungen über alle drei Basis-Services:

| Methode  | Endpoint           | Beschreibung                |
|----------|--------------------|-----------------------------|
| `POST`   | `/bookings`        | Neue Buchung erstellen      |
| `GET`    | `/bookings`        | Alle Buchungen auflisten    |
| `GET`    | `/bookings/{id}`   | Einzelne Buchung abrufen    |
| `DELETE` | `/bookings/{id}`   | Buchung stornieren          |
| `GET`    | `/health`          | Health Check                |
| `GET`    | `/openapi`         | OpenAPI-Spezifikation       |

## Testdaten

Die Services liefern konsistente Testdaten für drei Reiseziele:

| Flug | Hotel | Mietwagen |
|------|-------|-----------|
| LH123: Frankfurt → New York (450 EUR) | Manhattan Plaza Hotel (250 USD/night) | Ford Mustang (95 USD/day) |
| LH456: München → London (180 EUR) | London Bridge Inn (180 GBP/night) | Mini Cooper (65 GBP/day) |
| BA789: Berlin → Paris (120 EUR) | Paris Étoile Hotel (195 EUR/night) | Renault Clio (55 EUR/day) |

## Quickstart

### Services starten

```bash
# Alle Services mit Docker Compose starten
docker compose up -d

# Oder einzelne Services
docker compose up -d flight hotel car
```

### Services testen

```bash
# Flüge abfragen
curl http://localhost:8081/flights

# Hotels abfragen
curl http://localhost:8082/hotels

# Mietwagen abfragen
curl http://localhost:8083/cars

# Health Checks
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health
```

### Swagger UI

Die API-Dokumentation ist unter http://localhost:8084 verfügbar.

## Booking Service - Detaillierte Dokumentation

### Architektur

Der Booking Service fungiert als Orchestrator und koordiniert Buchungen über die drei Basis-Services:

```
┌─────────┐     ┌─────────────────┐     ┌──────────┐
│ Client  │────▶│ Booking Service │────▶│ Flight   │
└─────────┘     │                 │     └──────────┘
                │  (Orchestrator) │────▶┌──────────┐
                │                 │     │ Hotel    │
                │                 │     └──────────┘
                └─────────────────┘────▶┌──────────┐
                                        │ Car      │
                                        └──────────┘
```

### Buchungsablauf (User Journey)

```
1. Kunde sucht Flug      → GET /flights         (Flight Service)
2. Kunde sucht Hotel     → GET /hotels          (Hotel Service)
3. Kunde sucht Mietwagen → GET /cars            (Car Service)
4. Kunde bucht alles     → POST /bookings       (Booking Service)
5. Kunde prüft Buchung   → GET /bookings/{id}   (Booking Service)
6. Kunde storniert       → DELETE /bookings/{id} (Booking Service)
```

### Beispiel: Buchung erstellen

```bash
curl -X POST http://localhost:8085/bookings \
  -H "Content-Type: application/json" \
  -d '{
    "customerName": "Max Mustermann",
    "customerEmail": "max@example.com",
    "flightId": "LH123",
    "hotelId": "H1",
    "carId": "C1",
    "checkIn": "2024-06-01",
    "checkOut": "2024-06-05"
  }'
```

**Hinweis:** Alle Felder außer `customerName` sind optional. Man kann z.B. nur Flug + Hotel buchen.

### SAGA-Pattern (All-or-Nothing)

Der Booking Service implementiert das SAGA-Pattern für verteilte Transaktionen. Im Gegensatz zu "Partial Success" gilt hier das **All-or-Nothing-Prinzip**:

| Szenario | Verhalten |
|----------|-----------|
| Alle Teilbuchungen erfolgreich | Status: `CONFIRMED`, HTTP 201 |
| Eine Teilbuchung fehlgeschlagen | Status: `FAILED`, HTTP 409, automatischer Rollback |

**Workshop-Hinweis:** Die Teilnehmer implementieren den Booking Service selbst und lernen dabei das SAGA-Pattern.

### Buchungsstatus (Gesamtbuchung)

| Status      | Beschreibung                                           |
|-------------|--------------------------------------------------------|
| `PENDING`   | Buchung wird verarbeitet                               |
| `CONFIRMED` | Alle Teilbuchungen erfolgreich bestätigt               |
| `FAILED`    | Mindestens eine Teilbuchung fehlgeschlagen, Rollback durchgeführt |
| `CANCELLED` | Buchung vom Kunden storniert                           |

### Teilbuchungs-Status (Flight, Hotel, Car)

| Status          | Beschreibung                                      |
|-----------------|---------------------------------------------------|
| `CONFIRMED`     | Erfolgreich gebucht                               |
| `FAILED`        | Buchung fehlgeschlagen                            |
| `ROLLED_BACK`   | War gebucht, wurde durch SAGA-Rollback storniert  |
| `NOT_ATTEMPTED` | Wurde nicht versucht (wegen vorherigem Fehler)    |
| `CANCELLED`     | Vom Kunden storniert                              |

### SAGA-Workflow Beispiel

**Szenario: Hotel nicht verfügbar**

```
1. Client: POST /bookings {flight: LH123, hotel: H1, car: C1}
2. Booking Service:
   a) POST flight/bookings  → 201 OK (FB-001)
   b) POST hotel/bookings   → 409 Conflict (keine Zimmer)
   c) DELETE flight/bookings/FB-001 (Compensating Transaction)
   d) Car wird nicht mehr versucht
3. Response: 409 Conflict
```

**Response bei SAGA-Rollback:**

```json
{
  "id": "BK-20240601-002",
  "status": "FAILED",
  "error": "Hotel H1 not available for requested dates",
  "flight": {
    "status": "ROLLED_BACK",
    "bookingRef": "FB-001",
    "note": "Booking was created and then cancelled due to saga rollback"
  },
  "hotel": {
    "status": "FAILED",
    "error": "No rooms available"
  },
  "car": {
    "status": "NOT_ATTEMPTED",
    "note": "Booking was not attempted due to previous failure"
  }
}
```

### Konfiguration

Der Booking Service benötigt die URLs der Basis-Services:

```yaml
# docker-compose.yml
booking:
  build: ./booking
  ports:
    - "8085:8080"
  environment:
    - FLIGHT_SERVICE_URL=http://flight:8080
    - HOTEL_SERVICE_URL=http://hotel:8080
    - CAR_SERVICE_URL=http://car:8080
```
