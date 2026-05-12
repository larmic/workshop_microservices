---
name: new-booking-story
description: Legt eine neue Booking-Service Story an (Verzeichnisstruktur, Code, OpenAPI, HTTP-Tests, Makefile, GitHub Workflows, Docker Compose, Dashboard)
user_invocable: true
---

# Neue Booking Story anlegen

Erstelle eine neue BookingService Story mit allen zugehoerigen Dateien und Build-Konfigurationen.

## Argumente

Das erste Argument ist die Story-Nummer (z.B. `3` fuer Story 3). Wenn keine Nummer angegeben wurde, frage den User danach.

## Ablauf

### Vorbedingungen pruefen

1. Lese die Story-Nummer N aus den Argumenten
2. Berechne die vorherige Story-Nummer: P = N - 1
3. Pruefe ob `services/booking/storyN/` bereits existiert. Falls ja: Abbruch mit Fehlermeldung.
4. Pruefe ob `services/booking/storyP/` existiert. Falls nein: Abbruch mit Fehlermeldung ("Story P existiert nicht, kann nicht als Vorlage verwendet werden").
5. Berechne den externen Port: `8084 + N` (Story1=8085, Story2=8086, Story3=8087, usw.)

### Neue Dateien erstellen

Lese jeweils die Datei aus `services/booking/storyP/` (der vorherigen Story) als Template und erstelle die angepasste Version unter `services/booking/storyN/`:

1. **`services/booking/storyN/main.go`**
   - Kopiere von storyP/main.go
   - Aendere `Service: "booking-P"` zu `Service: "booking-N"` (Consul-Registrierung &mdash; sonst kollidieren Story P und Story N auf demselben Service-Namen)
   - Beachte: Die Imports auf `booking/storyP/...` werden vom "Wichtig"-Hinweis unten generisch durch `booking/storyN/...` ersetzt

2. **`services/booking/storyN/api/openapi.yaml`**
   - Kopiere von storyP/api/openapi.yaml
   - Aendere `version: P.0.0` zu `version: N.0.0`

3. **`services/booking/storyN/http/requests.http`**
   - Kopiere von storyP/http/requests.http
   - Aendere den Header-Kommentar von "Story P" zu "Story N"

4. **`services/booking/storyN/http/http-client.env.json`**
   - Kopiere von storyP/http/http-client.env.json
   - Aendere `bookingPort` in beiden Environments (local und docker) auf den berechneten Port (8084+N)

5. **`services/booking/storyN/http/Makefile`**
   - Kopiere 1:1 von storyP/http/Makefile (keine Aenderungen noetig)

6. **`services/booking/storyN/http/README.md`**
   - Kopiere von storyP/http/README.md
   - Aendere "Story P" zu "Story N" im Titel und Text

**Wichtig:** Kopiere auch alle weiteren Dateien und Unterverzeichnisse, die in `storyP/` existieren aber oben nicht explizit aufgefuehrt sind (z.B. zusaetzliche Handler-Dateien, Packages, etc.). Diese werden 1:1 kopiert, wobei Referenzen auf "storyP" durch "storyN" ersetzt werden.

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
      { "url": "http://localhost/api/booking-storyN/openapi", "name": "Booking Service (Story N)" }
      ```

11. **`services/docker-compose.reference.yml`**
    - Fuege einen neuen Service am Ende hinzu. **Wichtig:** Uebernimm den `environment`- und `depends_on`-Block 1:1 vom `booking-story{P}`-Service (ab Story 2 wird Consul verwendet, aeltere Stories nutzten statische `*_SERVICE_URL`-Variablen):
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
            CONSUL_URL: http://consul:8500
          depends_on:
            consul:
              condition: service_healthy
      ```
    (PORT = 8084 + N)

