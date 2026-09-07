---
name: new-booking-story
description: Legt eine neue Booking-Service Story an (Verzeichnisstruktur, Code, OpenAPI, Makefile, GitHub Workflows, Docker Compose, Traefik, Dashboard)
user_invocable: true
---

# Neue Booking Story anlegen

Erstelle eine neue BookingService Story mit allen zugehoerigen Dateien und Build-Konfigurationen.

## Argumente

Das erste Argument ist die Story-Nummer (z.B. `9` fuer Story 9). Wenn keine Nummer angegeben wurde, frage den User danach.

## Nummerierung und Ports

- Code-Stories sind heute `story1` und `story3` bis `story8`. **Story 2 ist eine Design-Session ohne Code** (REST vs. RESTful am Flipchart): kein Verzeichnis, kein Image, kein Port. Eine neue Story bekommt also die naechste freie Nummer nach der hoechsten Code-Story, die Vorlage ist die hoechste vorhandene Code-Story.
- Externer Host-Port: `8084 + N` (Story 1 = 8085, Story 3 = 8087, Story 8 = 8092). `8086` bleibt frei.
- Namensschema ueberall `booking-ref-storyN` (Compose-Service, Traefik-Router, Dashboard-Proxy, Env `BOOKING_REF_STORYN_URL`, Go-Variable `bookingRefStoryNURL`). Docker-Hub-Tag ist `storyN`.

## Ablauf

### Vorbedingungen pruefen

1. Lese die Story-Nummer N aus den Argumenten
2. Bestimme die Vorlage-Story P: die hoechste existierende `services/booking/storyP/` mit P < N (bei Luecken, wie Story 2, die naechste kleinere Code-Story)
3. Pruefe ob `services/booking/storyN/` bereits existiert. Falls ja: Abbruch mit Fehlermeldung.
4. Pruefe ob `docs/stories/story-0N-*.md` existiert. Falls nein: Hinweis ausgeben, dass die Story-Beschreibung fehlt (Kontext, User Story, Akzeptanzkriterien werden daraus uebernommen).
5. Berechne den externen Port: `8084 + N`

### Neue Dateien erstellen

Kopiere `services/booking/storyP/` komplett nach `services/booking/storyN/` (alle Unterverzeichnisse: `api/`, `handler/`, `circuitbreaker/`, `bulkhead/`, `saga/`, `README.md`, …) und passe an:

1. **`main.go`**
   - Import-Pfade `…/booking/storyP/...` zu `…/booking/storyN/...`
   - `Service: "booking-P"` zu `Service: "booking-N"` (Consul-Registrierung, sonst kollidieren Story P und Story N auf demselben Service-Namen)
   - Story-Kommentare im Kopf ("Story P: …") auf die neue Story anpassen

2. **`api/openapi.yaml`**
   - `version: P.0.0` zu `version: N.0.0`
   - `servers.url` auf `…/api/booking-ref-storyN`
   - `example: booking-N` beim Service-Namen (Info-Endpoint)
   - Beschreibungstexte, die "Story P" nennen

3. **Alle weiteren Dateien**: Referenzen auf `storyP` bzw. "Story P" durch `storyN` bzw. "Story N" ersetzen. Danach `gofmt -l services` und `go build -C services ./...`.

### Bestehende Dateien erweitern

4. **`services/Makefile`**
   - `.PHONY`-Zeile: `run-booking-ref-storyN` und `docker-build-booking-ref-storyN` ergaenzen
   - `docker-build-booking-ref-storyN` als Dependency zum `docker-build`-Target
   - Neues Target nach dem letzten `docker-build-booking-ref-story*`-Target:
     ```
     docker-build-booking-ref-storyN: ## Baut das BookingService Reference StoryN Docker-Image
     	docker build -f Dockerfile --build-arg SERVICE_PATH=booking/storyN --build-arg SERVICE_DESC="Workshop Microservices - Booking Service (Reference Story N)" --build-arg SERVICE_PORT=8080 -t workshop-microservices-booking:storyN .
     ```
   - `run-booking-ref-storyN` analog zu den bestehenden `run-*`-Targets

5. **`.github/workflows/build.yml`**
   - Neuen Job nach `build-booking-ref-storyP` einfuegen, Action-Versionen von den bestehenden Jobs uebernehmen (nicht raten, Renovate haelt sie aktuell):
     ```yaml
       build-booking-ref-storyN:
         name: Build BookingService Reference StoryN
         runs-on: ubuntu-latest
         defaults:
           run:
             working-directory: services
         steps:
           - uses: actions/checkout@<Version wie oben>

           - name: Setup Go
             uses: actions/setup-go@<Version wie oben>
             with:
               go-version-file: services/go.mod

           - name: Build
             run: go build -v -o bin/booking-ref-storyN ./booking/storyN
     ```

6. **`.github/workflows/docker.yml`**
   - Pfad-Filter `services/booking/storyN/**` in der `booking`-Filterliste ergaenzen
   - Matrix `story: [1, 3, 4, 5, 6, 7, 8]` um `N` erweitern. Der Job pusht das Image als Tag `storyN` ins gemeinsame Repo `workshop-microservices-booking`.

7. **`.github/workflows/dockerhub-description.yml`** und **`.github/dockerhub/booking.md`**
   - Short-Description (`story1, story3..story8, custom`) und Tag-Tabelle um `storyN` ergaenzen

