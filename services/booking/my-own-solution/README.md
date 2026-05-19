# Booking — My Own Solution (Kotlin / Ktor)

Eigene Umsetzung der **Story 1** in Kotlin mit [Ktor](https://ktor.io/).
Funktional identisch zur Go-Referenz unter `services/booking/story1/`.

## Eigenschaften

- **Sprache:** Kotlin 2.1 (JVM 21)
- **Framework:** Ktor 3 (CIO Server + CIO Client)
- **Startzeit:** ~0,5 s
- **Image-Größe:** ~120 MB (`eclipse-temurin:21-jre-alpine` + Fat-Jar)
- **Build:** Gradle (Shadow-Plugin für Fat-Jar)

## Endpoints

| Endpoint | Beschreibung |
|---|---|
| `GET /health` | Healthcheck, antwortet `{"status":"UP"}` |
| `GET /info` | Aktuelle Konfiguration als JSON |
| `GET /openapi` | Eingebettete OpenAPI-Spezifikation (YAML) |
| `GET /booking/offers` | Ruft Flight-, Hotel- und Car-Service parallel auf und kombiniert die Ergebnisse |

## Environment-Variablen

| Variable | Default |
|---|---|
| `FLIGHT_SERVICE_URL` | `http://localhost:8081` |
| `HOTEL_SERVICE_URL`  | `http://localhost:8082` |
| `CAR_SERVICE_URL`    | `http://localhost:8083` |

HTTP-Client-Timeout: 5000 ms (hardcoded, analog zur Go-Referenz).

## Build & Run

### Mit Docker (empfohlen)

```sh
docker build -t booking-my-own-solution .
docker run --rm -p 8080:8080 \
  -e FLIGHT_SERVICE_URL=http://host.docker.internal:8081 \
  -e HOTEL_SERVICE_URL=http://host.docker.internal:8082 \
  -e CAR_SERVICE_URL=http://host.docker.internal:8083 \
  booking-my-own-solution
```

### Lokal (mit Gradle)

```sh
gradle shadowJar
java -jar build/libs/booking-my-own-solution-all.jar
```

## Smoke-Test

```sh
curl -s localhost:8080/health
curl -s localhost:8080/info
curl -s localhost:8080/booking/offers | jq
```
