---
name: new-booking-story
description: Legt eine neue Booking-Service Story an (Verzeichnisstruktur, Code, OpenAPI, HTTP-Tests, Makefile, GitHub Workflows, Docker Compose)
user_invocable: true
---

# Neue Booking Story anlegen

Erstelle eine neue BookingService Story mit allen zugehoerigen Dateien und Build-Konfigurationen.

## Argumente

Das erste Argument ist die Story-Nummer (z.B. `3` fuer Story 3). Wenn keine Nummer angegeben wurde, frage den User danach.

## Ablauf

### Vorbedingungen pruefen

1. Lese die Story-Nummer N aus den Argumenten
2. Pruefe ob `services/booking/storyN/` bereits existiert. Falls ja: Abbruch mit Fehlermeldung.
3. Berechne den externen Port: `8084 + N` (Story1=8085, Story2=8086, Story3=8087, usw.)

### Neue Dateien erstellen

Lese jeweils die Datei aus `services/booking/story1/` als Template und erstelle die angepasste Version unter `services/booking/storyN/`:

1. **`services/booking/storyN/main.go`**
   - Kopiere 1:1 von story1/main.go (keine Aenderungen noetig)

2. **`services/booking/storyN/api/openapi.yaml`**
   - Kopiere von story1/api/openapi.yaml
   - Aendere `version: 1.0.0` zu `version: N.0.0`

3. **`services/booking/storyN/http/requests.http`**
   - Kopiere von story1/http/requests.http
   - Aendere den Header-Kommentar von "Story 1" zu "Story N"

4. **`services/booking/storyN/http/http-client.env.json`**
   - Kopiere von story1/http/http-client.env.json
   - Aendere `bookingPort` in beiden Environments (local und docker) auf den berechneten Port (8084+N)

5. **`services/booking/storyN/http/Makefile`**
   - Kopiere 1:1 von story1/http/Makefile (keine Aenderungen noetig)

6. **`services/booking/storyN/http/README.md`**
   - Kopiere von story1/http/README.md
   - Aendere "Story 1" zu "Story N" im Titel und Text

### Bestehende Dateien erweitern

7. **`services/Makefile`**
   - Fuege in der `.PHONY`-Zeile `run-booking-storyN` und `docker-build-booking-storyN` hinzu
   - Fuege `docker-build-booking-storyN` als Dependency zum `docker-build`-Target hinzu
   - Fuege ein neues Target nach dem letzten `docker-build-booking-story*`-Target hinzu:
     ```
     docker-build-booking-storyN: ## Baut das BookingService StoryN Docker-Image
     	docker build -f Dockerfile --build-arg SERVICE_PATH=booking/storyN --build-arg SERVICE_DESC="Workshop Microservices - Booking Service (Story N)" --build-arg SERVICE_PORT=8080 -t workshop-microservices-booking-storyN:latest .
     ```

8. **`.github/workflows/build.yml`**
   - Fuege einen neuen Job `build-booking-storyN` am Ende hinzu, nach dem Muster der bestehenden booking-story Jobs:
     ```yaml
       build-booking-storyN:
         name: Build BookingService StoryN
         runs-on: ubuntu-latest
         defaults:
           run:
             working-directory: services
         steps:
           - uses: actions/checkout@v6
           - name: Setup Go
             uses: actions/setup-go@v6
             with:
               go-version-file: services/go.mod
           - name: Build
             run: go build -v -o bin/booking-storyN ./booking/storyN
     ```

9. **`.github/workflows/docker.yml`**
   - Im `changes`-Job: Fuege `booking-storyN: ${{ steps.filter.outputs.booking-storyN }}` als neuen Output hinzu
   - Im `changes`-Job filter: Fuege einen neuen Filter-Eintrag hinzu:
     ```
     booking-storyN:
       - 'services/booking/storyN/**'
       - 'services/go.mod'
     ```
   - Fuege einen neuen Build-Job am Ende hinzu, nach dem Muster der bestehenden booking-story Jobs:
     ```yaml
       build-booking-storyN:
         needs: changes
         if: needs.changes.outputs.booking-storyN == 'true'
         runs-on: ubuntu-latest
         steps:
           - uses: actions/checkout@v6
           - uses: docker/login-action@v4
             with:
               username: ${{ secrets.DOCKERHUB_USERNAME }}
               password: ${{ secrets.DOCKERHUB_TOKEN }}
           - uses: docker/setup-buildx-action@v4
           - uses: docker/build-push-action@v7
             with:
               context: services
               file: services/Dockerfile
               push: true
               platforms: linux/amd64,linux/arm64
               build-args: |
                 SERVICE_PATH=booking/storyN
                 SERVICE_DESC=Workshop Microservices - Booking Service (Story N)
                 SERVICE_PORT=8080
               tags: ${{ secrets.DOCKERHUB_USERNAME }}/workshop-microservices-booking-storyN:latest
     ```

10. **`services/docker-compose.infra.yml`**
    - Fuege in der Swagger-UI `URLS`-Liste einen neuen Eintrag hinzu:
      ```
      { "url": "http://localhost:PORT/openapi", "name": "Booking Service (Story N)" }
      ```
    (PORT = 8084 + N)

11. **`services/docker-compose.reference.yml`**
    - Fuege einen neuen Service am Ende hinzu:
      ```yaml
        booking-storyN:
          build:
            context: .
            dockerfile: Dockerfile
            args:
              SERVICE_PATH: booking/storyN
          image: larmic/workshop-microservices-booking-storyN:latest
          ports:
            - "PORT:8080"
          environment:
            FLIGHT_SERVICE_URL: http://flight:8080
            HOTEL_SERVICE_URL: http://hotel:8080
            CAR_SERVICE_URL: http://car:8080
      ```
    (PORT = 8084 + N)

### Abschluss

Zeige eine Zusammenfassung:
- Welche Dateien erstellt wurden
- Welche Dateien geaendert wurden
- Der zugewiesene Port
- Hinweis: `cd services && go build -v -o /dev/null ./booking/storyN` zum Testen