8. **`services/docker-compose.infra.yml`**
   - Swagger-UI `URLS`-Liste:
     ```
     { "url": "http://localhost/api/booking-ref-storyN/openapi", "name": "Booking Service (Reference Story N)" },
     ```

9. **`services/docker-compose.reference.yml`**
   - Neuer Service am Ende, `environment`- und `depends_on`-Block 1:1 von `booking-ref-storyP` uebernehmen (ab Story 3 Consul, Story 1 nutzt statische `*_SERVICE_URL`-Variablen):
     ```yaml
       booking-ref-storyN:
         container_name: booking-ref-storyN
         build:
           context: .
           dockerfile: Dockerfile
           args:
             SERVICE_PATH: booking/storyN
         image: larmic/workshop-microservices-booking:storyN
         ports:
           - "PORT:8080"
         environment:
           CONSUL_URL: http://consul:8500
         depends_on:
           consul:
             condition: service_healthy
     ```
     (PORT = 8084 + N)

10. **`services/traefik/dynamic.yml`**
    - `http.routers`: Router `booking-ref-storyN` nach `booking-ref-storyP` (Muster: `rule: "PathPrefix(\`/api/booking-ref-storyN\`)"`, `entryPoints: [web]`, `middlewares: [booking-ref-storyN-stripprefix]`, `service: booking-ref-storyN`)
    - `http.middlewares`: `booking-ref-storyN-stripprefix` mit `stripPrefix.prefixes: ["/api/booking-ref-storyN"]`
    - `http.services`: `booking-ref-storyN` mit `loadBalancer.servers[0].url: "http://booking-ref-storyN:8080"`
    - Kontrolle: `grep -oE 'booking-ref-story[0-9]' services/traefik/dynamic.yml | sort | uniq -c` liefert fuer jede Story dieselbe Anzahl

### Dashboard erweitern

Das Dashboard (`services/dashboard/`) zeigt pro Story einen Stepper-Eintrag, einen API-Link und einen Inhaltsbereich. Das Startup-Overlay und das Health-Panel iterieren ueber die `bookingURLs`-Map in `main.go`, neue Stories erscheinen dort automatisch.

11. **`services/dashboard/main.go`**
    - Nach `bookingRefStoryPURL` eine neue URL-Variable:
      ```go
      bookingRefStoryNURL := getEnv("BOOKING_REF_STORYN_URL", "http://booking-ref-storyN:8080")
      ```
    - Eintrag in der `bookingURLs`-Map:
      ```go
      "booking-ref-storyN": bookingRefStoryNURL,
      ```
    - Proxy-Route fuer den Offers-Aufruf, falls die Story-Section einen REST-Button bekommt:
      ```go
      mux.HandleFunc("GET /api/booking-ref-storyN/offers", handler.ProxyHandler(bookingRefStoryNURL, http.MethodGet, "/booking/offers"))
      ```
      Story-spezifische Endpunkte (wie `saga-state` in Story 6) bei Bedarf separat.

12. **`services/dashboard/static/index.html`**
    - **API-Link** nach dem letzten `Booking Reference Story P`-Eintrag in der `.links-grid` (Sektion "Service APIs"):
      ```html
      <a class="link-chip" href="/api/booking-ref-storyN/openapi" target="_blank"><span>&#128214;</span> Booking Reference Story N</a>
      ```
    - **Stepper-Button** nach dem letzten `data-story-node="P"`-Button in `.stepper-nodes`:
      ```html
      <span class="stepper-connector"></span>
      <button type="button" class="stepper-node" data-story-node="N" onclick="showStory(N)">N</button>
      ```
    - **`STORY_META`-Eintrag** im JavaScript-Block (nach Eintrag P):
      ```javascript
      N: { title: "Story N: <Titel aus docs/stories/story-0N-*.md>",
           subtitle: "<Kurz-Subtitle, ein Satz>" },
      ```
    - **`STORY_COUNT`** auf `N` setzen.
    - **`BookingMode`-Literal** um `N: "reference"` ergaenzen (Reference/Custom-Umschalter).
    - **Neue `<section class="story-section" data-story="N" hidden>`** am Ende von `<main class="story-content">`. Minimaler Stub mit Story-Info-Block (Kontext, User Story, Akzeptanzkriterien aus `docs/stories/story-0N-*.md`) und optional einem Cheatsheet-Block. **Keine** story-spezifischen UI-Buttons im Stub, die fuegen die Workshop-Teilnehmer beim Bearbeiten der Story selbst hinzu. Vorlage:
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

### Doku und Slides

13. `docs/stories/story-0N-*.md`, `docs/questions/storyN.md`, Slides `services/slides/chapters/*-story-0N.md` und `*-story-0N-fragen.md` plus Eintrag in `services/slides/index.html`, Agenda (`03-agenda.md`), Zusammenfassungstabelle (`29-zusammenfassung.md`), Zeitleiste in `docs/themen.md`, `docs/troubleshooting.md` (Port-Tabelle), `README.md` und `CLAUDE.md` (Story-Bereiche).

### Abschluss

Zeige eine Zusammenfassung:
- Welche Dateien erstellt wurden
- Welche Dateien geaendert wurden
- Der zugewiesene Port
- Verifikation (vom Repo-Root, kein `cd`): `go build -C services ./...`, `gofmt -l services`, `docker compose --project-directory services -f services/docker-compose.yml -f services/docker-compose.infra.yml -f services/docker-compose.reference.yml config --quiet`