12. **`services/traefik/dynamic.yml`**
    - Fuege im `http.routers`-Block einen neuen Router nach `booking-story{P}` hinzu:
      ```yaml
          booking-storyN:
            rule: "PathPrefix(`/api/booking-storyN`)"
            entryPoints:
              - web
            middlewares:
              - booking-storyN-stripprefix
            service: booking-storyN
      ```
    - Fuege im `http.middlewares`-Block eine neue Middleware nach `booking-story{P}-stripprefix` hinzu:
      ```yaml
          booking-storyN-stripprefix:
            stripPrefix:
              prefixes:
                - "/api/booking-storyN"
      ```
    - Fuege im `http.services`-Block einen neuen Service nach `booking-story{P}` hinzu:
      ```yaml
          booking-storyN:
            loadBalancer:
              servers:
                - url: "http://booking-storyN:8080"
      ```

### Dashboard erweitern

Das Dashboard (`services/dashboard/`) zeigt pro Story einen Stepper-Eintrag, einen API-Link und einen Inhaltsbereich an. Ausserdem wartet die Startup-Overlay auf den Health-Check der neuen Story, damit der Workshop erst startet, wenn alle Booking-Services bereit sind.

13. **`services/dashboard/main.go`**
    - Fuege nach `bookingStory{P}URL` eine neue URL-Variable hinzu:
      ```go
      bookingStoryNURL := getEnv("BOOKING_STORYN_URL", "http://booking-storyN:8080")
      ```
    - Fuege in der `bookingURLs`-Map einen neuen Eintrag hinzu:
      ```go
      "booking-storyN": bookingStoryNURL,
      ```
    - **Hinweis:** Damit erscheint Story N automatisch im `/api/health-overview` und in der Startup-Overlay &mdash; keine weiteren Handler noetig, solange Story N keine eigenen Dashboard-API-Endpunkte braucht (Story-spezifische Endpunkte wie `saga-state` aus Story 5 werden bei Bedarf separat hinzugefuegt).

14. **`services/dashboard/static/index.html`**
    - **API-Link** nach dem letzten `Booking Story P`-Eintrag in der `.links-grid` (Sektion "Service APIs"):
      ```html
      <a class="link-chip" href="/api/booking-storyN/openapi" target="_blank"><span>&#128214;</span> Booking Story N</a>
      ```
    - **Stepper-Button** nach dem letzten `data-story-node="P"`-Button in `.stepper-nodes`:
      ```html
      <span class="stepper-connector"></span>
      <button type="button" class="stepper-node" data-story-node="N" onclick="showStory(N)">N</button>
      ```
    - **STORY_META-Eintrag** im JavaScript-Block (nach Eintrag P):
      ```javascript
      N: { title: "Story N: <Titel aus docs/stories/story-NN-*.md>",
           subtitle: "<Kurz-Subtitle, ein Satz>" },
      ```
    - **STORY_COUNT** inkrementieren:
      ```javascript
      const STORY_COUNT = N;
      ```
    - **Neue `<section class="story-section" data-story="N" hidden>`** am Ende von `<main class="story-content">` einfuegen. Minimaler Stub mit Story-Info-Block (Kontext, User Story, Akzeptanzkriterien aus `docs/stories/story-NN-*.md`) und optional einem Cheatsheet-Block. **Keine** story-spezifischen UI-Buttons im Stub &mdash; die fuegen die Workshop-Teilnehmer beim Bearbeiten der Story selbst hinzu. Vorlage (Story-Inhalte aus der Doku uebernehmen):
      ```html
      <section class="story-section" data-story="N" hidden>
          <div class="story-helpers" data-story-helpers="N">
              <details class="story-info">
                  <summary>Story lesen <span class="badge">User Story + Akzeptanzkriterien</span></summary>
                  <div class="story-info-body">
                      <h4>Kontext</h4>
                      <p>...</p>
                      <h4>User Story</h4>
                      <p>Als <em>...</em> moechte ich <em>...</em>, damit <em>...</em>.</p>
                      <h4>Akzeptanzkriterien</h4>
                      <ul>
                          <li>...</li>
                      </ul>
                  </div>
              </details>
          </div>
      </section>
      ```

### Abschluss

Zeige eine Zusammenfassung:
- Welche Dateien erstellt wurden
- Welche Dateien geaendert wurden
- Der zugewiesene Port
- Hinweis: `cd services && go build -v -o /dev/null ./booking/storyN` zum Testen
