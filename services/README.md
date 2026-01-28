# Services

Dokumentation der Microservices für das Travel-Booking-System.

## Service-Übersicht

| Service    | Port (Docker) | Interner Port | Beschreibung       |
|------------|---------------|---------------|--------------------|
| Flight     | 8081          | 8080          | Flugbuchungen      |
| Hotel      | 8082          | 8080          | Hotelbuchungen     |
| Car        | 8083          | 8080          | Mietwagenbuchungen |
| Swagger UI | 8084          | 8080          | API-Dokumentation  |

## API-Endpoints

Alle Services haben eine identische Endpoint-Struktur:

| Endpoint      | Beschreibung              |
|---------------|---------------------------|
| `GET /health` | Health Check              |
| `GET /flights`| Liste aller Flüge         |
| `GET /hotels` | Liste aller Hotels        |
| `GET /cars`   | Liste aller Mietwagen     |
| `GET /openapi`| OpenAPI-Spezifikation     |

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
