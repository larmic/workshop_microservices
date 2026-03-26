# Services

Microservices für das Travel-Booking-System (Flight, Hotel, Car, Booking).

## Service-Übersicht

| Service         | Traefik-Route              | API                                              |
|-----------------|----------------------------|--------------------------------------------------|
| Flight          | `/api/flight/**`           | [openapi.yaml](flight/api/openapi.yaml)          |
| Hotel           | `/api/hotel/**`            | [openapi.yaml](hotel/api/openapi.yaml)           |
| Car             | `/api/car/**`              | [openapi.yaml](car/api/openapi.yaml)             |
| Booking Story 1 | `/api/booking-story1/**`   | [openapi.yaml](booking/story1/api/openapi.yaml)  |
| Booking Story 2 | `/api/booking-story2/**`   | [openapi.yaml](booking/story2/api/openapi.yaml)  |

Alle Services sind über den Traefik Reverse Proxy auf Port 80 erreichbar:

```bash
curl http://localhost/api/flight/flights
curl http://localhost/api/hotel/hotels
curl http://localhost/api/car/cars
curl http://localhost/api/booking-story1/booking/offers
curl http://localhost/api/booking-story2/booking/offers
```

## Infrastruktur

| Service           | Port | Beschreibung                      |
|-------------------|------|-----------------------------------|
| Traefik           | 80   | API-Gateway / Reverse Proxy       |
| Traefik Dashboard | 8080 | Traefik Monitoring Dashboard      |
| Consul            | 8500 | Service Discovery & Health Checks |
| Swagger UI        | 80   | API-Dokumentation (via Traefik)   |

## Docker Compose Aufteilung

| Datei                          | Inhalt                                          |
|--------------------------------|-------------------------------------------------|
| `docker-compose.yml`           | Basis-Services (Flight, Hotel, Car)              |
| `docker-compose.infra.yml`     | Infrastruktur (Traefik, Consul, Swagger UI)      |
| `docker-compose.reference.yml` | Referenzlösungen (Booking Story 1, Story 2, …)   |

## Quickstart

```bash
# Alles starten (und zuvor bauen)
make docker-up

# Alles starten, aber Images von Docker Hub laden
make docker-up-hub
```

Alle verfügbaren Targets: `make help`

## Skalierung

Die Basis-Services (Flight, Hotel, Car) können über `--scale` skaliert werden.
Traefik verteilt die Requests per Round-Robin, und jede Instanz registriert sich
mit einer eigenen ID bei Consul.

```bash
docker compose -f docker-compose.yml -f docker-compose.infra.yml -f docker-compose.reference.yml up --build --scale flight=3
```
