# Services

[![Build Services](https://github.com/larmic/workshop_microservices/actions/workflows/build.yml/badge.svg)](https://github.com/larmic/workshop_microservices/actions/workflows/build.yml)
[![Docker Build & Push Go-Services](https://github.com/larmic/workshop_microservices/actions/workflows/docker.yml/badge.svg)](https://github.com/larmic/workshop_microservices/actions/workflows/docker.yml)
[![Deploy Pages (Feedback + Slides)](https://github.com/larmic/workshop_microservices/actions/workflows/deploy-pages.yml/badge.svg)](https://github.com/larmic/workshop_microservices/actions/workflows/deploy-pages.yml)
[![Update Docker Hub descriptions](https://github.com/larmic/workshop_microservices/actions/workflows/dockerhub-description.yml/badge.svg)](https://github.com/larmic/workshop_microservices/actions/workflows/dockerhub-description.yml)

Microservices für das Travel-Booking-System (Flight, Hotel, Car, Booking).

Die Badges zeigen den Status des jeweils letzten Laufs jeder Pipeline. Die Images
werden im Workflow „Docker Build & Push Go-Services" gebaut und nach Docker Hub
gepusht.

## Docker Images

| Image | Tag | Größe | Docker Hub |
|-------|-----|-------|------------|
| Flight    | `latest` | ![size](https://img.shields.io/docker/image-size/larmic/workshop-microservices-flight/latest?label=) | [larmic/workshop-microservices-flight](https://hub.docker.com/r/larmic/workshop-microservices-flight) |
| Hotel     | `latest` | ![size](https://img.shields.io/docker/image-size/larmic/workshop-microservices-hotel/latest?label=) | [larmic/workshop-microservices-hotel](https://hub.docker.com/r/larmic/workshop-microservices-hotel) |
| Car       | `latest` | ![size](https://img.shields.io/docker/image-size/larmic/workshop-microservices-car/latest?label=) | [larmic/workshop-microservices-car](https://hub.docker.com/r/larmic/workshop-microservices-car) |
| Booking   | `story1`…`story7`, `custom` | ![size](https://img.shields.io/docker/image-size/larmic/workshop-microservices-booking/story1?label=story1) | [larmic/workshop-microservices-booking](https://hub.docker.com/r/larmic/workshop-microservices-booking) |
| Dashboard | `latest` | ![size](https://img.shields.io/docker/image-size/larmic/workshop-microservices-dashboard/latest?label=) | [larmic/workshop-microservices-dashboard](https://hub.docker.com/r/larmic/workshop-microservices-dashboard) |

## Service-Übersicht

| Service         | Traefik-Route              | API                                              |
|-----------------|----------------------------|--------------------------------------------------|
| Flight          | `/api/flight/**`           | [openapi.yaml](flight/api/openapi.yaml)          |
| Hotel           | `/api/hotel/**`            | [openapi.yaml](hotel/api/openapi.yaml)           |
| Car             | `/api/car/**`              | [openapi.yaml](car/api/openapi.yaml)             |
| Booking Reference Story 1 | `/api/booking-ref-story1/**` | [openapi.yaml](booking/story1/api/openapi.yaml)  |
| Booking Reference Story 2 | `/api/booking-ref-story2/**` | [openapi.yaml](booking/story2/api/openapi.yaml)  |
| Booking Custom            | `/api/booking-custom/**`     | [openapi.yaml](booking/custom/src/main/resources/openapi.yaml) |

Alle Services sind über den Traefik Reverse Proxy auf Port 80 erreichbar:

```bash
curl http://localhost/api/flight/flights
curl http://localhost/api/hotel/hotels
curl http://localhost/api/car/cars
curl http://localhost/api/booking-ref-story1/booking/offers
curl http://localhost/api/booking-ref-story2/booking/offers
curl http://localhost/api/booking-custom/booking/offers
```

## Infrastruktur

| Service           | URL                          | Beschreibung                      |
|-------------------|------------------------------|-----------------------------------|
| Dashboard         | http://localhost              | Workshop Dashboard                |
| Traefik Dashboard | http://localhost:8080         | Traefik Monitoring Dashboard      |
| Consul            | http://localhost/consul       | Service Discovery & Health Checks |
| Swagger UI        | http://localhost/api          | API-Dokumentation (via Traefik)   |

## Docker Compose Aufteilung

| Datei                          | Inhalt                                          |
|--------------------------------|-------------------------------------------------|
| `docker-compose.yml`           | Basis-Services (Flight, Hotel, Car)              |
| `docker-compose.infra.yml`     | Infrastruktur (Traefik, Consul, Swagger UI)      |
| `docker-compose.reference.yml` | Referenzlösungen (Booking Reference Story 1…7)   |
| `docker-compose.custom.yml`    | Custom-Lösung des Teilnehmers (Booking Custom)   |

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
