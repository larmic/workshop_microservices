# Workshop-Vorbereitung

Diese Anleitung beschreibt, wie du deinen Arbeitsplatz für den Microservices-Workshop einrichtest.

## Voraussetzungen

Folgende Tools müssen auf deinem Rechner installiert sein:

- **Git**
- **Docker** (inkl. Docker Compose)
- **IDE / Editor** nach Wahl (z.B. IntelliJ IDEA, VS Code)
- **Framework** deiner Wahl lokal lauffähig (z.B. Spring Boot, Quarkus, Go, ...)

## 1. Repository klonen

```bash
git clone https://github.com/larmic/workshop_microservices.git
cd workshop_microservices
```

## 2. Backend-Services starten

Die drei Backend-Services (Flight, Hotel, Car) und eine Swagger UI werden per Docker Compose gestartet:

```bash
docker compose -f services/docker-compose.yml up -d
```

Nach dem Start laufen folgende Services:

| Service    | URL                        | Beschreibung              |
|------------|----------------------------|---------------------------|
| Flight     | http://localhost:8081       | Flugbuchungen             |
| Hotel      | http://localhost:8082       | Hotelbuchungen            |
| Car        | http://localhost:8083       | Mietwagenbuchungen        |
| Swagger UI | http://localhost:8084       | API-Dokumentation         |

## 3. Prüfen, ob alles läuft

Öffne die Swagger UI im Browser:

> http://localhost:8084

Dort siehst du die API-Dokumentation aller drei Backend-Services und kannst die Endpoints direkt ausprobieren.

Alternativ per Terminal:

```bash
# Health Checks
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health

# Beispiel: Verfügbare Flüge abfragen
curl http://localhost:8081/flights
```

## 4. Los geht's

Dein Arbeitsplatz ist eingerichtet. Starte jetzt mit [Story 1: Der erste cloud-native Booking-Service](stories/story-01-cloud-native-booking-service.md).
