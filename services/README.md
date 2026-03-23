# Services

Microservices für das Travel-Booking-System (Flight, Hotel, Car, Booking).

## Service-Übersicht

| Service         | Port | API                                              |
|-----------------|------|--------------------------------------------------|
| Flight          | 8081 | [openapi.yaml](flight/api/openapi.yaml)          |
| Hotel           | 8082 | [openapi.yaml](hotel/api/openapi.yaml)           |
| Car             | 8083 | [openapi.yaml](car/api/openapi.yaml)             |
| Booking Story 1 | 8085 | [openapi.yaml](booking/story1/api/openapi.yaml)  |
| Booking Story 2 | 8086 | [openapi.yaml](booking/story2/api/openapi.yaml)  |

## Infrastruktur

| Service    | Port | Beschreibung                          |
|------------|------|---------------------------------------|
| Consul     | 8500 | Service Discovery & Health Checks     |
| Swagger UI | 8084 | API-Dokumentation (Basis-Services)    |

## Docker Compose Aufteilung

| Datei                        | Inhalt                                        |
|------------------------------|-----------------------------------------------|
| `docker-compose.yml`         | Basis-Services (Flight, Hotel, Car)            |
| `docker-compose.infra.yml`   | Infrastruktur (Consul, Swagger UI)             |
| `docker-compose.reference.yml` | Referenzlösungen (Booking Story 1, Story 2, …) |

## Quickstart

```bash
# Alles starten (und zuvor bauen)
make docker-up

# Alles starten, aber Images von Docker Hub laden
make docker-up-hub
```

Alle verfügbaren Targets: `make help`
