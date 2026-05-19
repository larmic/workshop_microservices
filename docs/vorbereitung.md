# Workshop-Vorbereitung

Diese Anleitung beschreibt, wie du das Repo auscheckst, die Umgebung startest und mit dem Workshop loslegst.

## Voraussetzungen

- **Git**
- **Docker** inkl. Docker Compose
- **IDE / Editor** nach Wahl (optional, nur wenn du die Stories selbst implementieren willst)
- Eigene Sprache/Framework deiner Wahl, in der du den Booking-Service umsetzen möchtest (Java/Spring Boot, Quarkus, Go, Node.js, ...)

## 1. Repository klonen

```bash
git clone https://github.com/larmic/workshop_microservices.git
cd workshop_microservices/services
```

Alle weiteren Befehle werden aus dem Verzeichnis `services/` ausgeführt.

## 2. Umgebung starten

Es gibt zwei Wege, die Workshop-Umgebung hochzufahren — beide bringen die gleichen Container ans Laufen (Flight, Hotel, Car, Booking-Referenzlösungen, Traefik, Consul, Swagger UI, Dashboard).

### Variante A — Images von Docker Hub ziehen (empfohlen für Teilnehmer)

Schnell und ohne lokalen Build:

```bash
make docker-up-hub
```

Im Hintergrund: `docker compose ... pull` lädt fertige Images von Docker Hub (`larmic/workshop-microservices-*`), dann `docker compose ... up`.

### Variante B — Images lokal bauen

Wenn du Änderungen am Code oder an den Dockerfiles vornimmst, baust du die Images lokal:

```bash
make docker-up
```

Das entspricht `docker compose ... up --build` über alle drei Compose-Dateien.

### Ohne Makefile

Falls `make` nicht verfügbar ist, gehen beide Varianten auch direkt:

```bash
# Variante A (Docker Hub)
docker compose -f docker-compose.yml -f docker-compose.infra.yml -f docker-compose.reference.yml pull
docker compose -f docker-compose.yml -f docker-compose.infra.yml -f docker-compose.reference.yml up

# Variante B (lokal bauen)
docker compose -f docker-compose.yml -f docker-compose.infra.yml -f docker-compose.reference.yml up --build
```

Alle weiteren Make-Targets siehst du mit `make help`.

## 3. Dashboard öffnen

Das **Dashboard** ist das zentrale Tool für den Workshop. Hier startest, stoppst und skalierst du Services, siehst Health-Status und kommst in einem Klick zu den anderen Tools.

> **<http://localhost>**

![Workshop Dashboard](assets/dashboard.png)

## 4. Weitere URLs

| Tool              | URL                              | Zweck                                |
|-------------------|----------------------------------|--------------------------------------|
| Dashboard         | <http://localhost>               | Workshop-Steuerzentrale              |
| Swagger UI        | <http://localhost/api>           | API-Dokumentation aller Services     |
| Consul UI         | <http://localhost/consul>        | Service Discovery & Health           |
| Traefik Dashboard | <http://localhost:8080>          | Routing & Reverse-Proxy-Monitoring   |

Schneller Health-Check vom Terminal:

```bash
curl http://localhost/api/flight/health
curl http://localhost/api/hotel/health
curl http://localhost/api/car/health
```

## 5. Eigenen Booking-Service einklinken (Workshop-Setup)

Wenn deine Gruppe einen eigenen Booking-Service umsetzt, kannst du ihn neben der Referenz-Lösung ins Setup einklinken und im Dashboard pro Story zwischen Referenz und eurer Lösung umschalten.

### Konzept

Stories bauen kumulativ aufeinander auf — ihr erweitert iterativ **einen** Service, statt für jede Story ein neues Projekt anzulegen. Das Repo bleibt unverändert, eure Lösung lebt in einem **separaten Verzeichnis** (eigenes Repo, irgendwo auf der Platte). Compose baut euer Image über einen konfigurierbaren Pfad.

### Setup

1. **Eigenes Projekt anlegen** — irgendwo, mit einer `Dockerfile`-Datei in der Wurzel. Sprache/Framework frei. Der Service muss auf Port `8080` HTTP entgegennehmen.
2. **`.env` anlegen** im Verzeichnis `services/`:
   ```bash
   cp .env.example .env
   ```
   und in `.env` den absoluten Pfad zu eurem Projekt eintragen:
   ```bash
   WORKSHOP_BOOKING_PATH=/absolute/path/to/your/booking-service
   ```
3. **Stack mit Workshop-Service starten**:
   ```bash
   make docker-up-workshop
   ```
   Die Referenz-Images werden vom Docker Hub gezogen (nicht neu gebaut), nur euer Workshop-Service wird lokal gebaut.
4. **Im Dashboard pro Story** zwischen *Reference* und *Workshop* umschalten — die Buttons rufen dann den jeweils ausgewählten Service auf. Der Workshop-Toggle ist grau, solange euer Service nicht erreichbar ist.
5. **Code-Iteration**: Wenn ihr euren Service-Code geändert habt und das Image neu bauen wollt — ohne den ganzen Stack zu restarten:
   ```bash
   make docker-rebuild-workshop
   ```
   Das baut das Image und ersetzt nur den `booking-workshop`-Container. Alle anderen Container laufen weiter.

### Interne URLs

Euer Service läuft im selben Docker-Netz wie die Referenz-Implementierung. Die internen Hostnamen sind identisch:

| Backend  | URL                          |
|----------|------------------------------|
| Flight   | `http://flight:8080`         |
| Hotel    | `http://hotel:8080`          |
| Car      | `http://car:8080`            |
| Consul   | `http://consul:8500`         |

Diese Werte werden eurem Container als Environment-Variablen `FLIGHT_SERVICE_URL`, `HOTEL_SERVICE_URL`, `CAR_SERVICE_URL`, `CONSUL_URL` übergeben.

### Hinweise für Windows-Teilnehmer

- **Pfade** mit Forward-Slashes statt Backslashes: `WORKSHOP_BOOKING_PATH=C:/Users/foo/projects/my-booking-service`
- **Line endings** der `.env` müssen LF sein (nicht CRLF) — sonst landet `\r` am Pfadende und der Build schlägt mit "no such file" fehl. Editor entsprechend einstellen oder `.env` über WSL/Git Bash erzeugen.
- **Build-Performance**: das Projekt sollte im WSL2-Filesystem liegen (z.B. `~/projects/...`), nicht unter `/mnt/c/...`.

### Workshop-Setup stoppen

```bash
make docker-down-workshop
```

## 6. Los geht's

Dein Arbeitsplatz ist eingerichtet. Starte mit [Story 1: Der erste cloud-native Booking-Service](stories/story-01-cloud-native-booking-service.md).

---

## Trainer-Aufgabe: Feedback einsammeln

Feedback läuft über **GitHub Discussions** im Repo `larmic/workshop_microservices`. Der Link ist auf der Feedback-Slide (Slide 2) sowie als prominenter Button in der Dashboard-Kopfleiste hinterlegt.

Einmalige Einrichtung:

1. Im Repo: Settings → Features → ✓ Discussions
2. Kategorie *Workshop Feedback* anlegen (Format „Open-ended")

Pro Workshop-Run optional ein vorbereiteter Sammel-Thread:

- **Titel:** z.B. `Feedback Kickoff — YYYY-MM-DD`
- **Body:** Datum und die sechs Fragen aus der Disclaimer-Slide (roter Faden, Patterns, Zeitrahmen, Beispiele, Dashboard, Sonstiges)

Falls das Repo umzieht, die URL in `services/slides/chapters/02-feedback.md` und in `services/dashboard/static/index.html` (Feedback-Button) anpassen.
